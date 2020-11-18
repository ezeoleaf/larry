package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {
	cfg := Config{}

	app := &cli.App{
		Name:  "Tweet Random Repo",
		Usage: "Tweet random repositories",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "topic",
				Value:       "golang",
				Usage:       "topic for searching repos",
				Destination: &cfg.Topic,
			},
			&cli.StringFlag{
				Name:        "lang",
				Value:       "",
				Usage:       "language for searching repos",
				Destination: &cfg.Language,
			},
			&cli.StringFlag{
				Name:        "cfg",
				Value:       "./config.json",
				Usage:       "path to config file",
				Destination: &cfg.ConfigFile,
			},
			&cli.Int64Flag{
				Name:        "time",
				Value:       15,
				Usage:       "periodicity of tweet",
				Destination: &cfg.Periodicity,
			},
		},
		Action: func(c *cli.Context) error {
			cfg.SetConfigAccess()
			for {
				r := getRepo(cfg)
				tweetRepo(cfg, r)
				time.Sleep(time.Duration(cfg.Periodicity) * time.Minute)
			}
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
