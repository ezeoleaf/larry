package mock

import (
	"context"
	"time"
)

// CacheClientMock is a mock of CacheClient
type CacheClientMock struct {
	SetFn  func(key string, value interface{}, exp time.Duration) error
	GetFn  func(key string) (string, error)
	DelFn  func(key string) error
	ScanFn func(key string, action func(context.Context, string) error) error
}

// Set calls SetFn
func (c CacheClientMock) Set(key string, value interface{}, exp time.Duration) error {
	if c.SetFn == nil {
		return nil
	}

	return c.SetFn(key, value, exp)
}

// Get calls GetFn
func (c CacheClientMock) Get(key string) (string, error) {
	if c.GetFn == nil {
		return "", nil
	}

	return c.GetFn(key)
}

// Del calls DelFn
func (c CacheClientMock) Del(key string) error {
	if c.DelFn == nil {
		return nil
	}

	return c.DelFn(key)
}

// Scan calls ScanFn
func (c CacheClientMock) Scan(key string, action func(context.Context, string) error) error {
	if c.ScanFn == nil {
		return nil
	}

	return c.ScanFn(key, action)
}
