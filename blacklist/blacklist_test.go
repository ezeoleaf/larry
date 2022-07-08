package blacklist

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/ezeoleaf/larry/mock"
)

func TestLoad(t *testing.T) {
	expectedResults := []string{}

	for _, tc := range []struct {
		Name              string
		Key               string
		BlacklistFileName string
		ExpectedResults   []string
		CacheClient       mock.CacheClientMock
		ShouldError       bool
	}{
		{
			Name: "Should fail do to clear function",
			CacheClient: mock.CacheClientMock{
				DelFn: func(key string) error {
					return errors.New("some error")
				},
			},
			ShouldError:       true,
			Key:               "some-key",
			BlacklistFileName: "some-name.txt",
		},
		{
			Name: "Should fail do to not existing file",
			CacheClient: mock.CacheClientMock{
				ScanFn: func(key string, action func(context.Context, string) error) error {
					return nil
				},
				DelFn: func(key string) error {
					return nil
				},
			},
			ShouldError:       true,
			Key:               "some-key",
			BlacklistFileName: "some-file.txt",
		},
		{
			Name: "Should fail do to cache set failing",
			CacheClient: mock.CacheClientMock{
				ScanFn: func(key string, action func(context.Context, string) error) error {
					return nil
				},
				DelFn: func(key string) error {
					return nil
				},
				SetFn: func(key string, value interface{}, exp time.Duration) error {
					return errors.New("some errors")
				},
			},
			ShouldError:       true,
			Key:               "some-key",
			BlacklistFileName: "./testdata/blacklist.txt",
		},
		{
			Name: "Should not fail due to no file",
			CacheClient: mock.CacheClientMock{
				ScanFn: func(key string, action func(context.Context, string) error) error {
					return nil
				},
				DelFn: func(key string) error {
					return nil
				},
				SetFn: func(key string, value interface{}, exp time.Duration) error {
					return errors.New("some errors")
				},
			},
			ShouldError:       false,
			Key:               "some-key",
			BlacklistFileName: "",
		},
		{
			Name: "Should not fail",
			CacheClient: mock.CacheClientMock{
				ScanFn: func(key string, action func(context.Context, string) error) error {
					return nil
				},
				DelFn: func(key string) error {
					return nil
				},
				SetFn: func(key string, value interface{}, exp time.Duration) error {
					expectedResults = append(expectedResults, key)
					return nil
				},
			},
			ShouldError:       false,
			Key:               "some-key-",
			BlacklistFileName: "./testdata/blacklist.txt",
			ExpectedResults:   []string{"blacklist-some-key-test", "blacklist-some-key-test2"},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			blacklistClient := NewClient(tc.CacheClient)
			err := blacklistClient.Load(tc.BlacklistFileName, tc.Key)

			if err == nil && tc.ShouldError {
				t.Error("should error but not error received")
			} else if err != nil && !tc.ShouldError {
				t.Errorf("should not error but got %s", err)
			}

			if len(expectedResults) != len(tc.ExpectedResults) {
				t.Errorf("expected %v results, but got %v", len(expectedResults), len(tc.ExpectedResults))
			}

			for i, er := range tc.ExpectedResults {
				if er != expectedResults[i] {
					t.Errorf("expected %s in position %v, but got %s", er, i, expectedResults[i])
				}
			}

			expectedResults = []string{}
		})

	}
}

func TestLoadFromReader(t *testing.T) {

	expectedResults := []string{}

	for _, tc := range []struct {
		Name                  string
		BlacklistFileContents string
		Key                   string
		ExpectedResults       []string
		CacheClient           mock.CacheClientMock
		ShouldError           bool
	}{
		{
			Name: "Should fail do to not able to set",
			CacheClient: mock.CacheClientMock{
				SetFn: func(key string, value interface{}, exp time.Duration) error {
					return errors.New("some error")
				},
			},
			ShouldError: true,
			Key:         "some-key",
			BlacklistFileContents: `# Comment in first line
			777 # Comment after repo Id
			888
			# Blank line follows

				# Tab then comment
			999
			# Last line is comment`,
			ExpectedResults: []string{},
		},
		{
			Name: "Should not fail",
			CacheClient: mock.CacheClientMock{
				SetFn: func(key string, value interface{}, exp time.Duration) error {
					expectedResults = append(expectedResults, key)
					return nil
				},
			},
			ShouldError: false,
			Key:         "blacklist-golang-",
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
			blacklistClient := NewClient(tc.CacheClient)
			err := blacklistClient.LoadFromReader(strings.NewReader(tc.BlacklistFileContents), tc.Key)

			if err == nil && tc.ShouldError {
				t.Error("should error but not error received")
			} else if err != nil && !tc.ShouldError {
				t.Errorf("should not error but got %s", err)
			}

			if len(expectedResults) != len(tc.ExpectedResults) {
				t.Errorf("expected %v results, but got %v", len(expectedResults), len(tc.ExpectedResults))
			}

			for i, er := range tc.ExpectedResults {
				if er != expectedResults[i] {
					t.Errorf("expected %s in position %v, but got %s", er, i, expectedResults[i])
				}
			}

			expectedResults = []string{}
		})
	}
}
