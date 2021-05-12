package main

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
	"github.com/google/go-github/github"
)

func TestIsRepoNotInRedis(t *testing.T) {
	mr, _ := miniredis.Run()
	ro := &redis.Options{
		Addr: mr.Addr(),
	}

	rdb = NewRedisRepository(ro)

	cacheSize := 1

	var id int64 = 1

	repo := github.Repository{ID: &id}

	result := isRepoNotInRedis(&repo, cacheSize, "t")

	if !result {
		t.Errorf("Expected not in cache, got %v", result)
	}

	id = 2
	repo = github.Repository{ID: &id}

	result = isRepoNotInRedis(&repo, cacheSize, "t")

	if !result {
		t.Errorf("Expected not in cache, got %v", result)
	}

	id = 2
	repo = github.Repository{ID: &id}

	result = isRepoNotInRedis(&repo, cacheSize, "t")

	if result {
		t.Errorf("Expected in cache, got %v", result)
	}

	rdb.Del("t-2")

	result = isRepoNotInRedis(&repo, cacheSize, "t")

	if !result {
		t.Errorf("Expected not in cache because of expiration, got %v", result)
	}
}

func TestIsRepoNotInRedisWithOtherTopic(t *testing.T) {
	mr, _ := miniredis.Run()
	ro := &redis.Options{
		Addr: mr.Addr(),
	}

	rdb = NewRedisRepository(ro)

	cacheSize := 1

	var id int64 = 1

	r := github.Repository{ID: &id}

	result := isRepoNotInRedis(&r, cacheSize, "t")

	if !result {
		t.Errorf("Expected not in cache, got %v", result)
	}

	id = 2
	r = github.Repository{ID: &id}

	result = isRepoNotInRedis(&r, cacheSize, "t")

	if !result {
		t.Errorf("Expected not in cache, got %v", result)
	}

	id = 2
	r = github.Repository{ID: &id}

	result = isRepoNotInRedis(&r, cacheSize, "t2")

	if !result {
		t.Errorf("Expected not in cache due to different topic, got %v", result)
	}
}
