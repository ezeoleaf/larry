package main

import (
	"testing"
)

func TestGetFlags(t *testing.T) {
	mockConfig := Config{}
	flags := getFlags(&mockConfig)
	if flags == nil {
		t.Errorf("Expected flags, got %s", flags)
	}

	if len(flags) < 1 {
		t.Error("No flags found")
	}
}
