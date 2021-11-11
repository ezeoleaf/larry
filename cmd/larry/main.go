package main

import (
	"log"
	"os"
	"syscall"

	"github.com/ezeoleaf/larry/cache"
	"github.com/ezeoleaf/larry/config"
	"github.com/ezeoleaf/larry/domain"
	"github.com/ezeoleaf/larry/larry"
	"github.com/ezeoleaf/larry/providers"
	"github.com/ezeoleaf/larry/providers/github"
	"github.com/go-redis/redis/v8"
	"github.com/urfave/cli/v2"
)

var (
	logFatalf = log.Fatalf
	logFatal  = log.Fatal
)

func init() {
	// cfg = config.Config{}
}

var (
	githubAccessToken = envString("GITHUB_ACCESS_TOKEN", "")
	redisAddress      = envString("REDIS_ADDRESS", "localhost:6379")
)

func main() {
	cfg := config.Config{}

	app := &cli.App{
		Name:    "Larry",
		Usage:   "Twitter bot that publishes random information from providers",
		Flags:   larry.GetFlags(cfg),
		Authors: []*cli.Author{{Name: "Ezequiel Olea figueroa", Email: "ezeoleaf@gmail.com"}},
		Action: func(c *cli.Context) error {
			prov, err := getProvider(cfg)
			if err != nil {
				panic(err)
			}
			s := larry.Service{Provider: prov}
			return nil
			// e := cfg.SetConfigAccess()
			// if e != nil {
			// 	panic(e)
			// }
			// for {
			// 	err := s.Run()
			// 	if err != nil {
			// 		return err
			// 	}
			// 	// provider, pubs := getProviderAndPublishers()

			// 	// content := provider.GetContentToPublish()

			// 	// for _, ps := range pubs {
			// 	// 	ps.PublishContent(content)
			// 	// }

			// 	time.Sleep(time.Duration(cfg.Periodicity) * time.Minute)
			// }
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		logFatal(err)
	}
}

func getProvider(cfg config.Config) (domain.Client, error) {
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

// func getPublishers(p string) map[string]publishers.IPublish {
// 	pubs := make(map[string]publishers.IPublish)

// 	ps := strings.Split(p, ",")

// 	for _, v := range ps {
// 		v = strings.TrimSpace(v)

// 		if _, ok := pubs[v]; ok {
// 			continue
// 		}

// 		if v == publishers.Twitter {
// 			pubs[v] = twitter.NewTwitterPublisher(cfg)
// 		}
// 	}

// 	return pubs
// }

func envString(key string, fallback string) string {
	if value, ok := syscall.Getenv(key); ok {
		return value
	}
	return fallback
}
