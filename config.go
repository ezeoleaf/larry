package main

// Config is a struct that contains configuration for the app
type Config struct {
	Language  string
	Topic     string
	AccessCfg AccessConfig
}

type AccessConfig struct {
	GithubAccessToken     string `json:"github_access_token"`
	TwitterConsumerKey    string `json:"twitter_consumer_key"`
	TwitterConsumerSecret string `json:"twitter_consumer_secret"`
	TwitterAccessToken    string `json:"twitter_access_token"`
	TwitterAccessSecret   string `json:"twitter_access_secret"`
}
