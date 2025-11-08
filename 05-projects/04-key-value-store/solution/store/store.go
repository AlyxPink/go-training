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
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	if entry, exists := s.data[key]; exists {
		entry.Value = value
		entry.UpdatedAt = now
	} else {
		s.data[key] = &Entry{
			Value:     value,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
}

func (s *KVStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.data[key]
	if !exists {
		return "", false
	}

	// Check if key has expired
	if entry.ExpiresAt != nil && time.Now().After(*entry.ExpiresAt) {
		return "", false
	}

	return entry.Value, true
}

func (s *KVStore) Del(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[key]; !exists {
		return false
	}

	delete(s.data, key)
	return true
}

func (s *KVStore) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.data[key]
	if !exists {
		return false
	}

	// Check if key has expired
	if entry.ExpiresAt != nil && time.Now().After(*entry.ExpiresAt) {
		return false
	}

	return true
}

func (s *KVStore) Keys(pattern string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var keys []string
	now := time.Now()

	for key, entry := range s.data {
		// Skip expired keys
		if entry.ExpiresAt != nil && now.After(*entry.ExpiresAt) {
			continue
		}

		if matchPattern(key, pattern) {
			keys = append(keys, key)
		}
	}

	return keys
}

func (s *KVStore) Expire(key string, seconds int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, exists := s.data[key]
	if !exists {
		return false
	}

	// Check if already expired
	if entry.ExpiresAt != nil && time.Now().After(*entry.ExpiresAt) {
		return false
	}

	expiresAt := time.Now().Add(time.Duration(seconds) * time.Second)
	entry.ExpiresAt = &expiresAt
	return true
}

func (s *KVStore) TTL(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry, exists := s.data[key]
	if !exists {
		return -2
	}

	// Check if already expired
	if entry.ExpiresAt != nil && time.Now().After(*entry.ExpiresAt) {
		return -2
	}

	if entry.ExpiresAt == nil {
		return -1
	}

	ttl := time.Until(*entry.ExpiresAt).Seconds()
	if ttl < 0 {
		return -2
	}

	return int(ttl)
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
