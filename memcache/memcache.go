package memcache

import (
	"github.com/bradfitz/gomemcache/memcache"
	"strings"
)

type Memcache struct {
	cache *memcache.Client
}

func New(hosts []string) Memcache {
	for i, host := range hosts {
		hosts[i] = strings.TrimSpace(host)
	}

	return Memcache{
		cache: memcache.New(hosts...),
	}
}

func (m Memcache) Get(key string) []byte {
	item, _ := m.cache.Get(key)

	return item.Value
}

func (m Memcache) Put(key string, data []byte) error {
	return m.cache.Set(memcache.Item{
		Key: key,
		Value: data,
	})
}
