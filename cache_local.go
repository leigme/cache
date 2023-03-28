package cache

import (
	"github.com/allegro/bigcache"
)

type localCache struct {
	Options
	lc *bigcache.BigCache
}

func (l *localCache) Set(key string, value []byte) (ok bool) {
	if err := l.lc.Set(key, value); err == nil {
		ok = true
	}
	return
}

func (l *localCache) Get(key string) (value []byte) {
	if data, err := l.lc.Get(key); err == nil {
		value = data
	}
	return
}

func (l *localCache) Remove(key string) (ok bool) {
	if err := l.lc.Delete(key); err == nil {
		ok = true
	}
	return
}

func NewLocalCache(opts ...Option) (Cache, error) {
	dos := defaultOptions()
	for _, apply := range opts {
		apply(&dos)
	}
	bc, err := bigcache.NewBigCache(bigcache.DefaultConfig(dos.timeout))
	if err != nil {
		return nil, err
	}
	return &localCache{
		dos,
		bc,
	}, nil
}
