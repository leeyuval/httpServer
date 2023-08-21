package cacheDBs

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

// RedisDB implements the CacheDB interface using Redis as the cache database.
type RedisDB struct {
	client *redis.Client
}

// NewRedisDB creates a new instance of RedisDB.
func NewRedisDB(client *redis.Client) *RedisDB {
	return &RedisDB{client: client}
}

func (r *RedisDB) ExtractCachedData(ctx context.Context, cacheKey string, target interface{}) error {
	data, err := r.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func (r *RedisDB) StoreDataInCache(ctx context.Context, cacheKey string, data interface{}) error {
	cachedDataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	result := r.client.Set(ctx, cacheKey, cachedDataBytes, time.Hour*12)
	if result.Err() == nil {
		log.Println("Data stored in cache successfully.")
	} else {
		log.Println(result.Err())
	}

	return result.Err()
}
