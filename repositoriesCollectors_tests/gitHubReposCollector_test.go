package repositoriesCollectors

import (
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"httpServer/repositoriesCollectors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestServerConnectivity(t *testing.T) {
	t.Run("TestServerConnectivity", func(t *testing.T) {
		// Perform a request to the Gorilla Mux server
		resp, err := http.Get("http://localhost:8080/repositories/org/github")
		if err != nil {
			t.Fatalf("Error sending request to Gorilla Mux: %v", err)
		}
		defer resp.Body.Close()

		// Check the response status code
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
		}
	})
}

func TestRedisCache(t *testing.T) {
	t.Run("TestRedisCache", func(t *testing.T) {
		// Set the Redis server address as an environment variable
		os.Setenv("REDIS_ADDRESS", "my-redis:6379")

		// Create a test server using httptest
		testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Simulate the behavior of your server's route handler
			// You can customize this to match your server's behavior
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Test response"))
		}))
		defer testServer.Close()

		// Create a new context for the test
		ctx := context.Background()

		// Create a Redis client
		rdb := createRedisClient(ctx, t)
		pong, err := rdb.Ping(ctx).Result()
		if err != nil {
			t.Fatalf("Error pinging redis: %v", err)
		} else {
			log.Printf("Redis ping result: %v", pong)
		}
		defer rdb.Close()

		// Clear any existing cache entries
		clearCache(ctx, rdb, t)

		// Perform a request to the server to populate the cache
		_, err = http.Get(testServer.URL + "/repositories/org/github")
		if err != nil {
			t.Fatalf("Error sending request to server: %v", err)
		}

		// Perform the same request again and check if the response is fetched from the cache
		resp, err := http.Get(testServer.URL + "/repositories/org/github")
		if err != nil {
			t.Fatalf("Error sending request to server: %v", err)
		}
		defer resp.Body.Close()

		// Check the response status code
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
		}

		// Read the response body
		responseBody := make([]byte, 12) // Length of "Test response"
		_, err = resp.Body.Read(responseBody)
		if err != nil {
			t.Fatalf("Error reading response body: %v", err)
		}

		// Check the response body content
		if string(responseBody) != "Test response" {
			t.Errorf("Expected response body 'Test response', but got '%s'", responseBody)
		}
	})
}

func createRedisClient(ctx context.Context, t *testing.T) *redis.Client {
	// Get the Redis server address from the environment variable
	redisAddress := os.Getenv("REDIS_ADDRESS")
	if redisAddress == "" {
		t.Fatal("REDIS_ADDRESS environment variable is not set")
	}

	// Create a Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})

	return rdb
}

func clearCache(ctx context.Context, rdb *redis.Client, t *testing.T) {
	// Clear the cache using FLUSHALL command
	err := rdb.FlushAll(ctx).Err()
	if err != nil {
		t.Fatalf("Error clearing cache: %v", err)
	}
}

func BenchmarkGetRepositories(b *testing.B) {
	api := &repositoriesCollectors.GitHubReposCollector{}
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{})
	route := mux.NewRouter()
	api.ConfigureCollector(ctx, rdb, route)

	req := httptest.NewRequest("GET", "/repositories?org=github&phrase=go", nil)
	w := httptest.NewRecorder()

	for i := 0; i < b.N; i++ {
		api.GetRepositories(w, req)
	}
}
