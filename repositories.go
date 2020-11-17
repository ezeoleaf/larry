package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var repositories *github.RepositoriesSearchResult

func getRepositories(cfg Config) *github.RepositoriesSearchResult {
	if repositories == nil {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: "1d54942d4646ccd4eb673d8eedf569e554ae4f20"},
		)
		tc := oauth2.NewClient(ctx, ts)

		client := github.NewClient(tc)

		var e error

		qs := fmt.Sprintf("language:%s", cfg.Language)

		repositories, _, e = client.Search.Repositories(ctx, qs, nil)

		if e != nil {
			panic(e)
		}
	}

	return repositories
}

func getRepo(config Config) github.Repository {
	repos := getRepositories(config)
	repo := repos.Repositories[0]
	return repo
}
