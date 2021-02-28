package main

import (
	"testing"

	"github.com/google/go-github/github"
)

func TestIsRepoNotInCache(t *testing.T) {

	cacheSize := 1
	repositoriesCache = []int64{}

	var id int64 = 1

	r := github.Repository{ID: &id}

	result := isRepoNotInCache(&r, cacheSize)

	if !result {
		t.Errorf("Expected not in cache, got %v", result)
	}

	id = 2
	r = github.Repository{ID: &id}

	result = isRepoNotInCache(&r, cacheSize)

	if !result {
		t.Errorf("Expected not in cache, got %v", result)
	}

	result = isRepoNotInCache(&r, cacheSize)

	if result {
		t.Errorf("Expected in cache, got %v", result)
	}
}
