package blacklist

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/ezeoleaf/larry/cache"
	"github.com/go-redis/redis/v8"
)

func TestBlacklist(t *testing.T) {
	for _, tc := range []struct {
		Name                  string
		BlacklistFileContents string
		ExpectedResults       []string
	}{
		{
			Name: "Simple file",
			BlacklistFileContents: `111
			222`,
			ExpectedResults: []string{"blacklist-golang-111", "blacklist-golang-222"},
		},
		{
			Name: "Complex file",
			BlacklistFileContents: `# Comment in first line
			777 # Comment after repo Id
			888
			# Blank line follows

				# Tab then comment
			999
			# Last line is comment`,
			ExpectedResults: []string{"blacklist-golang-777", "blacklist-golang-888", "blacklist-golang-999"},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {

			mr, _ := miniredis.Run()
			ro := &redis.Options{
				Addr: mr.Addr(),
			}

			cc := cache.NewClient(ro)

			// test loading blacklist
			blacklistKeyPrefix := "blacklist-golang-"
			if err := readerLoader(cc, strings.NewReader(tc.BlacklistFileContents), blacklistKeyPrefix); err != nil {
				t.Error(err)
			}

			if count, err := keyCount(cc); err != nil {
				t.Error("Error retrieving key count")
			} else if len(tc.ExpectedResults) != count {
				t.Errorf("Key count found %d doesn't match expected count %d", count, len(tc.ExpectedResults))
			}

			for _, ex := range tc.ExpectedResults {
				if _, err := cc.Get(ex); !keyFound(err) {
					t.Errorf("No value found for expected result %s", ex)
				}
			}

			// test clearing blacklist
			if err := clear(cc, blacklistKeyPrefix); err != nil {
				t.Error(err)
			}

			for _, ex := range tc.ExpectedResults {
				if _, err := cc.Get(ex); keyFound(err) {
					t.Error(fmt.Sprintf("Expected key to be deleted %s", ex))
				}
			}
		})
	}
}

func keyFound(err error) bool {
	if err != redis.Nil {
		return true
	}
	return false
}

func keyCount(cc cache.Client) (int, error) {
	count := 0
	cc.Scan("*", func(ctx context.Context, s string) error {
		count++
		return nil
	})
	return count, nil
}
