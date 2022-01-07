package github

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/ezeoleaf/larry/cache"
	"github.com/ezeoleaf/larry/config"
	"github.com/ezeoleaf/larry/mock"
	"github.com/go-redis/redis/v8"
	"github.com/google/go-github/v39/github"
)

func TestQueryString(t *testing.T) {
	mr, _ := miniredis.Run()
	ro := &redis.Options{
		Addr: mr.Addr(),
	}

	rdb := cache.NewClient(ro)

	for _, tc := range []struct {
		Name        string
		mockConfig  config.Config
		returnValue string
	}{
		{
			Name: "Test get topic and language",
			mockConfig: config.Config{
				Language: "g",
				Topic:    "x",
			},
			returnValue: "a+topic:x+language:g",
		},
		{
			Name: "Test get topic",
			mockConfig: config.Config{
				Topic: "x",
			},
			returnValue: "a+topic:x",
		},
		{
			Name: "Test get language",
			mockConfig: config.Config{
				Language: "g",
			},
			returnValue: "a+language:g",
		},
		{
			Name:        "Test get nothing",
			mockConfig:  config.Config{},
			returnValue: "a+language:",
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			p := Provider{Config: tc.mockConfig, CacheClient: rdb}
			resp := p.getQueryString("a")

			if tc.returnValue != resp {
				t.Errorf("expected %v as value, got %v instead", tc.returnValue, resp)
			}
		})
	}
}

func TestIsRepoNotInRedis(t *testing.T) {
	for _, tc := range []struct {
		Name        string
		mockConfig  config.Config
		cacheClient mock.CacheClientMock
		returnValue bool
	}{
		{
			Name: "Test is not in cache",
			mockConfig: config.Config{
				Topic: "x",
			},
			cacheClient: mock.CacheClientMock{
				GetFn: func(key string) (string, error) {
					return "", redis.Nil
				},
				SetFn: func(key string, value interface{}, exp time.Duration) error {
					return nil
				},
			},
			returnValue: true,
		},
		{
			Name: "Test error in get",
			mockConfig: config.Config{
				Topic: "x",
			},
			cacheClient: mock.CacheClientMock{
				GetFn: func(key string) (string, error) {
					return "", nil
				},
				SetFn: func(key string, value interface{}, exp time.Duration) error {
					return nil
				},
			},
			returnValue: false,
		},
		{
			Name: "Test other error in get",
			mockConfig: config.Config{
				Topic: "x",
			},
			cacheClient: mock.CacheClientMock{
				GetFn: func(key string) (string, error) {
					return "", errors.New("some error")
				},
				SetFn: func(key string, value interface{}, exp time.Duration) error {
					return nil
				},
			},
			returnValue: false,
		},
		{
			Name: "Test error in set",
			mockConfig: config.Config{
				Topic: "x",
			},
			cacheClient: mock.CacheClientMock{
				GetFn: func(key string) (string, error) {
					return "", errors.New("some error")
				},
				SetFn: func(key string, value interface{}, exp time.Duration) error {
					return errors.New("could not save")
				},
			},
			returnValue: false,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			p := Provider{Config: tc.mockConfig, CacheClient: tc.cacheClient}
			resp := p.isRepoNotInCache(10)

			if tc.returnValue != resp {
				t.Errorf("expected %v as value, got %v instead", tc.returnValue, resp)
			}
		})
	}
}

