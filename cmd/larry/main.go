package main

import (
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/ezeoleaf/larry/blacklist"
	"github.com/ezeoleaf/larry/cache"
	"github.com/ezeoleaf/larry/config"
	"github.com/ezeoleaf/larry/larry"
	"github.com/ezeoleaf/larry/provider"
	"github.com/ezeoleaf/larry/provider/contentfile"
	"github.com/ezeoleaf/larry/provider/github"
	"github.com/ezeoleaf/larry/publisher"
	githubPub "github.com/ezeoleaf/larry/publisher/github"
	"github.com/ezeoleaf/larry/publisher/twitter"
	"github.com/go-redis/redis/v8"
	"github.com/urfave/cli/v2"
)

var (
	redisAddress = envString("REDIS_ADDRESS", "localhost:6379")

	githubAccessToken      = envString("GITHUB_ACCESS_TOKEN", "")
	githubPublishRepoOwner = envString("GITHUB_PUBLISH_REPO_OWNER", "")
	githubPublishRepoName  = envString("GITHUB_PUBLISH_REPO_NAME", "")
	githubPublishRepoFile  = envString("GITHUB_PUBLISH_REPO_FILE", "README.md")

	twitterConsumerKey    = envString("TWITTER_CONSUMER_KEY", "")
	twitterConsumerSecret = envString("TWITTER_CONSUMER_SECRET", "")
	twitterAccessToken    = envString("TWITTER_ACCESS_TOKEN", "")
	twitterAccessSecret   = envString("TWITTER_ACCESS_SECRET", "")
)

func main() {
	cfg := config.Config{}

	app := &cli.App{
		Name:  "Larry",
		Usage: "Bot that publishes information from providers to different publishers",
		Flags: larry.GetFlags(&cfg),
		Authors: []*cli.Author{
			{Name: "@ezeoleaf", Email: "ezeoleaf@gmail.com"},
			{Name: "@beesaferoot", Email: "hikenike6@gmail.com"},
			{Name: "@shubhcoder"},
			{Name: "@kannav02"},
			{Name: "@siddhant-k-code", Email: "siddhantkhare2694@gmail.com"},
			{Name: "@savagedev"},
		},
		Action: func(c *cli.Context) error {
			prov, err := getProvider(cfg)
			if err != nil {
				log.Fatal(err)
			}

			if prov == nil {
				log.Fatalf("could not initialize provider for %v", cfg.Provider)
			}

			pubs, err := getPublishers(cfg)
			if err != nil {
				log.Fatal(err)
			}

			if len(pubs) == 0 {
				log.Fatalln("no publishers initialized")
			}

			s := larry.Service{Provider: prov, Publishers: pubs}

			for {
				err := s.Run()
				if err != nil {
					log.Printf("Error in larry.Service.Run(): %v", err)
				}
				time.Sleep(time.Duration(cfg.Periodicity) * time.Minute)
			}
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatalln(err)
	}
}

func getProvider(cfg config.Config) (larry.Provider, error) {
	ro := &redis.Options{
		Addr:     redisAddress,
		Password: "", // no password set
		DB:       0,  // use default DB
	}

	cacheClient := cache.NewClient(ro)
	if err := blacklist.Load(cacheClient, cfg.BlacklistFile, cfg.GetCacheKeyPrefix()); err != nil {
		return nil, err
	}

	if cfg.Provider == provider.Github {
		np := github.NewProvider(githubAccessToken, cfg, cacheClient)
		return np, nil
	} else if cfg.Provider == provider.Contentfile {
		np, err := contentfile.NewProvider(cfg, cacheClient)
		return np, err
	}

	return nil, nil
}

func getPublishers(cfg config.Config) (map[string]larry.Publisher, error) {
	pubs := make(map[string]larry.Publisher)

	ps := strings.Split(cfg.Publishers, ",")

	for _, v := range ps {
		v = strings.ToLower(strings.TrimSpace(v))

		if _, ok := pubs[v]; ok {
			continue
		}

		if v == publisher.Twitter {
			accessKeys := twitter.AccessKeys{
				TwitterConsumerKey:    twitterConsumerKey,
				TwitterConsumerSecret: twitterConsumerSecret,
				TwitterAccessToken:    twitterAccessToken,
				TwitterAccessSecret:   twitterAccessSecret,
			}
			pubs[v] = twitter.NewPublisher(accessKeys, cfg)
		} else if v == publisher.Github {
			pubs[v] = githubPub.NewPublisher(githubAccessToken, cfg, githubPublishRepoOwner, githubPublishRepoName, githubPublishRepoFile)
		}
	}

	return pubs, nil
}

func envString(key string, fallback string) string {
	if value, ok := syscall.Getenv(key); ok {
		return value
	}
	return fallback
}
