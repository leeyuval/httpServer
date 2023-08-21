package repositoriesCollector_tests

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"httpServer/cacheDBs"
	"httpServer/repositoriesCollectors"
	"net/http"
	"net/http/httptest"
	"testing"
)

var ctx = context.Background()

func TestGetRepositories(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Create a new Gorilla Mux router
	r := mux.NewRouter()

	// Create Redis cache
	cache := cacheDBs.NewRedisDB(rdb)

	// Configure GitHub API collector
	var gitHubAPI repositoriesCollectors.GitHubReposCollector
	gitHubAPI.ConfigureCollector(ctx, rdb, r, cache)

	// Define routes
	gitHubAPI.SetupRoutes(r)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/repositories?org=github&phrase=go&page=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(gitHubAPI.GetRepositories)

	// Serve the HTTP request and record the response
	handler.ServeHTTP(rr, req)

	// Check the status code of the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the body of the response
	var jsonResponse repositoriesCollectors.GitHubJsonResponse
	err = json.NewDecoder(rr.Body).Decode(&jsonResponse)
	if err != nil {
		t.Fatal(err)
	}

	if len(jsonResponse.Items) == 0 {
		t.Errorf("handler returned empty items")
	}
}
