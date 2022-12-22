package mastodon

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ezeoleaf/larry/config"
	"github.com/ezeoleaf/larry/domain"
	"github.com/mattn/go-mastodon"
)

// Publisher represents the publisher client
type Publisher struct {
	Client       *mastodon.Client
	Config       config.Config
	capabilities map[string]int
	ctx          context.Context
}

// AccessKeys represents the keys and tokens needed for comunication with the client
type PublisherConfig struct {
	ClientCfg *mastodon.Config
	Username  string
	Password  string
}

// NewPublisher returns a new publisher
func NewPublisher(pCfg PublisherConfig, cfg config.Config) (Publisher, error) {
	log.Print("New Mastodon Publisher")

	var err error

	ctx := context.Background()

	client := mastodon.NewClient(pCfg.ClientCfg)

	// authenticate with username and password when necessary
	if len(pCfg.ClientCfg.AccessToken) == 0 {
		err = client.Authenticate(ctx, pCfg.Username, pCfg.Password)
	}

	p := Publisher{
		Config: cfg,
		Client: client,
		ctx:    ctx,
		capabilities: map[string]int{
			// this seems to be a default, but it would be nice to discover it later
			// see: https://github.com/mattn/go-mastodon/pull/167
			"max_characters": 500,
		},
	}

	return p, err
}

func (p Publisher) prepareToot(content *domain.Content) *mastodon.Toot {
	var status, spoilerText string

	if len(*content.Title) > 0 {
		if len(*content.URL) > 0 {
			status = fmt.Sprintf("%s\n%s\n\n", *content.Title, *content.URL)
		} else {
			status = *content.Title
		}
	}

	if len(*content.Subtitle) > 0 {
		spoilerText = *content.Subtitle
	}

	if len(content.ExtraData) > 0 {
		status = status + strings.Join(content.ExtraData, "\n")
	}

	// using the '3' for ...
	if len(status) > (p.capabilities["max_characters"]) {
		status = status[0:(p.capabilities["max_characters"]-3)] + "..."
	}

	return &mastodon.Toot{
		Status:      status,
		Visibility:  mastodon.VisibilityPublic,
		Sensitive:   false,
		Language:    p.Config.Language,
		SpoilerText: spoilerText,
	}
}

func (p Publisher) PublishContent(content *domain.Content) (bool, error) {

	status := p.prepareToot(content)

	if p.Config.SafeMode {
		log.Print("Running in Safe Mode")
		log.Print(content)
		log.Print(status)
		return true, nil
	}

	_, err := p.Client.PostStatus(p.ctx, status)
	if err != nil {
		log.Print(err)
		return false, err
	}

	log.Println("Content Published")
	return true, nil
}
