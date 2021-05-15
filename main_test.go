package main

import (
	"fmt"
	"testing"
)

func TestGetProviderAndPublishers(t *testing.T) {

	cfg.Provider = "github"
	cfg.Publishers = "twitter"

	pr, ps := getProviderAndPublishers()

	if pr == nil {
		t.Error("Expected github provider but no provider received")
	}

	if len(ps) == 0 {
		t.Error("Expected at least one publisher but no publisher received ")
	}
}

func TestGetEmptyProviderAndEmptyPublishers(t *testing.T) {

	origLogFatalf := logFatalf

	// After this test, replace the original fatal function
	defer func() { logFatalf = origLogFatalf }()

	errors := []string{}

	logFatalf = func(format string, args ...interface{}) {
		fmt.Println("alksndlaksndlaksndlaksndalksdnalksnd")
		if len(args) > 0 {
			errors = append(errors, fmt.Sprintf(format, args))
		} else {
			errors = append(errors, format)
		}
	}

	cfg.Provider = ""
	cfg.Publishers = ""

	pr, ps := getProviderAndPublishers()

	if len(errors) < 4 {
		t.Errorf("excepted four error, actual %v", len(errors))
	}

	if pr != nil {
		t.Error("Expected no provider but provider received")
	}

	if len(ps) > 0 {
		t.Error("Expected no publisher but publisher received")
	}
}

func TestGetMultiplePublishers(t *testing.T) {
	ps := getPublishers("twitter, github, twitter")

	if len(ps) != 1 {
		t.Errorf("Expected only one publisher. Got %v", len(ps))
	}
}
