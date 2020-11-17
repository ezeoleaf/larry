package main

import (
	"encoding/json"
	"io/ioutil"
)

// Config is a struct that contains configuration for the app
type Config struct {
	Language   string
	Topic      string
	ConfigFile string
	AccessCfg  AccessConfig
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
