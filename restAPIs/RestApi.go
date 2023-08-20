package restAPIs

import (
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"net/http"
)

type RestAPI interface {
	ConfigureRestAPI(ctx context.Context, rdb *redis.Client, route *mux.Router)
	GetRepositories(w http.ResponseWriter, r *http.Request)
}
