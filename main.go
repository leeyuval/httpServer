package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"httpServer/cacheDBs"
	"httpServer/repositoriesCollectors"
	"log"
	"net/http"
)

var ctx = context.Background()
var rdb *redis.Client

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr: "my-redis:6379",
	})

	cache := cacheDBs.NewRedisDB(rdb)

	r := mux.NewRouter()

	var gitHubAPI repositoriesCollectors.GitHubReposCollector
	gitHubAPI.ConfigureCollector(ctx, rdb, r, cache)

	r.HandleFunc("/repositories/org/{org}", gitHubAPI.GetRepositories).Methods("GET")
	r.HandleFunc("/repositories/org/{org}/q/{q}", gitHubAPI.GetRepositories).Methods("GET")

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
