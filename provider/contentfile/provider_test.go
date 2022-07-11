package contentfile

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/ezeoleaf/larry/cache"
	"github.com/ezeoleaf/larry/config"
	"github.com/ezeoleaf/larry/domain"
	"github.com/go-redis/redis/v8"
)

func TestGetContentFromFile(t *testing.T) {
	for _, tc := range []struct {
		Name             string
		CachedItems      []string
		BlacklistedItems []string
		ContentFile      string
		ExpectedContent  *domain.Content
		ExpectedError    string
	}{
		{
			Name:             "Success json",
			CachedItems:      []string{"title-0"},
			BlacklistedItems: []string{"title-1"},
			ContentFile:      "./testdata/test.json",
			ExpectedContent: &domain.Content{
				Title:     StringToPointer("title-2"),
				Subtitle:  StringToPointer("subtitle-2"),
				URL:       StringToPointer("url-2"),
				ExtraData: []string{"extradata-2-1", "extradata-2-2"},
			},
		},
		{
			Name:             "Success csv",
			CachedItems:      []string{"title-0"},
			BlacklistedItems: []string{"title-1"},
			ContentFile:      "./testdata/test.json",
			ExpectedContent: &domain.Content{
				Title:     StringToPointer("title-2"),
				Subtitle:  StringToPointer("subtitle-2"),
				URL:       StringToPointer("url-2"),
				ExtraData: []string{"extradata-2-1", "extradata-2-2"},
			},
		},
		{
			Name:             "Error no file extension",
			CachedItems:      []string{"title-0"},
			BlacklistedItems: []string{"title-1"},
			ContentFile:      "./testdata/test", // no file extension provided
			ExpectedContent:  nil,
			ExpectedError:    "no file extension provided, unable to determine file format",
		},
		{
			Name:             "Error invalid file extension",
			CachedItems:      []string{"title-0"},
			BlacklistedItems: []string{"title-1"},
			ContentFile:      "./testdata/test.txt", // this file extension is not supported
			ExpectedContent:  nil,
			ExpectedError:    "unsupported content file format: .txt",
		},
		{
			Name:             "Success json S3",
			CachedItems:      []string{"title-0"},
			BlacklistedItems: []string{"title-1"},
			ContentFile:      "s3://testdata/test.json", // S3, bucket, key
			ExpectedContent: &domain.Content{
				Title:     StringToPointer("title-2"),
				Subtitle:  StringToPointer("subtitle-2"),
				URL:       StringToPointer("url-2"),
				ExtraData: []string{"extradata-2-1", "extradata-2-2"},
			},
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

			cfg := config.Config{ContentFile: tc.ContentFile}
			p, err := NewProvider(cfg, cc, NewMockS3Client())
			if err != nil {
				if tc.ExpectedError != err.Error() {
					fmt.Println(err.Error())
					t.Error(err)
				}
				return
			}

			if content, err := p.GetContentToPublish(); err != nil {
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

					// check cache for returned object
					if _, err := p.CacheClient.Get(*tc.ExpectedContent.Title); err != nil {
						t.Errorf("expected %v not found in cache", *tc.ExpectedContent.Title)
					}
				}
			}
		})
	}
}

type MockS3Client struct {
}

func NewMockS3Client() MockS3Client {
	return MockS3Client{}
}

func (s MockS3Client) GetObject(bucket, key string) (io.Reader, error) {
	// convert bucket and key to a local filename
	f, err := os.OpenFile(fmt.Sprintf("./%s/%s", bucket, key), os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}
	return f, nil
}