func TestGetRepoUser(t *testing.T) {
	login := "t2"
	gu := github.User{Login: nil}
	guLog := github.User{Login: &login}
	for _, tc := range []struct {
		Name        string
		owner       *github.User
		username    *string
		userClient  mock.UserClientMock
		returnValue string
	}{
		{
			Name:        "Test empty owner",
			owner:       nil,
			returnValue: "",
		},
		{
			Name:        "Test empty owner login in owner",
			owner:       &gu,
			returnValue: "",
		},
		{
			Name: "Test error getting user",
			userClient: mock.UserClientMock{
				GetFn: func(ctx context.Context, user string) (*github.User, *github.Response, error) {
					return nil, nil, errors.New("some error")
				},
			},
			owner:       &guLog,
			returnValue: "",
		},
		{
			Name: "Test get twitter username not set",
			userClient: mock.UserClientMock{
				GetFn: func(ctx context.Context, user string) (*github.User, *github.Response, error) {
					return nil, nil, nil
				},
			},
			owner:       &guLog,
			returnValue: "",
		},
		{
			Name: "Test get twitter username set",
			userClient: mock.UserClientMock{
				GetFn: func(ctx context.Context, user string) (*github.User, *github.Response, error) {
					t := "twitterusername"
					u := github.User{TwitterUsername: &t}
					return &u, nil, nil
				},
			},
			owner:       &guLog,
			returnValue: "twitterusername",
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			p := Provider{GithubUserClient: tc.userClient}
			resp := p.getRepoUser(tc.owner)

			if tc.returnValue != resp {
				t.Errorf("expected %v as value, got %v instead", tc.returnValue, resp)
			}
		})
	}
}

func TestGetContent(t *testing.T) {
	login, lang, url, name, desc := "t2", "g", "url", "repo", "desc"
	count := 1
	gu := github.User{Login: &login}
	for _, tc := range []struct {
		Name        string
		repo        *github.Repository
		username    *string
		userClient  mock.UserClientMock
		mockConfig  config.Config
		returnValue string
	}{
		{
			Name:        "Test no repo no hashtags",
			mockConfig:  config.Config{},
			repo:        &github.Repository{Language: &lang, HTMLURL: &url},
			returnValue: "#g #github\nurl",
		},
		{
			Name:        "Test no repo with topic config for hashtags",
			mockConfig:  config.Config{Topic: "t"},
			repo:        &github.Repository{Language: &lang, HTMLURL: &url},
			returnValue: "#t #github\nurl",
		},
		{
			Name:        "Test no repo with language config for hashtags",
			mockConfig:  config.Config{Language: "l"},
			repo:        &github.Repository{Language: &lang, HTMLURL: &url},
			returnValue: "#l #github\nurl",
		},
		{
			Name:        "Test no repo with hashtags",
			mockConfig:  config.Config{Hashtags: "a,b,c"},
			repo:        &github.Repository{Language: &lang, HTMLURL: &url},
			returnValue: "#a #b #c\nurl",
		},
		{
			Name:        "Test with repo data and no hashtags",
			mockConfig:  config.Config{TweetLanguage: true},
			repo:        &github.Repository{Name: &name, Description: &desc, Language: &lang, HTMLURL: &url},
			returnValue: "repo: desc\nLang: g\n#g #github\nurl",
		},
		{
			Name:       "Test full with username",
			mockConfig: config.Config{TweetLanguage: true},
			userClient: mock.UserClientMock{
				GetFn: func(ctx context.Context, user string) (*github.User, *github.Response, error) {
					t := "twitterusername"
					u := github.User{TwitterUsername: &t}
					return &u, nil, nil
				},
			},
			repo:        &github.Repository{Name: &name, Description: &desc, Language: &lang, HTMLURL: &url, StargazersCount: &count, Owner: &gu},
			returnValue: "repo: desc\nLang: g\n⭐️ 1\n#g #github\nAuthor: @twitterusername\nurl",
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			p := Provider{GithubUserClient: tc.userClient,
				Config: tc.mockConfig}
			resp := p.getContent(tc.repo)

			if tc.returnValue != resp {
				t.Errorf("expected %v as value, got %v instead", tc.returnValue, resp)
			}
		})
	}
}

