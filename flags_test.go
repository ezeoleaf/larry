package main

import (
	"testing"

	"github.com/ezeoleaf/larry/config"
)

func TestGetFlags(t *testing.T) {
	mockConfig := config.Config{}
	flags := getFlags(&mockConfig)
	if flags == nil {
		t.Errorf("Expected flags, got %s", flags)
	}

	if len(flags) < 1 {
		t.Error("No flags found")
	}
}
