package blacklist

import (
	"bufio"
	"context"
	"io"
	"os"
	"strings"
	"time"

	"github.com/ezeoleaf/larry/cache"
)

const noExpiration = time.Duration(0)

// Client represent the repositories
type Client interface {
	Load(blacklistFileName, cacheKeyPrefix string) error
	LoadFromReader(handle io.Reader, keyPrefix string) error
}

// blacklist represent the blacklist  model
type blacklistClient struct {
	CacheClient cache.Client
}

// NewClient will create an object that represent the Blacklist interface
func NewClient(cacheClient cache.Client) Client {
	return &blacklistClient{
		CacheClient: cacheClient,
	}
}

func (bc *blacklistClient) Load(blacklistFileName, cacheKeyPrefix string) error {
	keyPrefix := "blacklist-" + cacheKeyPrefix

	if err := bc.clear(keyPrefix); err != nil {
		return err
	}

	if blacklistFileName != "" {
		f, err := os.OpenFile(blacklistFileName, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return err
		}
		defer f.Close()

		return bc.LoadFromReader(f, keyPrefix)
	}

	return nil
}

func (bc *blacklistClient) LoadFromReader(handle io.Reader, keyPrefix string) error {
	sc := bufio.NewScanner(handle)
	for sc.Scan() {
		parts := strings.Split(sc.Text(), "#")
		repoId := strings.TrimSpace(parts[0])
		if repoId != "" {
			if err := bc.CacheClient.Set(keyPrefix+repoId, true, noExpiration); err != nil {
				return err
			}
		}
	}
	if err := sc.Err(); err != nil {
		return err
	}

	return nil
}

func (bc *blacklistClient) clear(keyPrefix string) error {
	deleteKeyFn := func(ctx context.Context, key string) error {
		if err := bc.CacheClient.Del(key); err != nil {
			return err
		}
		return nil
	}

	err := bc.CacheClient.Scan(keyPrefix+"*", deleteKeyFn)
	return err
}
