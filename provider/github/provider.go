package github

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/ezeoleaf/larry/cache"
	"github.com/ezeoleaf/larry/config"
	"github.com/go-redis/redis/v8"
	"github.com/google/go-github/v39/github"
	"golang.org/x/oauth2"
)

type searchClient interface {
	Repositories(ctx context.Context, query string, opt *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error)
}
type userClient interface {
	Get(ctx context.Context, user string) (*github.User, *github.Response, error)
}

// Provider represents the provider client
type Provider struct {
	GithubSearchClient searchClient
	GithubUserClient   userClient
	CacheClient        cache.Client
	Config             config.Config
}

// NewProvider returns a new provider client
func NewProvider(apiKey string, cfg config.Config, cacheClient cache.Client) Provider {
	log.Print("New Github Provider")
	p := Provider{Config: cfg, CacheClient: cacheClient}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: apiKey},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	p.GithubSearchClient = github.NewClient(tc).Search
	p.GithubUserClient = github.NewClient(tc).Users

	return p
}

// GetContentToPublish returns a string with the content to publish to be used by the publishers
func (p Provider) GetContentToPublish() (string, error) {
	r, err := p.getRepo()
	if err != nil {
		return "", err
	}
	return p.getContent(r), nil
}

func (p Provider) getRepositories() ([]*github.Repository, *int, error) {
	// TODO: Improve
	so := github.SearchOptions{ListOptions: github.ListOptions{PerPage: 1}}

	repositories, _, e := p.GithubSearchClient.Repositories(context.Background(), p.getQueryString(), &so)

	if e != nil {
		return nil, nil, e
	}

	return repositories.Repositories, repositories.Total, nil
}

func (p Provider) getRepo() (*github.Repository, error) {
	_, total, err := p.getRepositories()
	if err != nil {
		return nil, err
	}

	var repo *github.Repository

	var found bool

	for !found {
		rand.Seed(time.Now().UTC().UnixNano())
		randPos := rand.Intn(*total / 100)

		repo = p.getSpecificRepo(randPos)

		found = repo != nil && p.isRepoNotInCache(*repo.ID)

		if found && *repo.Archived {
			found = false
			log.Print("Repository archived")
			log.Print(*repo.ID)
		}
	}

	return repo, nil
}

func (p Provider) getQueryString() string {
	var qs string

	if p.Config.Topic != "" && p.Config.Language != "" {
		qs = fmt.Sprintf("topic:%s+language:%s", p.Config.Topic, p.Config.Language)
	} else if p.Config.Topic != "" {
		qs = fmt.Sprintf("topic:%s", p.Config.Topic)
	} else {
		qs = fmt.Sprintf("language:%s", p.Config.Language)
	}

	return qs
}

func (p Provider) getSpecificRepo(pos int) *github.Repository {
	so := github.SearchOptions{ListOptions: github.ListOptions{PerPage: 1, Page: pos}}

	repositories, _, e := p.GithubSearchClient.Repositories(context.Background(), p.getQueryString(), &so)

	if e != nil {
		return nil
	}

	return repositories.Repositories[0]
}

func (p Provider) isRepoNotInCache(repoID int64) bool {
	k := p.Config.Topic + "-" + strconv.FormatInt(repoID, 10)
	_, err := p.CacheClient.Get(k)

	switch {
	case err == redis.Nil:
		err := p.CacheClient.Set(k, true, time.Duration(p.Config.Periodicity)*time.Minute)
		if err != nil {
			return false
		}

		return true
	case err != nil:
		log.Println("Get failed", err)
		return false
	}

	return false
}

func (p Provider) getContent(repo *github.Repository) string {
	hashtags, title, stargazers, author := "", "", "", ""

	hs := p.Config.GetHashtags()

	if len(hs) == 0 {
		if p.Config.Topic != "" {
			hashtags += "#" + p.Config.Topic + " "
		} else if p.Config.Language != "" {
			hashtags += "#" + p.Config.Language + " "
		} else if repo.Language != nil {
			hashtags += "#" + *repo.Language + " "
		}

		hashtags += "#github" + "\n"
	} else {
		for _, h := range hs {
			if hashtags != "" {
				hashtags += " "
			}
			hashtags += h
		}
		hashtags += "\n"
	}

	if repo.Name != nil {
		title += *repo.Name + ": "
	}

	if repo.Description != nil {
		title += *repo.Description + "\n"
	}

	if p.Config.TweetLanguage {
		if repo.Language != nil {
			title += "Lang: " + *repo.Language + "\n"
		}
	}

	if repo.StargazersCount != nil {
		stargazers += "⭐️ " + strconv.Itoa(*repo.StargazersCount) + "\n"
	}

	owner := p.getRepoUser(repo.Owner)
	if owner != "" {
		author += "Author: @" + owner + "\n"
	}

	return title + stargazers + hashtags + author + *repo.HTMLURL
}

func (p Provider) getRepoUser(owner *github.User) string {
	if owner == nil || owner.Login == nil {
		return ""
	}

	gUser, _, err := p.GithubUserClient.Get(context.Background(), *owner.Login)

	if err != nil {
		return ""
	}

	return gUser.GetTwitterUsername()
}
