package repositoriesCollectors

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"httpServer/cacheDBs"
	"httpServer/utils"
	"net/http"
	"strconv"
)

// ** Constants **

const (
	ItemsPerPage = 30
)

// ** Types **

// GitHubReposCollector collects and manages GitHub repository information.
type GitHubReposCollector struct {
	ctx     context.Context  // The context for handling requests.
	rdb     *redis.Client    // The Redis client for caching data.
	route   *mux.Router      // The router for defining API endpoints.
	cacheDB cacheDBs.CacheDB // The cache database for retrieving and storing data.
}

// GitHubJsonResponse defines the structure of the JSON response from the GitHub API.
type GitHubJsonResponse struct {
	TotalPages int                      `json:"total_count"`
	Items      []GitHubJsonResponseItem `json:"items"`
}

// GitHubJsonResponseItem represents a single item in the GitHubJsonResponse.
type GitHubJsonResponseItem struct {
	Name  string `json:"name"`
	Owner struct {
		Login string `json:"login"`
	} `json:"owner"`
	URL          string `json:"html_url"`
	CreationTime string `json:"created_at"`
	Stars        int    `json:"stargazers_count"`
}

// Repository represents a GitHub repository's information structure.
type Repository struct {
	Name         string
	Owner        string
	URL          string
	CreationTime string
	Stars        int
}

// ** Methods **

// ConfigureCollector configures the GitHubReposCollector with necessary dependencies.
func (api *GitHubReposCollector) ConfigureCollector(ctx context.Context, rdb *redis.Client, route *mux.Router, cacheDB cacheDBs.CacheDB) {
	api.ctx = ctx
	api.rdb = rdb
	api.route = route
	api.cacheDB = cacheDB
}

func (api *GitHubReposCollector) SetupRoutes(r *mux.Router) {
	r.HandleFunc("/repositories/org/{org}", api.GetRepositories).Methods("GET")
	r.HandleFunc("/repositories/org/{org}/q/{q}", api.GetRepositories).Methods("GET")
}

// GetRepositories is an HTTP handler that fetches and displays GitHub repositories based on organization and phrase.
func (api *GitHubReposCollector) GetRepositories(w http.ResponseWriter, r *http.Request) {
	org, phrase := getOrgNameAndPhrase(r)
	cacheKey := fmt.Sprintf("%s:%s", org, phrase)

	// Get page number from query parameters
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	pageNum, _ := strconv.Atoi(page)

	var jsonResponse GitHubJsonResponse
	err := api.cacheDB.ExtractCachedData(api.ctx, cacheKey+":"+page, &jsonResponse)
	if err != nil {
		apiURL := fmt.Sprintf("https://api.github.com/search/repositories?q=%s+in:name+org:%s&page=%s", phrase, org, page)

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

		// Store data in cache
		err = api.cacheDB.StoreDataInCache(api.ctx, cacheKey+":"+page, jsonResponse)
		if err != nil {
			http.Error(w, "Error storing data in cache", http.StatusInternalServerError)
			return
		}
	}

	title := generateHtmlTitle(org, phrase)

	paginationData := utils.GetPaginationData(pageNum, ItemsPerPage, jsonResponse.TotalPages)

	data := struct {
		Repositories []Repository
		Title        string
		Pagination   utils.PaginationData
	}{
		Repositories: convertToRepositories(jsonResponse.Items),
		Title:        title,
		Pagination:   paginationData,
	}

	// Render the HTML template
	err = utils.RenderHTMLTemplate(w, "templates/repositories.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ** Functions **

// convertToRepositories converts GitHubJsonResponse items to Repository items.
func convertToRepositories(items []GitHubJsonResponseItem) []Repository {
	var repositories []Repository
	for _, item := range items {
		repo := Repository{
			Name:         item.Name,
			Owner:        item.Owner.Login,
			URL:          item.URL,
			CreationTime: utils.FormatCreationTime(item.CreationTime),
			Stars:        item.Stars,
		}
		repositories = append(repositories, repo)
	}
	return repositories
}

// getOrgNameAndPhrase extracts organization and query phrase from the request.
func getOrgNameAndPhrase(r *http.Request) (string, string) {
	vars := mux.Vars(r)
	org := vars["org"]
	phrase, ok := vars["q"]
	if !ok {
		phrase = ""
	}
	return org, phrase
}

// generateHtmlTitle generates the HTML title for the repository listing page.
func generateHtmlTitle(org, phrase string) string {
	title := fmt.Sprintf("GitHub Repositories of '%s'", org)
	if phrase != "" {
		title = fmt.Sprintf("GitHub Repositories of '%s' including the phrase '%s'", org, phrase)
	}
	return title
}
