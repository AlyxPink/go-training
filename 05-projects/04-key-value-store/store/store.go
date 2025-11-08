package store

import (
	"strings"
	"sync"
	"time"
)

type Entry struct {
	Value     string
	ExpiresAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type KVStore struct {
	data map[string]*Entry
	mu   sync.RWMutex
}

func NewKVStore() *KVStore {
	return &KVStore{
		data: make(map[string]*Entry),
	}
}

func (s *KVStore) Set(key, value string) {
	// TODO: Implement thread-safe set
	panic("not implemented")
}

func (s *KVStore) Get(key string) (string, bool) {
	// TODO: Implement thread-safe get with expiry check
	panic("not implemented")
}

func (s *KVStore) Del(key string) bool {
	// TODO: Implement thread-safe delete
	panic("not implemented")
}

func (s *KVStore) Exists(key string) bool {
	// TODO: Implement exists check with expiry
	panic("not implemented")
}

func (s *KVStore) Keys(pattern string) []string {
	// TODO: Implement pattern matching
	panic("not implemented")
}

func (s *KVStore) Expire(key string, seconds int) bool {
	// TODO: Implement expiration
	panic("not implemented")
}

func (s *KVStore) TTL(key string) int {
	// TODO: Implement TTL check
	panic("not implemented")
}

func matchPattern(key, pattern string) bool {
	// Simple glob matching (* only)
	if pattern == "*" {
		return true
	}

	if !strings.Contains(pattern, "*") {
		return key == pattern
	}

	parts := strings.Split(pattern, "*")
	if len(parts) == 2 {
		return strings.HasPrefix(key, parts[0]) && strings.HasSuffix(key, parts[1])
	}

	return false
}
