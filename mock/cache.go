package mock

import "time"

type CacheClientMock struct {
	SetFn func(key string, value interface{}, exp time.Duration) error
	GetFn func(key string) (string, error)
	DelFn func(key string) error
}

func (c CacheClientMock) Set(key string, value interface{}, exp time.Duration) error {
	if c.SetFn == nil {
		return nil
	}

	return c.SetFn(key, value, exp)
}

func (c CacheClientMock) Get(key string) (string, error) {
	if c.GetFn == nil {
		return "", nil
	}

	return c.GetFn(key)
}

func (c CacheClientMock) Del(key string) error {
	if c.DelFn == nil {
		return nil
	}

	return c.DelFn(key)
}
