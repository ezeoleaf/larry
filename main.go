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
	"github.com/urfave/cli/v2"
)

var cfg config.Config

func init() {
	cfg = config.Config{}
}

func main() {
	app := &cli.App{
		Name:    "GobotTweet",
		Usage:   "Twitter bot that tweets random repositories",
		Flags:   getFlags(&cfg),
		Authors: []*cli.Author{{Name: "Ezequiel Olea figueroa", Email: "ezeoleaf@gmail.com"}},
		Action: func(c *cli.Context) error {
			e := cfg.SetConfigAccess()
			if e != nil {
				panic(e)
			}
			for {
				provider, _ := getProviderAndPublishers(&cfg)

				content := provider.GetContentToPublish()
				tweetContent(cfg, content)

				time.Sleep(time.Duration(cfg.Periodicity) * time.Minute)
			}
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func getProviderAndPublishers(c *config.Config) (providers.IContent, []publishers.IPublish) {
	if cfg.Provider == "" {
		log.Fatalf("%s is not a valid provider! %s", cfg.Provider, providers.GetValidProvidersToString())
	}

	if cfg.Publishers == "" {
		log.Fatalf("%s is not a valid publisher! %s", cfg.Provider, publishers.GetValidPublishersToString())
	}

	pr := getProvider(strings.ToLower(cfg.Provider))

	ps := getPublishers(strings.ToLower(cfg.Publishers))

	if pr == nil {
		log.Fatalf("%s is not a valid provider! %s", cfg.Provider, providers.GetValidProvidersToString())
	}

	if len(ps) == 0 {
		log.Fatalf("%s are not a valid publishers! %s", cfg.Publishers, publishers.GetValidPublishersToString())
	}

	return pr, ps
}

func getProvider(p string) providers.IContent {
	if cfg.Provider == providers.Github {
		return github.NewGithubRepository(cfg)
	}

	return nil
}

func getPublishers(p string) []publishers.IPublish {

	return []publishers.IPublish{}
}
