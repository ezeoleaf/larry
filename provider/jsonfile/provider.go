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

	size := 1
	reservoir := make([]domain.Content, size)
	rand.Seed(time.Now().UnixNano())

	// TODO: test with empty file or empty array
	decoder := json.NewDecoder(handle)
	decoder.Token()

	count := 0
	for decoder.More() {
		data := new(domain.Content)
		if err := decoder.Decode(data); err != nil {
			return nil, err
		}

		// TODO: test with missing title
		if data.Title == nil {
			log.Println("content missing title")
			continue
		}

		// check for content in cache/blacklist
		if p.isCached(*data.Title) {
			log.Printf("content cached: %s\n", *data.Title)
			continue
		} else if p.isBlacklisted(*data.Title) {
			log.Printf("content blacklisted: %s\n", *data.Title)
			continue
		}

		// reservoir sampling technique
		if count < size {
			// always fill the first x elements of the slice
			reservoir[count] = *data
		} else {
			// after the first x elements find a random slot in the slice
			j := rand.Intn(count + 1)
			if j < size {
				reservoir[j] = *data
			}
		}
		count++
	}

	// jsonString, _ := json.Marshal(reservoir)
	// fmt.Println(string(jsonString))

	// TODO: handle nothing found
	p.CacheClient.Set(*reservoir[0].Title, true, p.cacheExpirationMinutes())

	return &reservoir[0], nil
}

func StringToPointer(in string) *string {
	return &in
}

func (p Provider) isCached(title string) bool {
	_, err := p.CacheClient.Get(title)
	if err != redis.Nil {
		return true
	}
	return false
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
