package rateLimit

import (
	"context"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	Config "github.com/Dercraker/SearchEngine/internal/crawler/config"
)

type HostLimiter struct {
	mu          sync.Mutex
	nextAllowed time.Time
}

type Limiter struct {
	Cfg   Config.LimitConfig
	Count atomic.Int64

	globalMu   sync.Mutex
	globalNext time.Time

	mu  sync.Mutex
	lru *lruCache

	rngMu sync.Mutex
	rng   *rand.Rand
}

func NewRateLimiter(cfg Config.LimitConfig) *Limiter {
	seed := time.Now().UnixNano()
	return &Limiter{
		Cfg: cfg,
		lru: newLRUCache(cfg.MaxHost),
		rng: rand.New(rand.NewSource(seed)),
	}
}

func (l *Limiter) WaitGlobal(ctx context.Context) error {
	l.globalMu.Lock()
	defer l.globalMu.Unlock()

	now := time.Now()

	if !l.globalNext.IsZero() && now.Before(l.globalNext) {
		wait := l.globalNext.Sub(now)
		l.globalMu.Unlock()

		select {
		case <-time.After(wait):
		case <-ctx.Done():
			l.globalMu.Lock()
			return ctx.Err()
		}

		l.globalMu.Lock()
		now = time.Now()
	}

	delay := l.Cfg.HostDelay
	delay += l.randDuration(l.Cfg.Jitter)
	l.globalNext = now.Add(delay)

	return nil
}

func (l *Limiter) WaitHost(ctx context.Context, hl *HostLimiter) error {
	hl.mu.Lock()
	defer hl.mu.Unlock()

	now := time.Now()

	if !hl.nextAllowed.IsZero() && now.Before(hl.nextAllowed) {
		wait := hl.nextAllowed.Sub(now)
		hl.mu.Unlock()

		select {
		case <-time.After(wait):
		case <-ctx.Done():
			hl.mu.Lock()
			return ctx.Err()
		}

		hl.mu.Lock()
		now = time.Now()
	}

	delay := l.Cfg.HostDelay
	delay += l.randDuration(l.Cfg.Jitter)
	hl.nextAllowed = now.Add(delay)

	return nil
}

func (l *Limiter) GetHostLimiter(host string) *HostLimiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	if el, ok := l.lru.items[host]; ok {
		l.lru.ll.MoveToFront(el)
		return el.Value.(*lruEntry).hl
	}

	hl := &HostLimiter{}
	ent := &lruEntry{host: host, hl: hl}
	el := l.lru.ll.PushFront(ent)
	l.lru.items[host] = el

	if l.lru.ll.Len() > l.lru.max {
		last := l.lru.ll.Back()
		if last != nil {
			lastEntry := last.Value.(*lruEntry)
			delete(l.lru.items, lastEntry.host)
			l.lru.ll.Remove(last)
		}
	}

	return hl
}

func (l *Limiter) randDuration(max time.Duration) time.Duration {
	if max <= 0 {
		return 0
	}
	l.rngMu.Lock()
	n := l.rng.Int63n(int64(max) + 1)
	l.rngMu.Unlock()

	return time.Duration(n)
}

func NormalizeHost(host string) string {
	host = strings.ToLower(strings.TrimSpace(host))

	if i := strings.LastIndex(host, ":"); i >= 0 && strings.Contains(host[i+1:i], "0") || strings.ContainsAny(host[i+1:], "0123456789") {
		hostOnly := host[:i]
		portPart := host[i+1:]
		if portPart != "" && isAllDigits(portPart) {
			return hostOnly
		}
	}
	return host
}

func isAllDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return s != ""
}
