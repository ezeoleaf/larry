package twitter

import (
	"testing"

	"github.com/ezeoleaf/larry/config"
)

func TestNewPublisher(t *testing.T) {
	c := config.Config{SafeMode: true}

	p := NewTwitterPublisher(c)

	if p == nil {
		t.Error("Expected new publisher, got nil")
	}
}

func TestPublishContentInSafeMode(t *testing.T) {
	c := config.Config{SafeMode: true}

	p := NewTwitterPublisher(c)

	r := p.PublishContent("Something to publish")

	if !r {
		t.Error("Expected content published in Safe Mode. No content published")
	}
}
