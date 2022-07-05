package twitter

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ezeoleaf/larry/config"
	"github.com/ezeoleaf/larry/domain"
	"github.com/ezeoleaf/larry/publisher/twitter/oauth2"
)

// Publisher represents the publisher client
type Publisher struct {
	Client *oauth2.Twitter
	Config config.Config
}

// AccessKeys represents the keys and tokens needed for comunication with the client
type AccessKeys struct {
	TwitterClientID     string
	TwitterClientSecret string
}

// current limit of characters in tweet
const TweetLength int = 280

// NewPublisher returns a new publisher
func NewPublisher(accessKeys AccessKeys, cfg config.Config) Publisher {
	log.Print("New Twitter Publisher")

	oauthCfg := oauth2.NewConfig(accessKeys.TwitterClientID, accessKeys.TwitterClientSecret)

	ctx := context.Background()
	client, err := oauth2.NewClient(ctx, oauthCfg)
	if err != nil {
		log.Fatal("client",err)
	}

	p := Publisher{
		Config: cfg,
		Client: client,
	}

	return p
}

// prepareTweet converts a domain.Content in a string Tweet
func (p Publisher) prepareTweet(content *domain.Content) string {
	tweet := checkTweetData(content)

	return tweet
}

// changes description if generated tweet exceeds character limit
func checkTweetData(content *domain.Content) string {
	titleLen := len(*content.Title)
	subTitleLen := len(*content.Subtitle)
	urlLen := len(*content.URL)
	extraDataLen := len(strings.Join(content.ExtraData, " "))

	size := titleLen + subTitleLen + urlLen + extraDataLen + 3 // '3' = extra space in string literal of tweet

	if size > TweetLength {
		truncateValue := subTitleLen - ((size - TweetLength) + 5) // '5' = space for trailing "..."

		desiredDesc := *content.Subtitle
		desiredDesc = desiredDesc[:truncateValue]

		*content.Subtitle = strings.TrimSpace(desiredDesc) + " ..."
	}

    return fmt.Sprintf("%s: %s\n%s\n%s",
		*content.Title,
		*content.Subtitle,
		strings.Join(content.ExtraData, "\n"),
		*content.URL)
}

// PublishContent receives a content to publish and try to publish
func (p Publisher) PublishContent(content *domain.Content) (bool, error) {
	tweet := p.prepareTweet(content)

	if p.Config.SafeMode {
		log.Print("Running in Safe Mode")
		log.Print(tweet)
		return true, nil
	}

	_, err := p.Client.Update(tweet)

	if err != nil {
		log.Print(err)
		return false, err
	}

	log.Println("Content Published")
	return true, nil
}
