package main

import (
	"fmt"
	"strconv"

	"github.com/google/go-github/github"
)

func tweetRepo(cfg Config, repo github.Repository) {
	toTweet := *repo.FullName + ": " + *repo.Description + "\n⭐️: " + strconv.Itoa(*repo.StargazersCount) + "\n " + *repo.HTMLURL

	client := cfg.AccessCfg.GetTwitterClient()
	_, _, err := client.Statuses.Update(toTweet, nil)

	if err != nil {
		panic(err)
	}
	fmt.Println("Tweet Published")
}
