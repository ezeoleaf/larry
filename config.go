package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Config is a struct that contains configuration for the app
type Config struct {
	Language      string
	Topic         string
	ConfigFile    string
	AccessCfg     AccessConfig
	Periodicity   int
	Hashtags      string
	CacheSize     int
	TweetLanguage bool
	SafeMode      bool
}

// AccessConfig is a struct that contains configuration for the clients
type AccessConfig struct {
	GithubAccessToken     string `json:"github_access_token"`
	TwitterConsumerKey    string `json:"twitter_consumer_key"`
	TwitterConsumerSecret string `json:"twitter_consumer_secret"`
	TwitterAccessToken    string `json:"twitter_access_token"`
	TwitterAccessSecret   string `json:"twitter_access_secret"`
	DevMode               bool   `json:"dev_mode"` //If DevMode is true then it wont post any tweet //TODO: Move to arg as safe mode
}

// SetConfigAccess reads a configuration file and unmarshall to struct
func (c *Config) SetConfigAccess() error {
	file, _ := ioutil.ReadFile(c.ConfigFile)

	return json.Unmarshal([]byte(file), &c.AccessCfg)
}

// GetGithubClient returns a client for using Github API from the config struct
func (a *AccessConfig) GetGithubClient() (context.Context, *github.Client) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: a.GithubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return ctx, client
}

// GetTwitterClient returns a client for using Twitter API from the config struct
func (a *AccessConfig) GetTwitterClient() *twitter.Client {
	config := oauth1.NewConfig(a.TwitterConsumerKey, a.TwitterConsumerSecret)
	token := oauth1.NewToken(a.TwitterAccessToken, a.TwitterAccessSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	return twitter.NewClient(httpClient)
}

// GetHashtags return a list of hashtags from a comma separated string
func (c *Config) GetHashtags() []string {

	if c.Hashtags == "" {
		return []string{}
	}

	hs := strings.Split(c.Hashtags, ",")

	for i, h := range hs {
		hs[i] = fmt.Sprintf("#%s", strings.TrimSpace(h))
	}

	return hs
}
