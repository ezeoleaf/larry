package oauth2

import (
	"testing"

	oauth "golang.org/x/oauth2"
)

func TestNewConfig(t *testing.T) {
	for _, test := range []struct {
		name     string
		id       string
		secret   string
		expected *oauth.Config
	}{
		{
			name:   "Test to get config using id and secret",
			id:     "test",
			secret: "secret",
			expected: &oauth.Config{
				ClientID:     "test",
				ClientSecret: "secret",
			},
		},
		{
			name:   "Test to get config using empty id and secret",
			id:     "",
			secret: "",
			expected: &oauth.Config{
				ClientID:     "",
				ClientSecret: "",
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			config := NewConfig(test.id, test.secret)
			if config.ClientID != test.expected.ClientID {
				t.Errorf("expected %v got %v\n", config, test.expected.ClientID)
			}
			if config.ClientSecret != test.expected.ClientSecret {
				t.Errorf("expected %v got %v\n", config, test.expected.ClientSecret)
			}
		})
	}
}

func TestOpenURL(t *testing.T) {
    err := openUrl("")

    if err == nil {
        t.Error("openurl not handling all the cases\n")
    }
}
