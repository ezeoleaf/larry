package github

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/ezeoleaf/GobotTweet/cache"
	"github.com/ezeoleaf/GobotTweet/config"
	"github.com/ezeoleaf/GobotTweet/providers"
	"github.com/go-redis/redis/v8"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var repositories *github.RepositoriesSearchResult
var rdb cache.Repository
var cfg config.Config
var client githubClient
var ctx context.Context

// repository represent the repository model
type githubProvider struct {
	Provider interface{}
}

type githubClient interface {
	Repositories(ctx context.Context, query string, opt *github.SearchOptions) (*github.RepositoriesSearchResult, *github.Response, error)
}

func NewGithubRepository(config config.Config) providers.IContent {
	cfg = config
	ro := &redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}

	rdb = cache.NewRedisRepository(ro)

	setClient()

	return &githubProvider{}
}

func (g *githubProvider) GetContentToPublish() string {
	r := getRepo()
	return getContent(r)
}

func getContent(repo *github.Repository) string {
	hashtags, title, stargazers := "", "", ""

	hs := cfg.GetHashtags()

	if len(hs) == 0 {
		if cfg.Topic != "" {
			hashtags += "#" + cfg.Topic + " "
		} else if cfg.Language != "" {
			hashtags += "#" + cfg.Language + " "
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

	if cfg.TweetLanguage {
		if repo.Language != nil {
			title += "Lang: " + *repo.Language + "\n"
		}
	}

	if repo.StargazersCount != nil {
		stargazers += "⭐️ " + strconv.Itoa(*repo.StargazersCount) + "\n"
	}

	return title + stargazers + hashtags + *repo.HTMLURL
}

func setClient() {
	ctx = context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.AccessCfg.GithubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client = github.NewClient(tc).Search

}

func getRepositories() ([]github.Repository, int) {
	if repositories == nil {
		var e error
		var qs string

		if cfg.Topic != "" && cfg.Language != ""{
			qs = fmt.Sprintf("topic:%s+language:%s", cfg.Topic, cfg.Language)
		}else if cfg.Topic != "" {
			qs = fmt.Sprintf("topic:%s", cfg.Topic)
		} else {
			qs = fmt.Sprintf("language:%s", cfg.Language)
		}

		so := github.SearchOptions{ListOptions: github.ListOptions{PerPage: 1}}

		repositories, _, e = client.Repositories(ctx, qs, &so)

		if e != nil {
			panic(e)
		}
	}

	return repositories.Repositories, *repositories.Total
}

func getSpecificRepo(pos int) *github.Repository {
	var e error
	var qs string

	if cfg.Topic != "" {
		qs = fmt.Sprintf("topic:%s", cfg.Topic)
	} else {
		qs = fmt.Sprintf("language:%s", cfg.Language)
	}

	so := github.SearchOptions{ListOptions: github.ListOptions{PerPage: 1, Page: pos}}

	repositories, _, e = client.Repositories(ctx, qs, &so)

	if e != nil {
		return nil
	}

	return &repositories.Repositories[0]
}

func getRepo() *github.Repository {
	_, total := getRepositories()

	var repo *github.Repository

	var found bool

	for !found {
		rand.Seed(time.Now().UTC().UnixNano())
		randPos := rand.Intn(total / 100)

		repo = getSpecificRepo(randPos)

		found = repo != nil && isRepoNotInRedis(repo, cfg.CacheSize*cfg.Periodicity, cfg.Topic)

		if found && *repo.Archived {
			found = false
			log.Print("Repository archived")
			log.Print(*repo.ID)
		}
	}

	return repo
}

func isRepoNotInRedis(r *github.Repository, t int, topic string) bool {
	k := topic + "-" + strconv.FormatInt(*r.ID, 10)
	_, err := rdb.Get(k)

	switch {
	case err == redis.Nil:
		err := rdb.Set(k, true, time.Duration(t)*time.Minute)
		if err != nil {
			panic(err)
		}

		return true
	case err != nil:
		fmt.Println("Get failed", err)
	}

	return false
}
