package cache

import (
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redis/v8"
)

func TestGetKey(t *testing.T) {
	mr, _ := miniredis.Run()
	ro := &redis.Options{
		Addr: mr.Addr(),
	}

	r := NewRedisRepository(ro)

	_, e := r.Get("key")

	if e != redis.Nil {
		t.Error("No value expected for key not saved in cache")
	}

	r.Set("key", "v", 1)

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

	r := NewRedisRepository(ro)

	r.Set("key", "v", 1)

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

	r := NewRedisRepository(ro)

	r.Set("key", "v", 1)

	r.Del("key")

	_, e := r.Get("key")

	if e != redis.Nil {
		t.Error("Value expected for key")
	}

}
