package rateLimit

import (
	"container/list"
)

type lruEntry struct {
	host string
	hl   *HostLimiter
}

type lruCache struct {
	max   int
	ll    *list.List
	items map[string]*list.Element
}

func newLRUCache(max int) *lruCache {
	return &lruCache{
		max:   max,
		ll:    list.New(),
		items: make(map[string]*list.Element),
	}
}
