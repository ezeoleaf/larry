package csvfile

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/ezeoleaf/larry/cache"
	"github.com/ezeoleaf/larry/config"
	"github.com/ezeoleaf/larry/domain"
	"github.com/go-redis/redis/v8"
)

type Provider struct {
	CacheClient cache.Client
	Config      config.Config
}

func NewProvider(cfg config.Config, cacheClient cache.Client) Provider {
	log.Print("New Csvfile Provider")
	p := Provider{Config: cfg, CacheClient: cacheClient}
	return p
}

// GetContentToPublish returns content to publish to be used by the publishers
func (p Provider) GetContentToPublish() (*domain.Content, error) {
	return p.getContentFromFile(p.Config.ContentFile)
}

func (p Provider) getContentFromFile(csvFileName string) (*domain.Content, error) {
	if csvFileName != "" {
		f, err := os.OpenFile(csvFileName, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		return p.getContentFromReader(f)
	}

	return nil, fmt.Errorf("No csv file specified")
}

func (p Provider) getContentFromReader(handle io.Reader) (*domain.Content, error) {

	size := 1
	var reservoir []string
	rand.Seed(time.Now().UnixNano())

	count := 0
	skipHeader := p.Config.SkipCsvHeader
	csvReader := csv.NewReader(handle)
	for {
		rec, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// skip header line
		if skipHeader {
			skipHeader = false
			continue
		}

		if rec[0] == "" {
			log.Println("content missing title, skipping record")
			continue
		}

		// check for content in cache/blacklist
		if p.isCached(rec[0]) {
			continue
		} else if p.isBlacklisted(rec[0]) {
			log.Printf("content blacklisted: %s\n", rec[0])
			continue
		}

		// reservoir sampling technique
		if count < size {
			reservoir = rec
		} else {
			j := rand.Intn(count + 1)
			if j < size {
				reservoir = rec
			}
		}

		count++
	}

	if count > 0 {
		if content, err := convertCsvToContent(reservoir); err != nil {
			return nil, err
		} else {
			key := cacheKey(p.Config.GetCacheKeyPrefix(), *content.Title)
			p.CacheClient.Set(key, true, p.cacheExpirationMinutes())
			return content, nil
		}
	}

	return nil, nil
}

func convertCsvToContent(rec []string) (*domain.Content, error) {
	content := domain.Content{ExtraData: []string{}}
	if len(rec) > 0 {
		content.Title = StringToPointer(rec[0])
	}
	if len(rec) > 1 {
		content.Subtitle = StringToPointer(rec[1])
	}
	if len(rec) > 2 {
		content.URL = StringToPointer(rec[2])
	}
	if len(rec) > 3 {
		// number of extra data fields is variable for CSV
		content.ExtraData = make([]string, len(rec)-3)
		for i := 3; i < len(rec); i++ {
			content.ExtraData[i-3] = rec[i]
		}
	}
	return &content, nil
}

func StringToPointer(in string) *string {
	return &in
}

func (p Provider) isCached(title string) bool {
	key := cacheKey(p.Config.GetCacheKeyPrefix(), title)
	_, err := p.CacheClient.Get(key)
	if err != redis.Nil {
		return true
	}
	return false
}

func (p Provider) cacheExpirationMinutes() time.Duration {
	expirationMinutes := p.Config.CacheSize * p.Config.Periodicity
	if expirationMinutes < 0 {
		expirationMinutes = 0
	}
	return time.Duration(expirationMinutes) * time.Minute
}

func (p Provider) isBlacklisted(title string) bool {
	if _, err := p.CacheClient.Get("blacklist-" + title); err != redis.Nil {
		return true
	}
	return false
}

func cacheKey(cacheKeyPrefix string, title string) string {
	return cacheKeyPrefix + title
}
