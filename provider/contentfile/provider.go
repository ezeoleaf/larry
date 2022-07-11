package contentfile

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
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
	S3Client    s3FileReader
}

type ContentFileReader interface {
	getContentFromReader(handle io.Reader, skip func(string) bool) (*domain.Content, error)
}

func NewProvider(cfg config.Config, cacheClient cache.Client, s3client s3FileReader) (Provider, error) {

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
		S3Client:    s3client,
	}
	return p, nil
}

// GetContentToPublish returns content to publish to be used by the publishers
func (p Provider) GetContentToPublish() (*domain.Content, error) {
	u, err := url.Parse(p.Config.ContentFile)
	if err != nil {
		return nil, err
	}

	var content *domain.Content
	if strings.ToLower(u.Scheme) == "s3" {
		content, err = p.getContentFromS3Bucket(u.Hostname(), u.Path[1:])
	} else {
		content, err = p.getContentFromFile(p.Config.ContentFile)
	}

	if content != nil {
		p.addToCache(*content.Title)
	}

	return content, err
}

func (p Provider) getContentFromS3Bucket(bucket, key string) (*domain.Content, error) {
	if objectReader, err := p.S3Client.GetObject(bucket, key); err != nil {
		return nil, err
	} else {
		return p.FileReader.getContentFromReader(objectReader, p.skipCachedRecord)
	}
}

func (p Provider) getContentFromFile(fileName string) (*domain.Content, error) {
	if fileName != "" {
		f, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		return p.FileReader.getContentFromReader(f, p.skipCachedRecord)
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

func (p Provider) addToCache(title string) error {
	key := cacheKey(p.Config.GetCacheKeyPrefix(), title)
	return p.CacheClient.Set(key, true, p.cacheExpirationMinutes())
}
