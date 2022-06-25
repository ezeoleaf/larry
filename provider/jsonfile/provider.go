package jsonfile

import (
	"encoding/json"
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
	log.Print("New Json Provider")
	p := Provider{Config: cfg, CacheClient: cacheClient}
	return p
}

// GetContentToPublish returns content to publish to be used by the publishers
func (p Provider) GetContentToPublish() (*domain.Content, error) {
	return p.getContentFromFile(p.Config.LocalFile)
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

	return nil, nil
}

func (p Provider) getContentFromReader(handle io.Reader) (*domain.Content, error) {
	decoder := json.NewDecoder(handle)

	size := 1
	reservoir := make([]domain.Content, size)
	rand.Seed(time.Now().UnixNano())

	decoder.Token()

	count := 0
	for decoder.More() {
		data := new(domain.Content)
		decoder.Decode(data)

		// check cache/blacklist here for title
		_, err := p.CacheClient.Get(*data.Title)
		if err != redis.Nil {
			continue
		}
		if p.isBlacklisted(*data.Title) {
			continue
		}

		if count < size {
			reservoir[count] = *data
		} else {
			j := rand.Intn(count + 1)
			if j < size {
				reservoir[j] = *data
			}
		}
		count++
	}

	// jsonString, _ := json.Marshal(reservoir)
	// fmt.Println(string(jsonString))

	p.CacheClient.Set(*reservoir[0].Title, true, p.cacheExpirationMinutes())

	return &reservoir[0], nil
}

func StringToPointer(in string) *string {
	return &in
}

// TODO: this is repeated - show be in cache?
func (p Provider) cacheExpirationMinutes() time.Duration {
	expirationMinutes := p.Config.CacheSize * p.Config.Periodicity
	if expirationMinutes < 0 {
		expirationMinutes = 0
	}
	return time.Duration(expirationMinutes) * time.Minute
}

// TODO: this is repeated
func (p Provider) isBlacklisted(title string) bool {
	if _, err := p.CacheClient.Get("blacklist-" + title); err != redis.Nil {
		return true
	}
	return false
}
