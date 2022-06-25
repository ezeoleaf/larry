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

var noExpiration = time.Duration(0)

func Load(cc cache.Client, blacklistFileName, cacheKeyPrefix string) error {
	keyPrefix := "blacklist-" + cacheKeyPrefix

	if err := clear(cc, keyPrefix); err != nil {
		return err
	}

	if blacklistFileName != "" {
		f, err := os.OpenFile(blacklistFileName, os.O_RDONLY, os.ModePerm)
		if err != nil {
			return err
		}
		defer f.Close()

		return readerLoader(cc, f, keyPrefix)
	}

	return nil
}

func readerLoader(cc cache.Client, handle io.Reader, keyPrefix string) error {
	sc := bufio.NewScanner(handle)
	for sc.Scan() {
		parts := strings.Split(sc.Text(), "#")
		repoId := strings.TrimSpace(parts[0])
		if repoId != "" {
			if err := cc.Set(keyPrefix+repoId, "1", noExpiration); err != nil {
				return err
			}
		}
	}
	if err := sc.Err(); err != nil {
		return err
	}

	return nil
}

func clear(cc cache.Client, keyPrefix string) error {
	deleteKeyFn := func(ctx context.Context, key string) error {
		if err := cc.Del(key); err != nil {
			return err
		}
		return nil
	}

	err := cc.Scan(keyPrefix+"*", deleteKeyFn)
	return err
}
