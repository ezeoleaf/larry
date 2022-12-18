package mastodon

import (
	"fmt"
	"testing"

	"github.com/ezeoleaf/larry/config"
	"github.com/ezeoleaf/larry/domain"
	"github.com/mattn/go-mastodon"
	"github.com/stretchr/testify/assert"
)

func TestNewPublisher(t *testing.T) {
	c := config.Config{SafeMode: true}
	pc := PublisherConfig{
		ClientCfg: &mastodon.Config{
			Server:       "foo",
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			AccessToken:  "123",
		},
	}

	p, err := NewPublisher(pc, c)
	if err != nil {
		t.Errorf("didn't expect an error")
	}
	if p.Client == nil {
		t.Error("expected new publisher, got nil")
	}
}

func TestPublishContentInSafeMode(t *testing.T) {
	c := config.Config{SafeMode: true}
	pc := PublisherConfig{
		ClientCfg: &mastodon.Config{
			Server:       "foo",
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			AccessToken:  "123",
		},
	}

	p, err := NewPublisher(pc, c)
	assert.NoError(t, err)

	ti, s, u := "ti", "s", "u"

	cont := domain.Content{Title: &ti, Subtitle: &s, URL: &u}

	r, err := p.PublishContent(&cont)
	assert.NoError(t, err, "expected no error got %v", err)
	assert.True(t, r, "expected content published in Safe Mode. No content published")
}

func TestCheckTootDataInSafeMode(t *testing.T) {
	c := config.Config{SafeMode: true}
	pc := PublisherConfig{
		ClientCfg: &mastodon.Config{
			Server:       "foo",
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			AccessToken:  "123",
		},
	}

	p, err := NewPublisher(pc, c)
	assert.NoError(t, err)

	ti, u := "Lorem Ipsum", "https://loremipsum.io/generator/?n=3&t=s"
	extraData := []string{"50k", "Author: @unknown"}

	for _, tc := range []struct {
		Name           string
		Subtitle       string
		ExpectedResult *mastodon.Toot
	}{
		{
			Name:     "Test should return same content",
			Subtitle: "t",
			ExpectedResult: &mastodon.Toot{
				Status:      fmt.Sprintf("%s\n%s\n\n50k\nAuthor: @unknown", ti, u),
				Visibility:  "public",
				Sensitive:   false,
				SpoilerText: "t",
			},
		},
		{
			Name:     "Test should truncate subtitle",
			Subtitle: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Vitae sapien pellentesque habitant morbi tristique senectus et netus et. Nunc sed velit dignissim sodales.",
			ExpectedResult: &mastodon.Toot{
				Status:      fmt.Sprintf("%s\n%s\n\n50k\nAuthor: @unknown", ti, u),
				Visibility:  mastodon.VisibilityPublic,
				Sensitive:   false,
				SpoilerText: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Vitae sapien pellentesque habitant morbi tristique senectus et netus et. Nunc sed velit dignissim sodales.",
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			cont := &domain.Content{Title: &ti, Subtitle: &tc.Subtitle, URL: &u, ExtraData: extraData}

			resp := p.prepareToot(cont)

			assert.Equal(t, tc.ExpectedResult.SpoilerText, resp.SpoilerText)
			assert.Equal(t, tc.ExpectedResult.Status, resp.Status)
			assert.Equal(t, tc.ExpectedResult.Visibility, resp.Visibility)
			assert.Equal(t, tc.ExpectedResult.Sensitive, resp.Sensitive)
			assert.False(t, resp.Sensitive)

			assert.Equal(t, tc.ExpectedResult, resp, "resp should be %v, got %v", tc.ExpectedResult, resp)
		})
	}

}
