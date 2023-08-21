package httpServer_tests

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"httpServer/cacheDBs"
	"httpServer/repositoriesCollectors"
	"net/http"
	"net/http/httptest"
	"testing"
)

var ctx = context.Background()
var rdb *redis.Client

func TestGetRepositories(t *testing.T) {
	rdb = redis.NewClient(&redis.Options{
		Addr: "my-redis:6379",
	})

	cache := cacheDBs.NewRedisDB(rdb)

	r := mux.NewRouter()

	var gitHubAPI repositoriesCollectors.GitHubReposCollector
	gitHubAPI.ConfigureCollector(ctx, rdb, r, cache)

	r.HandleFunc("/repositories/org/{org}", gitHubAPI.GetRepositories).Methods("GET")
	r.HandleFunc("/repositories/org/{org}/q/{q}", gitHubAPI.GetRepositories).Methods("GET")

	req, err := http.NewRequest("GET", "/repositories/org/github", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"repositories":[]}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
