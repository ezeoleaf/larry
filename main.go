package main

import (
	"fmt"
	"log"
	"os"

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
				Value:       "golang",
				Usage:       "language for searching repos",
				Destination: &cfg.Language,
			},
		},
		Action: func(c *cli.Context) error {
			r := getRepo(cfg)
			fmt.Println(r)
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
