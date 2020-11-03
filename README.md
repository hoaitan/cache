# Cache - A multi level cache

`Cache` implements mult cache levels:

- Local cache: By using [freecache](https://github.com/coocood/freecache).
- Redis cache: By using Redis.
- Multi level cache: combine of multi caches.
- `id` cache: Convention cache by caching `id` (string) by `name` (string).

## Installation

Init go module before installation:

```bash
go get github.com/hoaitan/cache
```

## Usage

```go
// Test data
type Data struct {
    Number int
}

// Local cache
localCache := local.New(local.Config{
    Enable:     true,
    Size:       10 * 1024 * 1024, // 10 MB
    DefaultTTL: 10,               // TTL in seconds
})

// Redis cache
redisCache := redis.New(redis.Config{
    Enable:     true,
    Endpoint:   "localhost:6379", // Redis endpoint
    Timeout:    60,               // Redis connection timeout
    DefaultTTL: 60,               // TTL in seconds
}, "my-service")

// Multi cache
multiCache := multi.New(
    localCache, // 1st cache layer
    redisCache, // 2nd cache layer
    // ... more cache layers but not recommend because of performance.
)

// Close these caches when don't use any more
defer multiCache.Close()

// Test data
sampleData := &Data{1}
cachedData := new(Data)

// localCache, redisCache or multiCache implement same Cache interface
// so they has same way to use

// Load missing cache
multiCache.Get("key1", cachedData, func() error {
    fmt.Println("missing cache")

    // Set cache when missing with:
    // ttl = -1: using cache's default TTL
    // ttl =  0: never expired
    // ttl >  0: expired in seconds
    return multiCache.Set("key1", sampleData, -1)
})
fmt.Printf("missing cache: cachedData=%+v\n", cachedData)

// Load cached data
multiCache.Get("key1", cachedData, func() error {
    fmt.Println("missing cache")
    return fmt.Errorf("should never happen")
})
fmt.Printf("with cache: cachedData=%+v\n", cachedData)

// Invalid cache
multiCache.Delete("key1")
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)