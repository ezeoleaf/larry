package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/urfave/cli/v2"
)

func main() {
	config := oauth1.NewConfig("rvBUw3Jn2KDEA3YTV7AgqmL9f", "eVcymJznO9zv2C5qifJOSTcbqsuhX5At4qCX4LkvQ6CVE3kgkj")
	token := oauth1.NewToken("1328795920666931203-1btfihoHYnBdkyfvRV33Ahfr1ywUpB", "vs4Y6HcCxkOVASvzEwGOKgCXsy02rEEjWrMQNEy4wYse6")

	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	_ = twitter.NewClient(httpClient)

	// getRepo()

	cfg := Config{}

	app := &cli.App{
		Name:  "Tweet Random Repo",
		Usage: "Tweet random repositories",
		Flags: []cli.Flag{
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
