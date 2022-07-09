package oauth2

import (
	"context"
	"os"
	"testing"
	"time"

	oauth "golang.org/x/oauth2"
)

var now = time.Now()

func TestFileName(t *testing.T) {
	name := getFileName("test.env")

	home, err := os.UserHomeDir()
	if err != nil {
		t.Error(err)
	}
	home += "/.test.env"
	if name != home {
		t.Errorf("filename found %v instead of %v", name, home)
	}
}

func TestStore(t *testing.T) {
	filename = "./test_files/larry.env"

	err := store("TESTED", "TRUE")
	if err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	filename = "./test_files/larry.env"

	for _, test := range []struct {
		name     string
		key      string
		expected string
	}{
		{
			name:     "Get Success",
			key:      "TESTED",
			expected: "TRUE",
		},
		{
			name:     "Get Failure",
			key:      "ANONYMOUS",
			expected: "",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := get(test.key)

			if got != test.expected {
				t.Errorf("got %v and expected %v\n", got, test.expected)
			}
		})
	}
}

func TestGetFileContent(t *testing.T) {
	file, err := os.Open("./test_files/larry.env")
	if err != nil {
		t.Error("failed to open file")
	}
	defer file.Close()

	lines, err := getFileContent(file)
	if err != nil {
		t.Error(err)
	}

	if len(lines) != 5 {
		t.Errorf("found length=%v instead of 5\n", len(lines))
	}
}

func TestGetKey(t *testing.T) {
	content := []string{
		"#comment",
		"",
		"BOTNAME=larry",
	}
	for _, test := range []struct {
		name    string
		found   bool
		index   int
		key     string
		content []string
	}{
		{
			name:    "Test for successful search",
			found:   true,
			index:   2,
			key:     "BOTNAME",
			content: content,
		},
		{
			name:    "Test for unsuccessful search",
			found:   false,
			index:   -1,
			key:     "TEST",
			content: content,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			ok, i := getKey(test.content, test.key)
			if ok != test.found {
				t.Errorf("found is %v instead of %v", ok, test.found)
			}
			if i != test.index {
				t.Errorf("index is %v instead of %v", i, test.index)
			}
		})
	}
}

func TestStoreToken(t *testing.T) {
	filename = "./test_files/token.env"

	tok := &oauth.Token{
		AccessToken:  "accesstoken",
		RefreshToken: "refreshtoken",
		Expiry:       now,
	}

	err := storeToken(tok)
	if err != nil {
		t.Error(err)
	}
}

func TestGetToken(t *testing.T) {
	time, _ := time.Parse(time.RFC1123, now.Format(time.RFC1123))

	for _, test := range []struct {
		name     string
		filename string
		expected *oauth.Token
	}{
		{
			name:     "Test to perform unsuccessful token retrieval",
			filename: "./test_files/larry.env",
			expected: &oauth.Token{},
		},
		{
			name:     "Test to perform successful token retrieval",
			filename: "./test_files/token.env",
			expected: &oauth.Token{
				AccessToken:  "accesstoken",
				RefreshToken: "refreshtoken",
				Expiry:       time,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			filename = test.filename

			got := getToken()

			if *got != *test.expected {
				t.Errorf("got %v and expected %v\n", *got, *test.expected)
			}
		})
	}
}

func TestRegenerateToken(t *testing.T) {
    filename = "./test_files/token.env"

    ctx := context.Background()
    config := NewConfig("randomID", "superSecret")

    expected := getToken()

    tok, _ := regenerateToken(ctx, config, expected)

    if tok != expected {
        t.Errorf("expected %v; got %v\n", expected, tok)
    }
}
