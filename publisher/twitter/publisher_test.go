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
