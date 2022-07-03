package config

import (
	"fmt"
	"strings"
)

// Config is a struct that contains configuration for the app
type Config struct {
	Language      string
	Topic         string
	Periodicity   int
	Hashtags      string
	CacheSize     int
	TweetLanguage bool
	SafeMode      bool
	Provider      string
	Publishers    string
	BlacklistFile string
	ContentFile   string
	SkipCsvHeader bool
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

func (c *Config) GetCacheKeyPrefix() string {
	prefix := ""
	if c.Topic != "" {
		prefix = c.Topic + "-"
	}
	if c.Language != "" {
		prefix = prefix + c.Language + "-"
	}
	return prefix
}
