package github

import (
	"context"
	"errors"
	"testing"

	"github.com/alicebob/miniredis/v2"

	"github.com/ezeoleaf/GobotTweet/cache"
	"github.com/ezeoleaf/GobotTweet/config"
	"github.com/go-redis/redis/v8"
	"github.com/google/go-github/v39/github"
)

func TestIsRepoNotInRedis(t *testing.T) {
	mr, _ := miniredis.Run()
	ro := &redis.Options{
		Addr: mr.Addr(),
	}

	rdb = cache.NewRedisRepository(ro)

	cacheSize := 1

	var id int64 = 1

	repo := github.Repository{ID: &id}

	result := isRepoNotInRedis(&repo, cacheSize, "t")

	if !result {
		t.Errorf("Expected not in cache, got %v", result)
	}

	id = 2
	repo = github.Repository{ID: &id}

	result = isRepoNotInRedis(&repo, cacheSize, "t")

	if !result {
		t.Errorf("Expected not in cache, got %v", result)
	}

	id = 2
	repo = github.Repository{ID: &id}

	result = isRepoNotInRedis(&repo, cacheSize, "t")

	if result {
		t.Errorf("Expected in cache, got %v", result)
	}

	rdb.Del("t-2")

	result = isRepoNotInRedis(&repo, cacheSize, "t")

	if !result {
		t.Errorf("Expected not in cache because of expiration, got %v", result)
	}
}

func TestIsRepoNotInRedisWithOtherTopic(t *testing.T) {
	mr, _ := miniredis.Run()
	ro := &redis.Options{
		Addr: mr.Addr(),
	}

	rdb = cache.NewRedisRepository(ro)

	cacheSize := 1

	var id int64 = 1

	r := github.Repository{ID: &id}

	result := isRepoNotInRedis(&r, cacheSize, "t")

	if !result {
		t.Errorf("Expected not in cache, got %v", result)
	}

	id = 2
	r = github.Repository{ID: &id}

	result = isRepoNotInRedis(&r, cacheSize, "t")

	if !result {
		t.Errorf("Expected not in cache, got %v", result)
	}

	id = 2
	r = github.Repository{ID: &id}

	result = isRepoNotInRedis(&r, cacheSize, "t2")

	if !result {
		t.Errorf("Expected not in cache due to different topic, got %v", result)
	}
}

func TestGetContent(t *testing.T) {
	lang := "lang"
	name := "name"
	desc := "desc"
	star := 10
	url := "url"
	r := github.Repository{Language: &lang, Name: &name, Description: &desc, StargazersCount: &star, HTMLURL: &url}

	expected := "name: desc\n⭐️ 10\n#asd #github\nurl"

	cfg = config.Config{Topic: "asd"}

	result := getContent(&r)

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

	cfg = config.Config{Language: "lasd"}

	result := getContent(&r)

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

	cfg = config.Config{}

	result := getContent(&r)

	if expected != result {
		t.Errorf("Expected: %s, got %s", expected, result)
	}
}

func TestTweetRepoWithLangConfig(t *testing.T) {
	lang := "lang"
	name := "name"
	desc := "desc"
	star := 10
	url := "url"
	r := github.Repository{Language: &lang, Name: &name, Description: &desc, StargazersCount: &star, HTMLURL: &url}

	expected := "name: desc\nLang: lang\n⭐️ 10\n#lang #github\nurl"

	cfg = config.Config{TweetLanguage: true}

	result := getContent(&r)

	if expected != result {
		t.Errorf("Expected: %s, got %s", expected, result)
	}
}

func TestTweetRepoWithHashtags(t *testing.T) {
	lang := "lang"
	name := "name"
	desc := "desc"
	star := 10
	url := "url"
	r := github.Repository{Language: &lang, Name: &name, Description: &desc, StargazersCount: &star, HTMLURL: &url}

	expected := "name: desc\n⭐️ 10\n#a #b\nurl"

	cfg = config.Config{Hashtags: "a,b"}

	result := getContent(&r)

	if expected != result {
		t.Errorf("Expected: %s, got %s", expected, result)
	}
}

func TestTweetRepoWithAuthor(t *testing.T){
	lang := "lang"
	name := "name"
	url := "url"
	twitterUsername := "beesafe"
	author := github.User{TwitterUsername: &twitterUsername}
	r := github.Repository{Language: &lang, Name: &name, HTMLURL: &url, Owner: &author}
	
	expected := "name: #lang #github\nAuthor: @beesafe\nurl"
	
	cfg = config.Config{}

	result := getContent(&r)

	if expected != result {
		t.Errorf("Expected: %s, got %s", expected, result)
	}
}

type MockClient struct {
	RepositoriesFunc func(ctx context.Context, query string, opt *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error)
}

var GetRepositoriesFunc func(ctx context.Context, query string, opt *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error)

// Do is the mock client's `Do` func
func (m *MockClient) Repositories(ctx context.Context, query string, opt *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error) {
	return GetRepositoriesFunc(ctx, query, opt)
}

func TestGetRepositories(t *testing.T) {
	GetRepositoriesFunc = func(ctx context.Context, query string, opt *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error) {
		i := 1
		var repId int64 = 1
		reps := []*github.Repository{
			{ID: &repId},
		}
		r := github.RepositoriesSearchResult{Total: &i, Repositories: reps}
		return &r, nil, nil
	}

	client = &MockClient{}
	repos, total := getRepositories()

	if total != 1 {
		t.Errorf("Expected total one repository, got %v", total)
	}

	if len(repos) != 1 {
		t.Errorf("Expected one repo in slice, got %v", len(repos))
	}
}

func TestGetSpecificRepo(t *testing.T) {
	GetRepositoriesFunc = func(ctx context.Context, query string, opt *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error) {
		i := 1
		var repId int64 = 1
		reps := []*github.Repository{
			{ID: &repId},
		}
		r := github.RepositoriesSearchResult{Total: &i, Repositories: reps}
		return &r, nil, nil
	}

	client = &MockClient{}
	repo := getSpecificRepo(1)

	if repo == nil {
		t.Error("Expected repository got nil")
	}

	cfg.Topic = "topic"
	repo = getSpecificRepo(1)
	if repo == nil {
		t.Error("Expected repository got nil")
	}
}

func TestGetSpecificRepoError(t *testing.T) {
	GetRepositoriesFunc = func(ctx context.Context, query string, opt *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error) {
		return nil, nil, errors.New("Error")
	}

	client = &MockClient{}
	repo := getSpecificRepo(1)

	if repo != nil {
		t.Error("Expected nil got repository")
	}
}
