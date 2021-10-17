package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/ezeoleaf/GobotTweet/config"
	"github.com/ezeoleaf/GobotTweet/providers"
	"github.com/ezeoleaf/GobotTweet/providers/github"
	"github.com/ezeoleaf/GobotTweet/publishers"
	"github.com/ezeoleaf/GobotTweet/publishers/twitter"
	"github.com/urfave/cli/v2"
)

var (
	cfg       config.Config
	logFatalf = log.Fatalf
	logFatal  = log.Fatal
)

func init() {
	cfg = config.Config{}
}

func main() {
	app := &cli.App{
		Name:    "Larry",
		Usage:   "Twitter bot that tweets random repositories",
		Flags:   getFlags(&cfg),
		Authors: []*cli.Author{{Name: "Ezequiel Olea figueroa", Email: "ezeoleaf@gmail.com"}},
		Action: func(c *cli.Context) error {
			e := cfg.SetConfigAccess()
			if e != nil {
				panic(e)
			}
			for {
				provider, pubs := getProviderAndPublishers()

				content := provider.GetContentToPublish()

				for _, ps := range pubs {
					ps.PublishContent(content)
				}

				time.Sleep(time.Duration(cfg.Periodicity) * time.Minute)
			}
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		logFatal(err)
	}
}

func getProviderAndPublishers() (providers.IContent, map[string]publishers.IPublish) {
	if cfg.Provider == "" {
		logFatalf("%s is not a valid provider! %s", cfg.Provider, providers.GetValidProvidersToString())
	}

	if cfg.Publishers == "" {
		logFatalf("%s is not a valid publisher! %s", cfg.Provider, publishers.GetValidPublishersToString())
	}

	pr := getProvider(strings.ToLower(cfg.Provider))

	ps := getPublishers(strings.ToLower(cfg.Publishers))

	if pr == nil {
		logFatalf("%s is not a valid provider! %s", cfg.Provider, providers.GetValidProvidersToString())
	}

	if len(ps) == 0 {
		logFatalf("%s are not a valid publishers! %s", cfg.Publishers, publishers.GetValidPublishersToString())
	}

	return pr, ps
}

func getProvider(p string) providers.IContent {
	if p == providers.Github {
		return github.NewGithubRepository(cfg)
	}

	return nil
}

func getPublishers(p string) map[string]publishers.IPublish {
	pubs := make(map[string]publishers.IPublish)

	ps := strings.Split(p, ",")

	for _, v := range ps {
		v = strings.TrimSpace(v)

		if _, ok := pubs[v]; ok {
			continue
		}

		if v == publishers.Twitter {
			pubs[v] = twitter.NewTwitterPublisher(cfg)
		}
	}

	return pubs
}
