package main

import (
	"testing"

	"github.com/google/go-github/github"
)

func TestGetTweet(t *testing.T) {
	lang := "lang"
	name := "name"
	desc := "desc"
	star := 10
	url := "url"
	r := github.Repository{Language: &lang, Name: &name, Description: &desc, StargazersCount: &star, HTMLURL: &url}

	expected := "name: desc\n⭐️ 10\n#asd #github\nurl"

	mockConfig := Config{Topic: "asd"}

	result := getTweet(mockConfig, r)

	if expected != result {
		t.Errorf("Expected: %s, got %s", expected, result)
	}
}

func TestGetTweetNoTopic(t *testing.T) {
	lang := "lang"
	name := "name"
	desc := "desc"
	star := 10
	url := "url"
	r := github.Repository{Language: &lang, Name: &name, Description: &desc, StargazersCount: &star, HTMLURL: &url}

	expected := "name: desc\n⭐️ 10\n#lasd #github\nurl"

	mockConfig := Config{Language: "lasd"}

	result := getTweet(mockConfig, r)

	if expected != result {
		t.Errorf("Expected: %s, got %s", expected, result)
	}
}

func TestGetTweetNoTopicNoLang(t *testing.T) {
	lang := "lang"
	name := "name"
	desc := "desc"
	star := 10
	url := "url"
	r := github.Repository{Language: &lang, Name: &name, Description: &desc, StargazersCount: &star, HTMLURL: &url}

	expected := "name: desc\n⭐️ 10\n#lang #github\nurl"

	mockConfig := Config{}

	result := getTweet(mockConfig, r)

	if expected != result {
		t.Errorf("Expected: %s, got %s", expected, result)
	}
}

func TestTweetRepoWithoutURL(t *testing.T) {
	c := Config{}

	lang := "lang"
	name := "name"
	desc := "desc"
	star := 10

	r := github.Repository{Language: &lang, Name: &name, Description: &desc, StargazersCount: &star}

	result := tweetRepo(c, r)

	if result {
		t.Error("Expected: false, got true")
	}
}
