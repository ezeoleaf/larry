package twitter

import (
	"testing"

	"github.com/ezeoleaf/larry/config"
	"github.com/ezeoleaf/larry/domain"
)

func TestNewPublisher(t *testing.T) {
	c := config.Config{SafeMode: true}
	ak := AccessKeys{}

	p := NewPublisher(ak, c)

	if p.Client == nil {
		t.Error("expected new publisher, got nil")
	}
}

func TestPublishContentInSafeMode(t *testing.T) {
	c := config.Config{SafeMode: true}
	ak := AccessKeys{}

	p := NewPublisher(ak, c)

	ti, s, u := "ti", "s", "u"

	cont := domain.Content{Title: &ti, Subtitle: &s, URL: &u}

	r, err := p.PublishContent(&cont)

	if !r {
		t.Error("expected content published in Safe Mode. No content published")
	}

	if err != nil {
		t.Errorf("expected no error got %v", err)
	}
}

func TestCheckTweetDataInSafeMode(t *testing.T) {
	c := config.Config{SafeMode: true}
	ak := AccessKeys{}

	p := NewPublisher(ak, c)

	subtitle := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
	subtitle += "Vitae sapien pellentesque habitant morbi tristique senectus et netus et. Nunc sed velit dignissim sodales."

	ti, s, u := "Lorem Ipsum", subtitle, "https://loremipsum.io/generator/?n=3&t=s"
	extraData := []string{"50k", "Author: @unknown"}

	cont := &domain.Content{Title: &ti, Subtitle: &s, URL: &u, ExtraData: extraData}

	resp := p.prepareTweet(cont)

	if len(resp) > TweetLength {
		t.Errorf("Tweet length is %v, which is greater than %v", len(resp), TweetLength)
	}
}
