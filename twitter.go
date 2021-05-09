package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/google/go-github/github"
)

func getTweet(cfg Config, repo *github.Repository) string {
	hashtags, title, stargazers := "", "", ""

	hs := cfg.GetHashtags()

	if len(hs) == 0 {
		if cfg.Topic != "" {
			hashtags += "#" + cfg.Topic + " "
		} else if cfg.Language != "" {
			hashtags += "#" + cfg.Language + " "
		} else if repo.Language != nil {
			hashtags += "#" + *repo.Language + " "
		}

		hashtags += "#github" + "\n"
	} else {
		for _, h := range hs {
			if hashtags != "" {
				hashtags += " "
			}
			hashtags += h
		}
		hashtags += "\n"
	}

	if repo.Name != nil {
		title += *repo.Name + ": "
	}

	if repo.Description != nil {
		title += *repo.Description + "\n"
	}

	if cfg.TweetLanguage {
		if repo.Language != nil {
			title += "Lang: " + *repo.Language + "\n"
		}
	}

	if repo.StargazersCount != nil {
		stargazers += "⭐️ " + strconv.Itoa(*repo.StargazersCount) + "\n"
	}

	return title + stargazers + hashtags + *repo.HTMLURL
}

func tweetRepo(cfg Config, repo *github.Repository) bool {
	if repo.HTMLURL == nil {
		return false
	}

	tweet := getTweet(cfg, repo)

	if cfg.SafeMode {
		log.Print("Running in Safe Mode")
		log.Print(tweet)
		return true
	}

	client := cfg.AccessCfg.GetTwitterClient()
	_, _, err := client.Statuses.Update(tweet, nil)

	if err != nil {
		log.Print(err)
		return false
	}

	fmt.Println("Tweet Published")
	return true
}
