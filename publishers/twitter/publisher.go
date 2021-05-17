package twitter

import (
	"fmt"
	"log"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/ezeoleaf/GobotTweet/config"
	"github.com/ezeoleaf/GobotTweet/publishers"
)

var cfg config.Config
var client *twitter.Client

// twitterProvider represent the repository model
type twitterPublisher struct {
	Publisher interface{}
}

func NewTwitterPublisher(config config.Config) publishers.IPublish {
	log.Print("New Twitter Publisher")

	cfg = config

	setClient()

	return &twitterPublisher{}
}

func setClient() {
	oauthCfg := oauth1.NewConfig(cfg.AccessCfg.TwitterConsumerKey, cfg.AccessCfg.TwitterConsumerSecret)
	oauthToken := oauth1.NewToken(cfg.AccessCfg.TwitterAccessToken, cfg.AccessCfg.TwitterAccessSecret)

	client = twitter.NewClient(oauthCfg.Client(oauth1.NoContext, oauthToken))
}

func (g *twitterPublisher) PublishContent(c string) bool {
	if cfg.SafeMode {
		log.Print("Running in Safe Mode")
		log.Print(c)
		return true
	}

	_, _, err := client.Statuses.Update(c, nil)

	if err != nil {
		log.Print(err)
		return false
	}

	fmt.Println("Content Published")
	return true
}
