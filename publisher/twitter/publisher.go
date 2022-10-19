package twitter

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/ezeoleaf/larry/config"
	"github.com/ezeoleaf/larry/domain"
	"golang.org/x/oauth2/clientcredentials"
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
}

// current limit of characters in tweet
const TweetLength int = 280

// NewPublisher returns a new publisher
func NewPublisher(accessKeys AccessKeys, cfg config.Config) Publisher {
	log.Print("New Twitter Publisher")

	config := &clientcredentials.Config{
		ClientID:     accessKeys.TwitterConsumerKey,
		ClientSecret: accessKeys.TwitterConsumerSecret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}

	client := twitter.NewClient(config.Client(context.Background()))

	p := Publisher{
		Config: cfg,
		Client: client,
	}

	return p
}

// prepareTweet converts a domain.Content in a string Tweet
func (p Publisher) prepareTweet(content *domain.Content) string {
	checkTweetData(content)

	tweet := fmt.Sprintf("%s: %s\n%s\n%s",
		*content.Title,
		*content.Subtitle,
		strings.Join(content.ExtraData, "\n"),
		*content.URL)

	return tweet
}

// changes description if generated tweet exceeds character limit
func checkTweetData(content *domain.Content) {
    titleLen := len(*content.Title)
    subTitleLen := len(*content.Subtitle)
    urlLen := len(*content.URL)
    extraDataLen := len(strings.Join(content.ExtraData, " "))

    size := titleLen + subTitleLen + urlLen + extraDataLen + 3  // '3' = extra space in string literal of tweet

    if size > TweetLength {
        truncateValue := subTitleLen - ((size - TweetLength) + 5)   // '5' = space for trailing "..."

        desiredDesc := *content.Subtitle
        desiredDesc = desiredDesc[:truncateValue]

        *content.Subtitle = strings.TrimSpace(desiredDesc) + " ..."
    }
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

