package twitter

import (
	"fmt"
	"log"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/ezeoleaf/larry/config"
)

type Publisher struct {
	Client *twitter.Client
	Config config.Config
}

type AccessKeys struct {
	TwitterConsumerKey    string
	TwitterConsumerSecret string
	TwitterAccessToken    string
	TwitterAccessSecret   string
}

func NewPublisher(accessKeys AccessKeys, cfg config.Config) Publisher {
	log.Print("New Twitter Publisher")

	oauthCfg := oauth1.NewConfig(accessKeys.TwitterConsumerKey, accessKeys.TwitterConsumerSecret)
	oauthToken := oauth1.NewToken(accessKeys.TwitterAccessToken, accessKeys.TwitterAccessSecret)

	client := twitter.NewClient(oauthCfg.Client(oauth1.NoContext, oauthToken))

	p := Publisher{
		Config: cfg,
		Client: client,
	}

	return p
}

func (p Publisher) PublishContent(content string) (bool, error) {
	if p.Config.SafeMode {
		log.Print("Running in Safe Mode")
		log.Print(content)
		return true, nil
	}

	_, _, err := p.Client.Statuses.Update(content, nil)

	if err != nil {
		log.Print(err)
		return false, err
	}

	fmt.Println("Content Published")
	return true, nil
}
