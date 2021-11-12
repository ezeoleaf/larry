package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// Client represent the repositories
type Client interface {
	Set(key string, value interface{}, exp time.Duration) error
	Get(key string) (string, error)
	Del(key string) error
}

// repository represent the repository model
type repository struct {
	Client *redis.Client
}

// NewRedisRepository will create an object that represent the Repository interface
func NewClient(ro *redis.Options) Client {
	return &repository{redis.NewClient(ro)}
}

// Set attaches the redis repository and set the data
func (r *repository) Set(key string, value interface{}, exp time.Duration) error {
	return r.Client.Set(ctx, key, value, exp).Err()
}

// Get attaches the redis repository and get the data
func (r *repository) Get(key string) (string, error) {
	get := r.Client.Get(ctx, key)
	return get.Result()
}

func (r *repository) Del(key string) error {
	return r.Client.Del(ctx, key).Err()
}
