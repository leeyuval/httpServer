package repositoriesCollector_tests

import (
	"context"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"httpServer/cacheDBs"
	"testing"
)

func TestStoreAndExtractCachedData(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer mr.Close()

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	rdb := cacheDBs.NewRedisDB(client)

	ctx := context.Background()
	cacheKey := "testKey"
	expectedValue := "testValue"

	err = rdb.StoreDataInCache(ctx, cacheKey, expectedValue)
	if err != nil {
		t.Fatalf("failed to store data in cache: %v", err)
	}

	var actualValue string
	err = rdb.ExtractCachedData(ctx, cacheKey, &actualValue)
	if err != nil {
		t.Fatalf("failed to extract cached data: %v", err)
	}

	if expectedValue != actualValue {
		t.Errorf("expected value %q, got %q", expectedValue, actualValue)
	}
}
