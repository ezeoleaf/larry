package main

import (
	"testing"

	"github.com/google/go-github/github"
)

func TestIsRepoNotInRedis(t *testing.T) {
	cacheSize := 1

	var id int64 = 1

	r := github.Repository{ID: &id}

	result := isRepoNotInRedis(&r, cacheSize)

	if !result {
		t.Errorf("Expected not in cache, got %v", result)
	}

	id = 2
	r = github.Repository{ID: &id}

	result = isRepoNotInRedis(&r, cacheSize)

	if !result {
		t.Errorf("Expected not in cache, got %v", result)
	}

	id = 2
	r = github.Repository{ID: &id}

	result = isRepoNotInRedis(&r, cacheSize)

	if result {
		t.Errorf("Expected in cache, got %v", result)
	}

	rdb.Del(ctx, "2")

	result = isRepoNotInRedis(&r, cacheSize)

	if !result {
		t.Errorf("Expected not in cache because of expiration, got %v", result)
	}
}
