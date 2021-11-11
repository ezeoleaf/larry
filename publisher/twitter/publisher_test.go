package twitter

import (
	"testing"

	"github.com/ezeoleaf/larry/config"
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

	r, err := p.PublishContent("Something to publish")

	if !r {
		t.Error("expected content published in Safe Mode. No content published")
	}

	if err != nil {
		t.Errorf("expected no error got %v", err)
	}
}
