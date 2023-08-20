package restAPIs

import (
	"encoding/json"
	"fmt"
	"html/template"
	"httpServer/utils"
	"math"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

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

func (api *GitHubRestAPI) GetRepositories(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	org := vars["org"]
	phrase, ok := vars["q"]
	if ok {
		fmt.Fprintf(w, "You've requested the repositories of: %s with the phrase: %s\n", org, phrase)
	} else {
		fmt.Fprintf(w, "You've requested the repositories of: %s\n", org)
	}
	apiURL := fmt.Sprintf("https://api.github.com/search/repositories?q=%s+in:name+org:%s", phrase, org)
	fmt.Println(apiURL)

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

	// Use pagination utilities
	currentPage := paginationHandler.GetCurrentPage(r)
	perPage := 10 // Number of items per page
	totalItems := len(jsonResponse.Items)
	totalPages := int(math.Ceil(float64(totalItems) / float64(perPage)))
	startIndex, endIndex := paginationHandler.GetPageIndices(currentPage, perPage, totalItems)
	paginatedResponse := jsonResponse.Items[startIndex:endIndex]
	pageNumbers := paginationHandler.GetPageNumbers(currentPage, totalPages)

	tmpl, err := template.ParseFiles("templates/repositories.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	queryParams := paginationHandler.GetQueryParams(r)

	data := struct {
		Repositories []struct {
			Name  string `json:"name"`
			Owner struct {
				Login string `json:"login"`
			} `json:"owner"`
			URL          string `json:"html_url"`
			CreationTime string `json:"created_at"`
			Stars        int    `json:"stargazers_count"`
		}
		CurrentPage int
		PageNumbers []int
		BaseUrl     string
		QueryParams map[string]string
	}{
		Repositories: paginatedResponse,
		CurrentPage:  currentPage,
		PageNumbers:  pageNumbers,
		BaseUrl:      r.URL.Path,
		QueryParams:  queryParams,
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
