package cache

import (
	"sync"
	"time"
)

type Entry[T any] struct {
	Value                T
	FetchedAt, ExpiresAt time.Time
	Stale                bool
}
type Cache[T any] interface {
	Get(key string) (Entry[T], bool)
	Set(key string, value T, ttl time.Duration)
}
type Memory[T any] struct {
	mu    sync.RWMutex
	items map[string]Entry[T]
	now   func() time.Time
}

func NewMemory[T any]() *Memory[T] {
	return &Memory[T]{items: make(map[string]Entry[T]), now: time.Now}
}
func (m *Memory[T]) Get(key string) (Entry[T], bool) {
	m.mu.RLock()
	entry, ok := m.items[key]
	m.mu.RUnlock()
	if !ok {
		return Entry[T]{}, false
	}
	entry.Stale = !m.now().Before(entry.ExpiresAt)
	return entry, true
}
func (m *Memory[T]) Set(key string, value T, ttl time.Duration) {
	now := m.now()
	m.mu.Lock()
	m.items[key] = Entry[T]{Value: value, FetchedAt: now, ExpiresAt: now.Add(ttl)}
	m.mu.Unlock()
}
