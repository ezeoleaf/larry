package main

import (
	"testing"
)

func TestSetConfigAccessNoFile(t *testing.T) {
	mockConfig := Config{}
	e := mockConfig.SetConfigAccess()
	if e == nil {
		t.Errorf("Expected error, got %s", e)
	}
}

func TestSetConfigAccessValidFile(t *testing.T) {
	mockConfig := Config{ConfigFile: "./config.example.json"}
	e := mockConfig.SetConfigAccess()
	if e != nil {
		t.Errorf("Expected nil, got %s", e)
	}
}

func TestGetHashtags(t *testing.T) {
	mockConfig := Config{Hashtags: "a,b,c "}
	expected := []string{"#a", "#b", "#c"}
	hs := mockConfig.GetHashtags()

	for i, h := range hs {
		if h != expected[i] {
			t.Errorf("Expected %s, got %s", expected[i], h)
		}
	}
}

func TestGetNoHashtags(t *testing.T) {
	mockConfig := Config{}
	hs := mockConfig.GetHashtags()

	if len(hs) > 0 {
		t.Errorf("Expected 0, got %v", len(hs))
	}
}
