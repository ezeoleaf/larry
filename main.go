package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/urfave/cli/v2"
)

var rdb *redis.Client
var ctx = context.Background()

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func main() {
	cfg := Config{}

	app := &cli.App{
		Name:  "GobotTweet",
		Usage: "Twitter bot that tweets random repositories",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "topic",
				Aliases:     []string{"t"},
				Value:       "",
				Usage:       "topic for searching repos",
				Destination: &cfg.Topic,
			},
			&cli.StringFlag{
				Name:        "lang",
				Aliases:     []string{"l"},
				Value:       "",
				Usage:       "language for searching repos",
				Destination: &cfg.Language,
			},
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Value:       "./config.json",
				Usage:       "path to config file",
				Destination: &cfg.ConfigFile,
			},
			&cli.IntFlag{
				Name:        "time",
				Aliases:     []string{"x"},
				Value:       15,
				Usage:       "periodicity of tweet in minutes",
				Destination: &cfg.Periodicity,
			},
			&cli.IntFlag{
				Name:        "cache",
				Aliases:     []string{"r"},
				Value:       50,
				Usage:       "size of cache for no repeating repositories",
				Destination: &cfg.CacheSize,
			},
			&cli.StringFlag{
				Name:        "hashtag",
				Aliases:     []string{"ht"},
				Value:       "",
				Usage:       "list of comma separated hashtags",
				Destination: &cfg.Hashtags,
			},
		},
		Action: func(c *cli.Context) error {
			e := cfg.SetConfigAccess()
			if e != nil {
				panic(e)
			}
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
