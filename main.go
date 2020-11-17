package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	// var s string
	// var cfg Config
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
		},
		Action: func(c *cli.Context) error {
			cfg.SetConfigAccess()
			r := getRepo(cfg)
			fmt.Println(*r.URL)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
