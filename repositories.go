package main

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var repositories *github.RepositoriesSearchResult

func getRepositories(cfg Config) ([]github.Repository, int) {
	if repositories == nil {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: cfg.AccessCfg.GithubAccessToken},
		)
		tc := oauth2.NewClient(ctx, ts)

		client := github.NewClient(tc)

		var e error
		var qs string

		if cfg.Topic != "" {
			qs = fmt.Sprintf("topic:%s", cfg.Topic)
		} else {
			qs = fmt.Sprintf("language:%s", cfg.Language)
		}

		repositories, _, e = client.Search.Repositories(ctx, qs, nil)

		if e != nil {
			panic(e)
		}
	}

	return repositories.Repositories, *repositories.Total
}

func getRepo(config Config) github.Repository {
	repos, total := getRepositories(config)
	fmt.Println(len(repos))
	randPos := rand.Intn(total)
	repo := repos[randPos]
	return repo
}
