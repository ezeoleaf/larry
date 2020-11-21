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
