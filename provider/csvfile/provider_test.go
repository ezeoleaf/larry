package csvfile

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

func TestGetContentFromReader(t *testing.T) {
	for _, tc := range []struct {
		Name             string
		CachedItems      []string
		BlacklistedItems []string
		ContentFile      string
		SkipHeader       bool
		ExpectedContent  *domain.Content
		ExpectedError    string
	}{
		{
			Name:             "Test success",
			CachedItems:      []string{"title-0"},
			BlacklistedItems: []string{"title-1"},
			ContentFile: `"title-0","subtitle-0","url-0","extradata-0-1","extradata-0-2"
title-1,subtitle-1,url-1,extradata-1-1,extradata-1-2
"title-2","subtitle-2","url-2","extradata-2-1","extradata-2-2,embedded comma"`,
			ExpectedContent: &domain.Content{
				Title:     StringToPointer("title-2"),
				Subtitle:  StringToPointer("subtitle-2"),
				URL:       StringToPointer("url-2"),
				ExtraData: []string{"extradata-2-1", "extradata-2-2,embedded comma"},
			},
		},
		{
			Name:             "Test record has too many extra data",
			CachedItems:      []string{"title-0"},
			BlacklistedItems: []string{"title-1"},
			ContentFile: `"title-0","subtitle-0","url-0","extradata-0-1","extradata-0-2"
title-1,subtitle-1,url-1,extradata-1-1,extradata-1-2
"title-2","subtitle-2","url-2","extradata-2-1","extradata-2-2,embedded comma","this will cause an error"`,
			ExpectedContent: nil,
			ExpectedError:   "record on line 3: wrong number of fields",
		},
		{
			Name:             "Test record has too few extra data",
			CachedItems:      []string{"title-0"},
			BlacklistedItems: []string{"title-1"},
			ContentFile: `"title-0","subtitle-0","url-0","extradata-0-1","extradata-0-2"
title-1,subtitle-1,url-1,extradata-1-1,extradata-1-2
"title-2","subtitle-2","url-2","extradata-2-1"`,
			ExpectedContent: nil,
			ExpectedError:   "record on line 3: wrong number of fields",
		},
		{
			Name:             "Test last record extra data is blank",
			CachedItems:      []string{"title-0"},
			BlacklistedItems: []string{"title-1"},
			ContentFile: `"title-0","subtitle-0","url-0","extradata-0-1","extradata-0-2"
title-1,subtitle-1,url-1,extradata-1-1,extradata-1-2
"title-2","subtitle-2","url-2","extradata-2-1",""`,
			ExpectedContent: &domain.Content{
				Title:     StringToPointer("title-2"),
				Subtitle:  StringToPointer("subtitle-2"),
				URL:       StringToPointer("url-2"),
				ExtraData: []string{"extradata-2-1", ""},
			},
		},
		{
			Name:             "Test skip header",
			CachedItems:      []string{"title-1"},
			BlacklistedItems: []string{"title-2"},
			ContentFile: `Title,Subtitle,URL,ExtraData1,ExtraData2
"title-0","subtitle-0","url-0","extradata-0-1","extradata-0-2"
title-1,subtitle-1,url-1,extradata-1-1,extradata-1-2
"title-2","subtitle-2","url-2","extradata-2-1","extradata-2-2,embedded comma"`,
			SkipHeader: true,
			ExpectedContent: &domain.Content{
				Title:     StringToPointer("title-0"),
				Subtitle:  StringToPointer("subtitle-0"),
				URL:       StringToPointer("url-0"),
				ExtraData: []string{"extradata-0-1", "extradata-0-2"},
			},
		},
		{
			Name:             "Test skip header, number of header fields don't match other records",
			CachedItems:      []string{"title-1"},
			BlacklistedItems: []string{"title-2"},
			ContentFile: `Title,Subtitle,URL,ExtraData1
"title-0","subtitle-0","url-0","extradata-0-1","extradata-0-2"
title-1,subtitle-1,url-1,extradata-1-1,extradata-1-2
"title-2","subtitle-2","url-2","extradata-2-1","extradata-2-2,embedded comma"`,
			SkipHeader:      true,
			ExpectedContent: nil,
			ExpectedError:   "record on line 2: wrong number of fields",
		},
		{
			Name:             "Test empty file",
			CachedItems:      []string{},
			BlacklistedItems: []string{},
			ContentFile:      "",
			ExpectedContent:  nil,
		},
		{
			Name:             "Test just title",
			CachedItems:      []string{},
			BlacklistedItems: []string{},
			ContentFile:      `"title-0"`,
			ExpectedContent: &domain.Content{
				Title:     StringToPointer("title-0"),
				ExtraData: []string{},
			},
		},
		{
			Name:             "Test no title",
			CachedItems:      []string{},
			BlacklistedItems: []string{},
			ContentFile:      `"","title-0"`,
			ExpectedContent:  nil,
		},
		{
			Name:             "Test many ExtraData",
			CachedItems:      []string{},
			BlacklistedItems: []string{},
			ContentFile:      `title-0,,,ExtraData1,ExtraData2,ExtraData3,ExtraData4,ExtraData5`,
			ExpectedContent: &domain.Content{
				Title:     StringToPointer("title-0"),
				Subtitle:  StringToPointer(""),
				URL:       StringToPointer(""),
				ExtraData: []string{"ExtraData1", "ExtraData2", "ExtraData3", "ExtraData4", "ExtraData5"},
			},
		},
		{
			Name:             "Test malformed file",
			CachedItems:      []string{},
			BlacklistedItems: []string{},
			ContentFile:      `"title-0""`, // extra double quote
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
				cc.Set(item, "1", 0)
			}
			for _, item := range tc.BlacklistedItems {
				cc.Set("blacklist-"+item, "1", 0)
			}

			cfg := config.Config{SkipCsvHeader: tc.SkipHeader}
			p := Provider{Config: cfg, CacheClient: cc}

			if content, err := p.getContentFromReader(strings.NewReader(tc.ContentFile)); err != nil {
				if tc.ExpectedError != err.Error() {
					fmt.Println(err)
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
					t.Errorf("expected %v as value, got %v instead", *&tc.ExpectedContent.Title, *content.Title)
				} else {
					// compare returned object
					expected, _ := json.Marshal(tc.ExpectedContent)
					got, _ := json.Marshal(content)
					if string(expected) != string(got) {
						t.Errorf("expected %v as value, got %v instead", string(expected), string(got))
					}

					// check cache for returned object
					if _, err := p.CacheClient.Get(*tc.ExpectedContent.Title); err != nil {
						t.Errorf("expected %v not found in cache", *tc.ExpectedContent.Title)
					}
				}
			}
		})
	}
}
