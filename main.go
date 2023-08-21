package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"httpServer/cacheDBs"
	"httpServer/repositoriesCollectors"
	"log"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()

	// Configure Redis client based on environment
	rdb := configureRedisClient()

	// Create Redis cache
	cache := cacheDBs.NewRedisDB(rdb)

	// Create a new Gorilla Mux router
	r := mux.NewRouter()

	// Configure GitHub API collector
	var gitHubAPI repositoriesCollectors.GitHubReposCollector
	gitHubAPI.ConfigureCollector(ctx, rdb, r, cache)

	// Define routes
	gitHubAPI.SetupRoutes(r)

	// Start the HTTP server
	startServer(r)
}

func configureRedisClient() *redis.Client {
	runningInDocker := os.Getenv("RUNNING_IN_DOCKER")
	redisAddr := "localhost:6379" // Default address for local environment

	if runningInDocker == "true" {
		redisAddr = "my-redis:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	return rdb
}

func startServer(r *mux.Router) {
	addr := ":8080"
	fmt.Printf("Server available at http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
