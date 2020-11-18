package main

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Config is a struct that contains configuration for the app
type Config struct {
	Language    string
	Topic       string
	ConfigFile  string
	AccessCfg   AccessConfig
	Periodicity int64
}

type AccessConfig struct {
	GithubAccessToken     string `json:"github_access_token"`
	TwitterConsumerKey    string `json:"twitter_consumer_key"`
	TwitterConsumerSecret string `json:"twitter_consumer_secret"`
	TwitterAccessToken    string `json:"twitter_access_token"`
	TwitterAccessSecret   string `json:"twitter_access_secret"`
}

func (c *Config) SetConfigAccess() {
	file, _ := ioutil.ReadFile(c.ConfigFile)

	_ = json.Unmarshal([]byte(file), &c.AccessCfg)
}

func (a *AccessConfig) GetGithubClient() (context.Context, *github.Client) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: a.GithubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return ctx, client
}

func (a *AccessConfig) GetTwitterClient() *twitter.Client {
	config := oauth1.NewConfig(a.TwitterConsumerKey, a.TwitterConsumerSecret)
	token := oauth1.NewToken(a.TwitterAccessToken, a.TwitterAccessSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	return twitter.NewClient(httpClient)
}
