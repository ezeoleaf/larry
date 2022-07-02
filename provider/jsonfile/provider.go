package jsonfile

import (
	"encoding/json"
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
	log.Print("New Jsonfile Provider")
	p := Provider{Config: cfg, CacheClient: cacheClient}
	return p
}

// GetContentToPublish returns content to publish to be used by the publishers
func (p Provider) GetContentToPublish() (*domain.Content, error) {
	return p.getContentFromFile(p.Config.ContentFile)
}

func (p Provider) getContentFromFile(jsonFileName string) (*domain.Content, error) {
	if jsonFileName != "" {
		f, err := os.OpenFile(jsonFileName, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		return p.getContentFromReader(f)
	}

	return nil, fmt.Errorf("No json file specified")
}

func (p Provider) getContentFromReader(handle io.Reader) (*domain.Content, error) {

	size := 1
	reservoir := domain.Content{}
	rand.Seed(time.Now().UnixNano())

	decoder := json.NewDecoder(handle)
	if _, err := decoder.Token(); err != nil {
		if err.Error() == "EOF" {
			return nil, nil
		}
		return nil, err
	}

	count := 0
	for decoder.More() {
		data := new(domain.Content)
		if err := decoder.Decode(data); err != nil {
			return nil, err
		}

		if data.Title == nil || *data.Title == "" {
			log.Println("content missing title, skipping record")
			continue
		}

		// check for content in cache/blacklist
		if p.isCached(*data.Title) {
			continue
		} else if p.isBlacklisted(*data.Title) {
			log.Printf("content blacklisted: %s\n", *data.Title)
			continue
		}

		// reservoir sampling technique
		if count < size {
			reservoir = *data
		} else {
			j := rand.Intn(count + 1)
			if j < size {
				reservoir = *data
			}
		}
		count++
	}

	if count > 0 {
		key := cacheKey(p.Config.GetCacheKeyPrefix(), *reservoir.Title)
		p.CacheClient.Set(key, true, p.cacheExpirationMinutes())
		return &reservoir, nil
	}

	return nil, nil
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
