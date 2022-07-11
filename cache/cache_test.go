package cache

import (
	"context"
	"errors"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
)

func TestGetKey(t *testing.T) {
	mr, _ := miniredis.Run()
	ro := &redis.Options{
		Addr: mr.Addr(),
	}

	r := NewClient(ro)

	_, e := r.Get("key")

	if e != redis.Nil {
		t.Error("No value expected for key not saved in cache")
	}

	err := r.Set("key", "v", 1)
	if err != nil {
		t.Error("could not set key")
	}

	_, e = r.Get("key")

	if e != nil {
		t.Error("Value expected for key")
	}
}

func TestSetKey(t *testing.T) {
	mr, _ := miniredis.Run()
	ro := &redis.Options{
		Addr: mr.Addr(),
	}

	r := NewClient(ro)

	err := r.Set("key", "v", 1)
	if err != nil {
		t.Error("could not set key")
	}

	v, e := r.Get("key")

	if e != nil {
		t.Error("Value expected for key")
	}

	if v != "v" {
		t.Errorf("v expected got %s", v)
	}
}

func TestDetKey(t *testing.T) {
	mr, _ := miniredis.Run()
	ro := &redis.Options{
		Addr: mr.Addr(),
	}

	r := NewClient(ro)

	err := r.Set("key", "v", 1)
	if err != nil {
		t.Error("could not set key")
	}

	err = r.Del("key")
	if err != nil {
		t.Error("could not del key")
	}

	_, e := r.Get("key")

	if e != redis.Nil {
		t.Error("Value expected for key")
	}
}

func TestScanKey(t *testing.T) {
	mr, _ := miniredis.Run()
	ro := &redis.Options{
		Addr: mr.Addr(),
	}

	r := NewClient(ro)

	err := r.Set("key", "v", 1)
	if err != nil {
		t.Error("could not set key")
	}

	err = r.Scan("key", func(ctx context.Context, key string) error {
		return errors.New("some error")
	})

	if err == nil {
		t.Error("expected error but got none")
	}

	err = r.Set("key", "v", 1)
	if err != nil {
		t.Error("could not set key")
	}

	err = r.Scan("key", func(ctx context.Context, key string) error {
		return nil
	})

	if err != nil {
		t.Errorf("expected no error but got %v", err)
	}

}
