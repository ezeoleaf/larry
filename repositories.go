package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/go-github/github"
)

const cacheSize = 25

var repositories *github.RepositoriesSearchResult
var repositoriesCache []int64

func getRepositories(cfg Config) ([]github.Repository, int) {
	if repositories == nil {
		ctx, client := cfg.AccessCfg.GetGithubClient()

		var e error
		var qs string

		if cfg.Topic != "" {
			qs = fmt.Sprintf("topic:%s", cfg.Topic)
		} else {
			qs = fmt.Sprintf("language:%s", cfg.Language)
		}

		so := github.SearchOptions{ListOptions: github.ListOptions{PerPage: 1}}

		repositories, _, e = client.Search.Repositories(ctx, qs, &so)

		if e != nil {
			panic(e)
		}
	}

	return repositories.Repositories, *repositories.Total
}

func getSpecificRepo(cfg Config, pos int) *github.Repository {
	ctx, client := cfg.AccessCfg.GetGithubClient()

	var e error
	var qs string

	if cfg.Topic != "" {
		qs = fmt.Sprintf("topic:%s", cfg.Topic)
	} else {
		qs = fmt.Sprintf("language:%s", cfg.Language)
	}

	so := github.SearchOptions{ListOptions: github.ListOptions{PerPage: 1, Page: pos}}

	repositories, _, e = client.Search.Repositories(ctx, qs, &so)

	if e != nil {
		return nil
	}

	return &repositories.Repositories[0]
}

func getRepo(config Config) *github.Repository {
	_, total := getRepositories(config)

	var repo *github.Repository

	var found bool

	for !found {
		rand.Seed(time.Now().UTC().UnixNano())
		randPos := rand.Intn(total / 100)

		repo = getSpecificRepo(config, randPos)

		found = repo != nil && isRepoNotInCache(repo)
	}

	return repo
}

func isRepoNotInCache(r *github.Repository) bool {
	fmt.Println(*r.ID)
	for _, x := range repositoriesCache {
		if x == *r.ID {
			return false
		}
	}

	if len(repositoriesCache) == cacheSize {
		repositoriesCache = repositoriesCache[1:]
	}

	repositoriesCache = append(repositoriesCache, *r.ID)

	return true
}

func init() {
	repositoriesCache = []int64{}
}
