package repositoriesCollectors

import (
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"httpServer/cacheDBs"
	"net/http"
)

// ReposCollector defines the interface for repositories information collector.
type ReposCollector interface {
	// ConfigureCollector configures the collector with the provided context, Redis client, and router.
	ConfigureCollector(ctx context.Context, rdb *redis.Client, route *mux.Router, cacheDB cacheDBs.CacheDB)

	// GetRepositories handles the HTTP request to retrieve repositories based on defined query params.
	GetRepositories(w http.ResponseWriter, r *http.Request)
}
