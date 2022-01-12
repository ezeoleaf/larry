package twitter

import (
	"fmt"
	"log"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/ezeoleaf/larry/config"
	"github.com/ezeoleaf/larry/domain"
)

// Publisher represents the publisher client
type Publisher struct {
	Client *twitter.Client
	Config config.Config
}

// AccessKeys represents the keys and tokens needed for comunication with the client
type AccessKeys struct {
	TwitterConsumerKey    string
	TwitterConsumerSecret string
	TwitterAccessToken    string
	TwitterAccessSecret   string
}

// NewPublisher returns a new publisher
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

// prepareTweet convers a domain.Content in a string Tweet
func (p Publisher) prepareTweet(content *domain.Content) string {
	tweet := fmt.Sprintf("%s: %s\n%s\n%s",
		*content.Title,
		*content.Subtitle,
		strings.Join(content.ExtraData, "\n"),
		*content.URL)

	return tweet
}

// PublishContent receives a content to publish and try to publish
func (p Publisher) PublishContent(content *domain.Content) (bool, error) {
	tweet := p.prepareTweet(content)

	if p.Config.SafeMode {
		log.Print("Running in Safe Mode")
		log.Print(tweet)
		return true, nil
	}

	_, _, err := p.Client.Statuses.Update(tweet, nil)

	if err != nil {
		log.Print(err)
		return false, err
	}

	log.Println("Content Published")
	return true, nil
}
