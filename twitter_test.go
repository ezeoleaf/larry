package main

import (
	"testing"

	"github.com/ezeoleaf/GobotTweet/config"
)

func TestTweetRepoInSafeMode(t *testing.T) {
	c := config.Config{SafeMode: true}

	result := tweetContent(c, "Something to publish")

	if !result {
		t.Error("Expected: true, got false")
	}
}
