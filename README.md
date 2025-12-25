# Cache - A multi-level cache library for Go

A flexible and efficient caching library that supports multiple cache layers with a unified interface.

## Features

- **Local Cache**: In-memory caching using [freecache](https://github.com/coocood/freecache) with gob encoding
- **Redis Cache**: Distributed caching using [go-redis/redis/v8](https://github.com/go-redis/redis) with JSON encoding
- **Multi-level Cache**: Combine multiple cache layers for optimal performance
- **ID Cache**: Convention-based caching for name-to-ID mapping with custom load functions
- **Flexible TTL**: Support for default, infinite, and custom TTL configurations
- **Health Checks**: Built-in readiness and enablement checks
- **Namespace Support**: Key prefixing for Redis and utility functions for key generation

## Installation

```bash
go get github.com/hoaitan/cache
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/hoaitan/cache/local"
    "github.com/hoaitan/cache/redis"
    "github.com/hoaitan/cache/multi"
)

type Data struct {
    Number int
}

func main() {
    // Create local cache (10 MB, 10s default TTL)
    localCache := local.New(local.Config{
        Enable:     true,
        Size:       10 * 1024 * 1024, // 10 MB (minimum 512 KB)
        DefaultTTL: 10,               // TTL in seconds
    })

    // Create Redis cache with key prefix
    redisCache := redis.New(redis.Config{
        Enable:     true,
        Endpoint:   "localhost:6379", // Redis endpoint
        Timeout:    60,               // Connection timeout in seconds
        DefaultTTL: 60,               // TTL in seconds
    }, "my-service") // Key prefix for namespacing

    // Create multi-level cache (local -> redis)
    multiCache := multi.New(
        localCache, // 1st layer (fastest)
        redisCache, // 2nd layer (shared)
        // Add more layers if needed (not recommended for performance)
    )

    // Always close caches when done
    defer multiCache.Close()

    // Prepare test data
    sampleData := &Data{Number: 1}
    cachedData := new(Data)

    // Get with cache miss callback
    multiCache.Get("key1", cachedData, func() error {
        fmt.Println("Cache miss - loading data...")

        // Set cache on miss with:
        // ttl = -1: use cache's default TTL
        // ttl =  0: never expire
        // ttl >  0: expire in N seconds
        return multiCache.Set("key1", sampleData, -1)
    })
    fmt.Printf("First call: %+v\n", cachedData)

    // Get cached data (no callback needed)
    multiCache.Get("key1", cachedData, nil)
    fmt.Printf("Second call (cached): %+v\n", cachedData)

    // Delete from all cache levels
    multiCache.Delete("key1")
}
```

## API Reference

All cache implementations (local, redis, multi) implement the `Cache` interface:

### Core Methods

```go
// Set data with TTL (in seconds)
// ttl = -1: use default TTL from config
// ttl =  0: never expire
// ttl >  0: expire in N seconds
Set(key string, data interface{}, ttl int) error

// Get cached data with optional miss callback
// If cache miss and fn is provided, fn() will be called
Get(key string, ptr interface{}, fn MissCacheFn) error

// Delete cache entry
// Returns true if key was deleted, false otherwise
Delete(key string) (bool, error)

// Check if key exists in cache
IsExist(key string) (bool, error)

// Flush all cache entries
// Returns number of entries flushed
// Note: Not supported for Redis cache (returns NotSupportedErr)
Flush() (int, error)

// Check if cache backend is ready
// For Redis: performs ping check
// For Local: always returns true
IsReady() bool

// Check if cache is enabled
IsEnable() bool

// Close cache connection
Close() error
```

### Utility Functions

```go
// Create namespaced cache key from parts
// Example: cache.MakeKey("user", "123") -> "user:123"
func MakeKey(parts ...string) string
```

## Advanced Usage

### ID Cache

The ID cache provides a specialized interface for caching name-to-ID mappings with custom load functions:

```go
import (
    "github.com/hoaitan/cache/id"
    "github.com/hoaitan/cache/local"
)

// Create underlying cache
baseCache := local.New(local.Config{
    Enable:     true,
    Size:       5 * 1024 * 1024,
    DefaultTTL: 3600,
})

// Create ID cache with namespace
idCache := id.New(baseCache, "user-ids")

// Set custom load function
idCache.SetLoadFn(func(name string) (string, error) {
    // Load ID from database or external service
    return loadUserIDFromDB(name)
})

// Get or set ID (automatically loads if missing)
userID, err := idCache.GetOrSet("john.doe")
if err != nil {
    // Handle error
}

// Check if name exists in cache
exists, err := idCache.IsExist("john.doe")

// Delete cached entry
err = idCache.Delete("john.doe")
```

**ID Cache Features:**
- Default TTL of 24 hours
- Automatic loading via custom function
- Namespace support for key organization
- Built on top of standard Cache interface

### Health Checks

```go
// Check if cache backend is ready
if !multiCache.IsReady() {
    log.Fatal("Cache is not ready")
}

// Check if cache is enabled
if multiCache.IsEnable() {
    // Use cache
} else {
    // Skip cache, load directly
}
```

### Key Management

```go
import "github.com/hoaitan/cache"

// Create namespaced keys
userKey := cache.MakeKey("users", "123")           // "users:123"
sessionKey := cache.MakeKey("session", "abc", "v2") // "session:abc:v2"

// Check existence
exists, err := redisCache.IsExist(userKey)

// Delete specific keys
deleted, err := redisCache.Delete(userKey)
```

### Flushing Cache

```go
// Flush local cache (clears all entries)
count, err := localCache.Flush()
fmt.Printf("Flushed %d entries\n", count)

// Note: Redis cache does not support Flush() for performance reasons
count, err := redisCache.Flush()
// Returns cache.NotSupportedErr
```

### Error Handling

```go
import "github.com/hoaitan/cache"

// Handle cache operations
err := multiCache.Set("key", data, 60)
if err != nil {
    // Handle encoding or connection errors
}

// Handle unsupported operations
count, err := redisCache.Flush()
if err == cache.NotSupportedErr {
    log.Println("Flush not supported for Redis cache")
}
```

## Configuration Details

### Local Cache Config

```go
type Config struct {
    Enable     bool // Enable/disable cache
    Size       int  // Cache size in bytes (minimum 512 KB)
    DefaultTTL int  // Default TTL in seconds
}
```

### Redis Cache Config

```go
type Config struct {
    Enable     bool   // Enable/disable cache
    Endpoint   string // Redis server address (host:port)
    Timeout    int    // Dial/Read/Write timeout in seconds
    DefaultTTL int    // Default TTL in seconds
}
```

## Multi-Level Cache Behavior

When using multi-level cache:

- **Get**: Returns data from the first cache layer that contains the key
- **Set**: Writes data to all cache layers
- **Delete**: Removes key from all cache layers
- **IsExist**: Returns true if key exists in any layer
- **Flush**: Flushes all layers that support it
- **IsReady**: Returns true only if all layers are ready
- **IsEnable**: Returns true if at least one layer is enabled
- **Close**: Closes all cache layers

## Best Practices

1. **Layer Order**: Place faster caches first (e.g., local before Redis)
2. **TTL Strategy**: Use shorter TTL for local cache, longer for Redis
3. **Error Handling**: Always check errors, especially for network-based caches
4. **Resource Cleanup**: Always defer `Close()` to prevent resource leaks
5. **Namespace Keys**: Use `MakeKey()` or Redis key prefix to avoid collisions
6. **Health Checks**: Use `IsReady()` before critical operations
7. **Performance**: Limit multi-level cache to 2-3 layers maximum

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)