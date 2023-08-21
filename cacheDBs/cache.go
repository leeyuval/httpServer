package cacheDBs

import (
	"context"
)

// CacheDB defines the interface for cache databases.
type CacheDB interface {
	ExtractCachedData(ctx context.Context, cacheKey string, target interface{}) error
	StoreDataInCache(ctx context.Context, cacheKey string, data interface{}) error
}
