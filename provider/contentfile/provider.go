package contentfile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/ezeoleaf/larry/cache"
	"github.com/ezeoleaf/larry/config"
	"github.com/ezeoleaf/larry/domain"
	"github.com/go-redis/redis/v8"
)

type Provider struct {
	CacheClient cache.Client
	Config      config.Config
	FileReader  ContentFileReader
}

type ContentFileReader interface {
	getContentFromReader(handle io.Reader, skip func(string) bool) (*domain.Content, error)
}

func NewProvider(cfg config.Config, cacheClient cache.Client) (Provider, error) {

	ext := filepath.Ext(cfg.ContentFile)
	if ext == "" {
		return Provider{}, fmt.Errorf("no file extension provided, unable to determine file format")
	}

	var fileReader ContentFileReader
	switch ext {
	case ".json":
		fileReader = NewJsonFileReader()
	case ".csv":
		fileReader = NewCsvFileReader(cfg.SkipCsvHeader)
	default:
		return Provider{}, fmt.Errorf("unsupported content file format: %s", ext)
	}

	p := Provider{
		Config:      cfg,
		CacheClient: cacheClient,
		FileReader:  fileReader,
	}
	return p, nil
}

// GetContentToPublish returns content to publish to be used by the publishers
func (p Provider) GetContentToPublish() (*domain.Content, error) {
	return p.getContentFromFile(p.Config.ContentFile)
}

func (p Provider) getContentFromFile(fileName string) (*domain.Content, error) {
	if fileName != "" {
		f, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		if content, err := p.FileReader.getContentFromReader(f, p.skipCachedRecord); err != nil {
			return nil, err
		} else {
			p.addToCache(*content.Title)
			return content, nil
		}
	}

	return nil, fmt.Errorf("no content file specified")
}

func StringToPointer(in string) *string {
	return &in
}

func (p Provider) skipCachedRecord(title string) bool {
	if p.isCached(title) {
		return true
	} else if p.isBlacklisted(title) {
		return true
	}
	return false
}

func (p Provider) isCached(title string) bool {
	key := cacheKey(p.Config.GetCacheKeyPrefix(), title)
	_, err := p.CacheClient.Get(key)
	return err != redis.Nil
}

func (p Provider) isBlacklisted(title string) bool {
	if _, err := p.CacheClient.Get("blacklist-" + title); err != redis.Nil {
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

func cacheKey(cacheKeyPrefix string, title string) string {
	return cacheKeyPrefix + title
}

func (p Provider) addToCache(title string) {
	key := cacheKey(p.Config.GetCacheKeyPrefix(), title)
	p.CacheClient.Set(key, true, p.cacheExpirationMinutes())
}
