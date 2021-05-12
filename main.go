package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/urfave/cli/v2"
)

var rdb Repository
var ctx = context.Background()
var cfg Config

func init() {
	ro := &redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	}

	rdb = NewRedisRepository(ro)

	cfg = Config{}
}

func main() {
	app := &cli.App{
		Name:  "GobotTweet",
		Usage: "Twitter bot that tweets random repositories",
		Flags: getFlags(&cfg),
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
