package contentfile

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/ezeoleaf/larry/cache"
	"github.com/ezeoleaf/larry/config"
	"github.com/ezeoleaf/larry/domain"
	"github.com/go-redis/redis/v8"
)

func TestGetJsonContentFromReader(t *testing.T) {
	for _, tc := range []struct {
		Name             string
		CachedItems      []string
		BlacklistedItems []string
		ContentFile      string
		ExpectedContent  *domain.Content
		ExpectedError    string
	}{
		{
			Name:             "Test success",
			CachedItems:      []string{"title-1"},
			BlacklistedItems: []string{"title-2"},
			ContentFile:      `[{"Title":"title-0","Subtitle":"subtitle-0","URL":"url-0","ExtraData":["extradata-0-1","extradata-0-2"]},{"Title":"title-1","Subtitle":"subtitle-1","URL":"url-1","ExtraData":["extradata-1-1","extradata-1-2"]},{"Title":"title-2","Subtitle":"subtitle-2","URL":"url-2","ExtraData":["extradata-2-1","extradata-2-2"]}]`,
			ExpectedContent: &domain.Content{
				Title:     StringToPointer("title-0"),
				Subtitle:  StringToPointer("subtitle-0"),
				URL:       StringToPointer("url-0"),
				ExtraData: []string{"extradata-0-1", "extradata-0-2"},
			},
		},
		{
			Name:             "Test missing title",
			CachedItems:      []string{},
			BlacklistedItems: []string{"title-2"},
			ContentFile:      `[{"Title":"","Subtitle":"subtitle-0","URL":"url-0","ExtraData":["extradata-0-1","extradata-0-2"]},{"Title":"title-1","Subtitle":"subtitle-1","URL":"url-1","ExtraData":["extradata-1-1","extradata-1-2"]},{"Title":"title-2","Subtitle":"subtitle-2","URL":"url-2","ExtraData":["extradata-2-1","extradata-2-2"]}]`,
			ExpectedContent: &domain.Content{
				Title:     StringToPointer("title-1"),
				Subtitle:  StringToPointer("subtitle-1"),
				URL:       StringToPointer("url-1"),
				ExtraData: []string{"extradata-1-1", "extradata-1-2"},
			},
		},
		{
			Name:             "Test empty array",
			CachedItems:      []string{"title-1"},
			BlacklistedItems: []string{"title-2"},
			ContentFile:      `[]`,
			ExpectedContent:  nil,
		},
		{
			Name:             "Test empty file",
			CachedItems:      []string{"title-1"},
			BlacklistedItems: []string{"title-2"},
			ContentFile:      ``,
			ExpectedContent:  nil,
		},
		{
			Name:             "Test malformed file",
			CachedItems:      []string{},
			BlacklistedItems: []string{},
			ContentFile:      `"title-0"`, // CSV file provided instead of JSON
			ExpectedContent:  nil,
			ExpectedError:    `parse error on line 1, column 11: extraneous or missing " in quoted-field`,
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {

			mr, _ := miniredis.Run()
			ro := &redis.Options{
				Addr: mr.Addr(),
			}
			cc := cache.NewClient(ro)

			for _, item := range tc.CachedItems {
				err := cc.Set(item, "1", 0)
				if err != nil {
					t.Error("could not set key")
				}
			}
			for _, item := range tc.BlacklistedItems {
				err := cc.Set("blacklist-"+item, "1", 0)
				if err != nil {
					t.Error("could not set key")
				}
			}

			cfg := config.Config{ContentFile: "./test.json"}
			p, err := NewProvider(cfg, cc, NewMockS3Client())
			if err != nil {
				fmt.Println(err)
				t.Error(err)
			}

			if content, err := p.FileReader.getContentFromReader(strings.NewReader(tc.ContentFile), p.skipCachedRecord); err != nil {
				if tc.ExpectedError != err.Error() {
					fmt.Println(err.Error())
					t.Error(err)
				}
			} else {
				if content == nil && tc.ExpectedContent == nil {
					// success
				} else if content == nil && tc.ExpectedContent != nil {
					t.Errorf("expected %v as value, got nil instead", *tc.ExpectedContent.Title)
				} else if content != nil && tc.ExpectedContent == nil {
					t.Errorf("expected nil as value, got %v instead", *content.Title)
				} else if *content.Title != *tc.ExpectedContent.Title {
					t.Errorf("expected %v as value, got %v instead", tc.ExpectedContent.Title, *content.Title)
				} else {
					// compare returned object
					expected, _ := json.Marshal(tc.ExpectedContent)
					got, _ := json.Marshal(content)
					if string(expected) != string(got) {
						t.Errorf("expected %v as value, got %v instead", string(expected), string(got))
					}
				}
			}
		})
	}
}
