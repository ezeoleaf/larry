package main

import (
	"log"
	"os"
	"time"

	"github.com/ezeoleaf/GobotTweet/config"
	"github.com/ezeoleaf/GobotTweet/providers"
	"github.com/ezeoleaf/GobotTweet/providers/github"
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
				var provider providers.IContent
				if cfg.Provider == providers.Github {
					provider = github.NewGithubRepository(cfg)
				}

				if provider != nil {
					content := provider.GetContentToPublish()
					tweetContent(cfg, content)
				} else {
					log.Fatal("No valid provider! " + providers.GetValidProvidersToString())
				}

				time.Sleep(time.Duration(cfg.Periodicity) * time.Minute)
			}
		},
	}

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}