func TestGetSpecificRepo(t *testing.T) {
	var id1 int64
	var id2 int64
	id1, id2 = 1, 2
	repo1 := github.Repository{ID: &id1}
	repo2 := github.Repository{ID: &id2}
	for _, tc := range []struct {
		Name         string
		searchClient mock.SearchClientMock
		returnValue  *github.Repository
	}{
		{
			Name: "Test get no repo",
			searchClient: mock.SearchClientMock{
				RepositoriesFn: func(ctx context.Context, query string, opt *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error) {
					return nil, nil, errors.New("some error")
				},
			},
			returnValue: nil,
		},
		{
			Name: "Test get repo1 from list with two repos",
			searchClient: mock.SearchClientMock{
				RepositoriesFn: func(ctx context.Context, query string, opt *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error) {
					resp := github.RepositoriesSearchResult{Repositories: []*github.Repository{&repo1, &repo2}}
					ghResp := github.Response{LastPage: 1}
					return &resp, &ghResp, nil
				},
			},
			returnValue: &repo1,
		},
		{
			Name: "Test get repo2 from list with one repos",
			searchClient: mock.SearchClientMock{
				RepositoriesFn: func(ctx context.Context, query string, opt *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error) {
					resp := github.RepositoriesSearchResult{Repositories: []*github.Repository{&repo2}}
					ghResp := github.Response{LastPage: 1}
					return &resp, &ghResp, nil
				},
			},
			returnValue: &repo2,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			p := Provider{GithubSearchClient: tc.searchClient}
			resp := p.getSpecificRepo("a", 1)

			if tc.returnValue == nil && resp != nil {
				t.Errorf("expected no value, got %v instead", resp)
			} else if resp == nil && tc.returnValue != nil {
				t.Errorf("expected %v value, got no value instead", tc.returnValue)
			} else if resp != nil && tc.returnValue != nil {
				if &resp.ID != &tc.returnValue.ID {
					t.Errorf("expected %v value as repo ID, got %v instead", &tc.returnValue.ID, &resp.ID)
				}
			}
		})
	}
}

func TestGetRepo(t *testing.T) {
	var id1 int64
	var id2 int64
	id1, id2 = 1, 2
	archived, notArchived := true, false
	repo1 := github.Repository{ID: &id1, Archived: &notArchived}
	repo2 := github.Repository{ID: &id2, Archived: &archived}
	for _, tc := range []struct {
		Name         string
		searchClient mock.SearchClientMock
		cacheClient  mock.CacheClientMock
		returnValue  *github.Repository
		shouldError  bool
	}{
		{
			Name: "Test get no repo error in search",
			searchClient: mock.SearchClientMock{
				RepositoriesFn: func(ctx context.Context, query string, opt *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error) {
					return nil, nil, errors.New("some error")
				},
			},
			returnValue: nil,
			shouldError: true,
		},
		{
			Name: "Test get no error but no repos in search",
			searchClient: mock.SearchClientMock{
				RepositoriesFn: func(ctx context.Context, query string, opt *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error) {
					ghResp := github.Response{LastPage: 0}
					return nil, &ghResp, nil
				},
			},
			returnValue: nil,
			shouldError: true,
		},
		{
			Name: "Test get repo not in cache",
			searchClient: mock.SearchClientMock{
				RepositoriesFn: func(ctx context.Context, query string, opt *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error) {
					t := 1000
					resp := github.RepositoriesSearchResult{
						Repositories: []*github.Repository{&repo1, &repo2},
						Total:        &t,
					}
					ghResp := github.Response{LastPage: 1}
					return &resp, &ghResp, nil
				},
			},
			cacheClient: mock.CacheClientMock{
				GetFn: func(key string) (string, error) {
					return "", redis.Nil
				},
				SetFn: func(key string, value interface{}, exp time.Duration) error {
					return nil
				},
			},
			returnValue: &repo1,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			p := Provider{
				GithubSearchClient: tc.searchClient,
				CacheClient:        tc.cacheClient,
			}
			resp, err := p.getRepo()

			if tc.shouldError && err == nil {
				t.Error("expected error, got no error instead")
			}

			if tc.returnValue == nil && resp != nil {
				t.Errorf("expected no value, got %v instead", resp)
			} else if resp == nil && tc.returnValue != nil {
				t.Errorf("expected %v value, got no value instead", tc.returnValue)
			} else if resp != nil && tc.returnValue != nil {
				if &resp.ID != &tc.returnValue.ID {
					t.Errorf("expected %v value as repo ID, got %v instead", &tc.returnValue.ID, &resp.ID)
				}
			}
		})
	}
}
