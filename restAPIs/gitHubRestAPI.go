package restAPIs

import (
	"encoding/json"
	"fmt"
	"html/template"
	"httpServer/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

const PerPage = 10

// GitHubRestAPI is an implementation of RestAPI interface for GitHub repositories.
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

func (api *GitHubRestAPI) ConfigureRestAPI(ctx context.Context, rdb *redis.Client, route *mux.Router) {
	api.ctx = ctx
	api.rdb = rdb
	api.route = route
}

type Repository struct {
	Name         string
	Owner        string
	URL          string
	CreationTime string
	Stars        int
}

func paginateRepositories(jsonResponse GitHubJsonResponse, page int) ([]Repository, []int, int) {
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

func (api *GitHubRestAPI) GetRepositories(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	org := vars["org"]
	phrase, ok := vars["q"]
	if !ok {
		phrase = ""
	}

	title := fmt.Sprintf("GitHub Repositories of '%s'", org)

	if ok && phrase != "" {
		title = fmt.Sprintf("GitHub Repositories of '%s' including the phrase '%s'", org, phrase)
	}

	apiURL := fmt.Sprintf("https://api.github.com/search/repositories?q=%s+in:name+org:%s", phrase, org)

	// Fetch repositories from Redis cache if available
	cacheKey := org + ":" + phrase
	cacheResult, err := api.rdb.Get(api.ctx, cacheKey).Result()
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cacheResult))
		return
	}

	// Fetch repositories from GitHub API
	resp, err := http.Get(apiURL)
	if err != nil {
		http.Error(w, "Error fetching data from GitHub API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var jsonResponse GitHubJsonResponse

	err = json.NewDecoder(resp.Body).Decode(&jsonResponse)
	if err != nil {
		http.Error(w, "Error decoding GitHub API response", http.StatusInternalServerError)
		return
	}

	// Cache the JSON response in Redis
	cacheTTL := 1 // Cache for 1 second
	if phrase != "" {
		// Longer cache TTL when phrase is included
		cacheTTL = 10 // Cache for 10 seconds
	}
	cacheResultJSON, _ := json.Marshal(jsonResponse)
	api.rdb.Set(api.ctx, cacheKey, cacheResultJSON, time.Duration(cacheTTL)*time.Second)

	// Pagination
	pageParam := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageParam)
	if page <= 0 {
		page = 1
	}

	// Paginate repositories using the new function
	paginatedResponse, pageNumbers, totalPages := paginateRepositories(jsonResponse, page)

	tmpl, err := template.ParseFiles("templates/repositories.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
