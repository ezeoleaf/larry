package config

import (
	"testing"
)

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

func TestGetCacheKeyPrefix(t *testing.T) {
	for _, tc := range []struct {
		Name           string
		MockConfig     Config
		ExpectedPrefix string
	}{
		{
			Name:           "Should return empty",
			MockConfig:     Config{},
			ExpectedPrefix: "",
		},
		{
			Name:           "Should return only topic",
			MockConfig:     Config{Topic: "topic"},
			ExpectedPrefix: "topic-",
		},
		{
			Name:           "Should return only language",
			MockConfig:     Config{Language: "lang"},
			ExpectedPrefix: "lang-",
		},
		{
			Name:           "Should return topic and language",
			MockConfig:     Config{Topic: "topic", Language: "lang"},
			ExpectedPrefix: "topic-lang-",
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			prefix := tc.MockConfig.GetCacheKeyPrefix()

			if prefix != tc.ExpectedPrefix {
				t.Errorf("expected %s, got %s", tc.ExpectedPrefix, prefix)
			}
		})

	}

}
