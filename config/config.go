package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
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
	Provider      string
	Publishers    string
}

// AccessConfig is a struct that contains configuration for the clients
type AccessConfig struct {
	GithubAccessToken     string `json:"github_access_token"`
	TwitterConsumerKey    string `json:"twitter_consumer_key"`
	TwitterConsumerSecret string `json:"twitter_consumer_secret"`
	TwitterAccessToken    string `json:"twitter_access_token"`
	TwitterAccessSecret   string `json:"twitter_access_secret"`
}

// SetConfigAccess reads a configuration file and unmarshall to struct
func (c *Config) SetConfigAccess() error {
	file, _ := ioutil.ReadFile(c.ConfigFile)

	return json.Unmarshal([]byte(file), &c.AccessCfg)
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
