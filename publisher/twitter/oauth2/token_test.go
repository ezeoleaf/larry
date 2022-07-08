package oauth2

import (
	"os"
	"testing"

	oauth "golang.org/x/oauth2"
)

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
	filename = "./test_files/.env"

	err := store("TESTED", "TRUE")
	if err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	filename = "./test_files/.env"

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
	file, err := os.Open("./test_files/.env")
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

func TestGetToken(t *testing.T) {
	filename = "./test_files/.env"

    expected := &oauth.Token{}

    got := getToken()

    if *got != *expected {
        t.Errorf("got %v and expected %v\n", *got, *expected)
    }
}
