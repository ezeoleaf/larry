package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/google/go-github/github"
)

func tweetRepo(cfg Config, repo github.Repository) {
	if repo.HTMLURL == nil {
		return
	}

	hashtags, title, stargazers := "", "", ""

	if repo.Language != nil {
		hashtags += "#" + *repo.Language + " "
	}

	hashtags += "#github" + "\n"

	if repo.Name != nil {
		title += *repo.Name + ": "
	}

	if repo.Description != nil {
		title += *repo.Description + "\n"
	}

	if repo.StargazersCount != nil {
		stargazers += "⭐️: " + strconv.Itoa(*repo.StargazersCount) + "\n"
	}

	toTweet := title + stargazers + hashtags + *repo.HTMLURL

	client := cfg.AccessCfg.GetTwitterClient()
	_, _, err := client.Statuses.Update(toTweet, nil)

	if err != nil {
		log.Print(err)
	}
	fmt.Println("Tweet Published")
}
