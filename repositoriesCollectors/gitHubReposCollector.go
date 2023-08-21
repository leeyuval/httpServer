package repositoriesCollectors

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

// ** Types **

// GitHubReposCollector collects and manages GitHub repository information.
type GitHubReposCollector struct {
	ctx   context.Context // The context for handling requests.
	rdb   *redis.Client   // The Redis client for caching data.
	route *mux.Router     // The router for defining API endpoints.
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
func (api *GitHubReposCollector) ConfigureCollector(ctx context.Context, rdb *redis.Client, route *mux.Router) {
	api.ctx = ctx
	api.rdb = rdb
	api.route = route
}

// getCachedData retrieves cached GitHub JSON response data based on a cache key.
// It returns the cached data or an error if retrieval or decoding fails.
func (api *GitHubReposCollector) getCachedData(cacheKey string) (GitHubJsonResponse, error) {
	cachedData, err := api.rdb.Get(api.ctx, cacheKey).Result()
	if err != nil {
		return GitHubJsonResponse{}, err
	}

	var jsonResponse GitHubJsonResponse
	err = json.Unmarshal([]byte(cachedData), &jsonResponse)
	if err != nil {
		return GitHubJsonResponse{}, err
	}

	return jsonResponse, nil
}

// storeDataInCache stores GitHub JSON response data in the cache using a cache key.
func (api *GitHubReposCollector) storeDataInCache(cacheKey string, data GitHubJsonResponse) {
	cachedDataBytes, _ := json.Marshal(data)
	result := api.rdb.Set(api.ctx, cacheKey, cachedDataBytes, time.Hour*12)
	if result.Err() == nil {
		log.Println("Data stored in cache successfully.")
	} else {
		log.Println(result.Err())
	}
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

	// Check if data is in cache
	jsonResponse, err := api.getCachedData(cacheKey + ":" + page)
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
		api.storeDataInCache(cacheKey+":"+page, jsonResponse)
	}

	title := generateHtmlTitle(org, phrase)

	data := struct {
		Repositories []Repository
		Title        string
		Page         int
		PrevPage     int
		NextPage     int
		HasPrev      bool
		HasNext      bool
		TotalPages   int
	}{
		Repositories: convertToRepositories(jsonResponse.Items),
		Title:        title,
		Page:         pageNum,
		PrevPage:     pageNum - 1,
		NextPage:     pageNum + 1,
		HasPrev:      pageNum > 1,
		HasNext:      pageNum < jsonResponse.TotalPages,
		TotalPages:   jsonResponse.TotalPages,
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
