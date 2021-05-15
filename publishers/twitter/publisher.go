package twitter

import (
	"fmt"
	"log"

	"github.com/ezeoleaf/GobotTweet/config"
	"github.com/ezeoleaf/GobotTweet/publishers"
)

var cfg config.Config

// twitterProvider represent the repository model
type twitterPublisher struct {
	Publisher interface{}
}

func NewTwitterPublisher(config config.Config) publishers.IPublish {
	log.Print("New Twitter Publisher")

	cfg = config

	return &twitterPublisher{}
}

func (g *twitterPublisher) PublishContent(c string) bool {
	if cfg.SafeMode {
		log.Print("Running in Safe Mode")
		log.Print(c)
		return true
	}

	client := cfg.AccessCfg.GetTwitterClient()
	_, _, err := client.Statuses.Update(c, nil)

	if err != nil {
		log.Print(err)
		return false
	}

	fmt.Println("Content Published")
	return true
}
