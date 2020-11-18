package main

import (
	"fmt"
	"strconv"

	"github.com/google/go-github/github"
)

func tweetRepo(cfg Config, repo github.Repository) {

	hashtags := "#" + *repo.Language + " #github" + "\n"
	title := *repo.FullName + ": " + *repo.Description + "\n"
	stargazers := "⭐️: " + strconv.Itoa(*repo.StargazersCount) + "\n"

	toTweet := title + stargazers + hashtags + *repo.HTMLURL

	client := cfg.AccessCfg.GetTwitterClient()
	_, _, err := client.Statuses.Update(toTweet, nil)

	if err != nil {
		panic(err)
	}
	fmt.Println("Tweet Published")
}
