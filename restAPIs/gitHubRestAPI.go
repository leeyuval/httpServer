package restAPIs

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"httpServer/utils"
	"log"
	"net/http"
	"strconv"
	"time"
)

const PerPage = 10

type GitHubRestAPI struct {
	ctx   context.Context
	rdb   *redis.Client
	route *mux.Router
}

type GitHubJsonResponse struct {
	Items []struct {
		Name  string `json:"name"`
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
		URL          string `json:"html_url"`
		CreationTime string `json:"created_at"`
		Stars        int    `json:"stargazers_count"`
	} `json:"items"`
}

type Repository struct {
	Name         string
	Owner        string
	URL          string
	CreationTime string
	Stars        int
}

func (api *GitHubRestAPI) ConfigureRestAPI(ctx context.Context, rdb *redis.Client, route *mux.Router) {
	api.ctx = ctx
	api.rdb = rdb
	api.route = route
}

func (api *GitHubRestAPI) GetRepositories(w http.ResponseWriter, r *http.Request) {
	org, phrase := getOrgNameAndPhrase(r)
	cacheKey := fmt.Sprintf("%s:%s", org, phrase)
	pageParam := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageParam)
	if page < 1 {
		page = 1
	}

	var jsonResponse GitHubJsonResponse
	var paginatedResponse []Repository
	var pageNumbers []int
	var totalPages int

	// Check if data is in cache
	cachedData, err := api.rdb.Get(api.ctx, cacheKey).Result()
	if err == nil {
		err := json.Unmarshal([]byte(cachedData), &jsonResponse)
		if err != nil {
			http.Error(w, "Error decoding cached data", http.StatusInternalServerError)
			return
		}

		paginatedResponse, pageNumbers, totalPages = paginate(jsonResponse, page)
	} else {
		apiURL := fmt.Sprintf("https://api.github.com/search/repositories?q=%s+in:name+org:%s", phrase, org)

		// Fetch repositories from GitHub API
		resp, err := http.Get(apiURL)
		if err != nil {
			http.Error(w, "Error fetching data from GitHub API", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
		if err != nil {
			http.Error(w, "Error decoding GitHub API response", http.StatusInternalServerError)
			return
		}

		paginatedResponse, pageNumbers, totalPages = paginate(jsonResponse, page)

		// Store data in cache
		cachedDataBytes, _ := json.Marshal(jsonResponse)
		result := api.rdb.Set(api.ctx, cacheKey, cachedDataBytes, time.Hour*12)
		if result.Err() == nil {
			log.Println("Data stored in cache successfully.")
		} else {
			log.Println(result.Err())
		}
	}

	title := generateHtmlTitle(org, phrase)

	data := struct {
		Repositories []Repository
		Title        string
		CurrentPage  int
		PageNumbers  []int
		TotalPages   int
	}{
		Repositories: paginatedResponse,
		Title:        title,
		CurrentPage:  page,
		PageNumbers:  pageNumbers,
		TotalPages:   totalPages,
	}

	// Render the HTML template
	err = utils.RenderHTMLTemplate(w, "templates/repositories.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func paginate(jsonResponse GitHubJsonResponse, page int) ([]Repository, []int, int) {
	startIndex := (page - 1) * PerPage
	endIndex := startIndex + PerPage
	if endIndex > len(jsonResponse.Items) {
		endIndex = len(jsonResponse.Items)
	}
	paginatedResponse := jsonResponse.Items[startIndex:endIndex]

	var repositories []Repository
	for _, repo := range paginatedResponse {
		repositories = append(repositories, Repository{
			Name:         repo.Name,
			Owner:        repo.Owner.Login,
			URL:          repo.URL,
			CreationTime: utils.FormatCreationTime(repo.CreationTime),
			Stars:        repo.Stars,
		})
	}

	totalPages := (len(jsonResponse.Items) + PerPage - 1) / PerPage
	pageNumbers := make([]int, totalPages)
	for i := range pageNumbers {
		pageNumbers[i] = i + 1
	}

	return repositories, pageNumbers, totalPages
}

func getOrgNameAndPhrase(r *http.Request) (string, string) {
	vars := mux.Vars(r)
	org := vars["org"]
	phrase, ok := vars["q"]
	if !ok {
		phrase = ""
	}
	return org, phrase
}

func generateHtmlTitle(org, phrase string) string {
	title := fmt.Sprintf("GitHub Repositories of '%s'", org)
	if phrase != "" {
		title = fmt.Sprintf("GitHub Repositories of '%s' including the phrase '%s'", org, phrase)
	}
	return title
}
