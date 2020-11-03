package redis

type Config struct {
	Enable     bool
	Endpoint   string
	Timeout    int // in seconds
	DefaultTTL int // in seconds
}
