package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/ezeoleaf/larry/cache"
	"github.com/ezeoleaf/larry/config"
	"github.com/ezeoleaf/larry/larry"
	"github.com/ezeoleaf/larry/providers"
	"github.com/ezeoleaf/larry/providers/github"
	"github.com/ezeoleaf/larry/publishers"
	"github.com/ezeoleaf/larry/publishers/twitter"
	"github.com/go-redis/redis/v8"
	"github.com/urfave/cli/v2"
)

var (
	githubAccessToken     = envString("GITHUB_ACCESS_TOKEN", "")
	redisAddress          = envString("REDIS_ADDRESS", "localhost:6379")
	twitterConsumerKey    = envString("TWITTER_CONSUMER_KEY", "")
	twitterConsumerSecret = envString("TWITTER_CONSUMER_SECRET", "")
	twitterAccessToken    = envString("TWITTER_ACCESS_TOKEN", "")
	twitterAccessSecret   = envString("TWITTER_ACCESS_SECRET", "")
)

func main() {
	cfg := config.Config{}

	app := &cli.App{
		Name:    "Larry",
		Usage:   "Twitter bot that publishes random information from providers",
		Flags:   larry.GetFlags(&cfg),
		Authors: []*cli.Author{{Name: "Ezequiel Olea figueroa", Email: "ezeoleaf@gmail.com"}},
		Action: func(c *cli.Context) error {
			fmt.Println(cfg)
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

			s := larry.Service{Provider: prov, Publishers: pubs, Config: cfg}

			for {
				err := s.Run()
				if err != nil {
					return err
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
	if cfg.Provider == providers.Github {
		np := github.NewProvider(githubAccessToken, cfg, cacheClient)
		return np, nil
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

		if v == publishers.Twitter {
			accessKeys := twitter.AccessKeys{
				TwitterConsumerKey:    twitterConsumerKey,
				TwitterConsumerSecret: twitterConsumerSecret,
				TwitterAccessToken:    twitterAccessToken,
				TwitterAccessSecret:   twitterAccessSecret,
			}
			pubs[v] = twitter.NewPublisher(accessKeys, cfg)
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
