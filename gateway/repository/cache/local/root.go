package local

import (
	"context"
	"github.com/allegro/bigcache/v3"
	"go.uber.org/atomic"
	"time"
)

type cacheLocal struct {
	cache     *bigcache.BigCache
	isRunning *atomic.Bool
}

type CacheLocalImpl interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	IsRunning() bool
	Run()
	Stop()
	Reset() error
}

func NewLocalCache(ttl time.Duration) (CacheLocalImpl, error) {
	cache, e := bigcache.New(context.Background(), bigcache.DefaultConfig(ttl))
	if e != nil {
		return nil, e
	}
	return &cacheLocal{cache: cache, isRunning: atomic.NewBool(false)}, nil
}

func (l *cacheLocal) Set(key string, value []byte) error {
	return l.cache.Set(key, value)
}

func (l *cacheLocal) Get(key string) ([]byte, error) {
	result, e := l.cache.Get(key)
	if e != nil {
		return nil, e
	}
	return result, nil
}

func (l *cacheLocal) IsRunning() bool {
	return l.isRunning.Load()
}

func (l *cacheLocal) Run() {
	l.isRunning.Store(true)
}

func (l *cacheLocal) Stop() {
	l.isRunning.Store(false)
}

func (l *cacheLocal) Reset() error {
	return l.cache.Reset()
}
