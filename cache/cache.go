package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// Client represent the repositories
type Client interface {
	Set(key string, value interface{}, exp time.Duration) error
	Get(key string) (string, error)
	Del(key string) error
	Scan(key string, action func(context.Context, string) error) error
}

// repository represent the repository model
type repository struct {
	Client *redis.Client
}

// NewClient will create an object that represent the Repository interface
func NewClient(ro *redis.Options) Client {
	return &repository{redis.NewClient(ro)}
}

// Set attaches the redis repository and set the data
func (r *repository) Set(key string, value interface{}, exp time.Duration) error {
	return r.Client.Set(context.Background(), key, value, exp).Err()
}

// Get attaches the redis repository and get the data
func (r *repository) Get(key string) (string, error) {
	get := r.Client.Get(context.Background(), key)
	return get.Result()
}

func (r *repository) Del(key string) error {
	return r.Client.Del(context.Background(), key).Err()
}

func (r *repository) Scan(key string, action func(context.Context, string) error) error {
	ctx := context.Background()
	iter := r.Client.Scan(ctx, 0, key, 0).Iterator()
	for iter.Next(ctx) {
		if err := action(ctx, iter.Val()); err != nil {
			return err
		}
	}
	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}
