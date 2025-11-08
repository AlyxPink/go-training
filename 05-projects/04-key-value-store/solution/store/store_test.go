package store

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKVStore_SetGet(t *testing.T) {
	s := NewKVStore()

	// Test basic set/get
	s.Set("key1", "value1")
	val, ok := s.Get("key1")
	require.True(t, ok, "key should exist")
	assert.Equal(t, "value1", val)

	// Test get non-existent key
	_, ok = s.Get("nonexistent")
	assert.False(t, ok, "key should not exist")
}

func TestKVStore_Del(t *testing.T) {
	s := NewKVStore()

	// Set and delete
	s.Set("key1", "value1")
	deleted := s.Del("key1")
	assert.True(t, deleted, "delete should return true for existing key")

	// Verify deletion
	_, ok := s.Get("key1")
	assert.False(t, ok, "key should not exist after deletion")

	// Delete non-existent key
	deleted = s.Del("nonexistent")
	assert.False(t, deleted, "delete should return false for non-existent key")
}

func TestKVStore_Exists(t *testing.T) {
	s := NewKVStore()

	// Key doesn't exist
	assert.False(t, s.Exists("key1"))

	// Key exists
	s.Set("key1", "value1")
	assert.True(t, s.Exists("key1"))

	// Key deleted
	s.Del("key1")
	assert.False(t, s.Exists("key1"))
}

func TestKVStore_Keys(t *testing.T) {
	s := NewKVStore()

	// Setup data
	s.Set("user:1", "alice")
	s.Set("user:2", "bob")
	s.Set("post:1", "hello")
	s.Set("post:2", "world")

	tests := []struct {
		name     string
		pattern  string
		expected []string
	}{
		{
			name:     "all keys",
			pattern:  "*",
			expected: []string{"user:1", "user:2", "post:1", "post:2"},
		},
		{
			name:     "user keys",
			pattern:  "user:*",
			expected: []string{"user:1", "user:2"},
		},
		{
			name:     "post keys",
			pattern:  "post:*",
			expected: []string{"post:1", "post:2"},
		},
		{
			name:     "exact match",
			pattern:  "user:1",
			expected: []string{"user:1"},
		},
		{
			name:     "no match",
			pattern:  "admin:*",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys := s.Keys(tt.pattern)
			assert.ElementsMatch(t, tt.expected, keys)
		})
	}
}

func TestKVStore_Expire(t *testing.T) {
	s := NewKVStore()

	// Set key and expiration
	s.Set("key1", "value1")
	success := s.Expire("key1", 1)
	assert.True(t, success, "expire should succeed for existing key")

	// Key should still exist
	val, ok := s.Get("key1")
	require.True(t, ok)
	assert.Equal(t, "value1", val)

	// Wait for expiration
	time.Sleep(1100 * time.Millisecond)

	// Key should be expired
	_, ok = s.Get("key1")
	assert.False(t, ok, "key should be expired")

	// Expire non-existent key
	success = s.Expire("nonexistent", 10)
	assert.False(t, success, "expire should fail for non-existent key")
}

func TestKVStore_TTL(t *testing.T) {
	s := NewKVStore()

	// TTL for non-existent key
	ttl := s.TTL("nonexistent")
	assert.Equal(t, -2, ttl, "TTL should be -2 for non-existent key")

	// TTL for key without expiration
	s.Set("key1", "value1")
	ttl = s.TTL("key1")
	assert.Equal(t, -1, ttl, "TTL should be -1 for key without expiration")

	// TTL for key with expiration
	s.Expire("key1", 10)
	ttl = s.TTL("key1")
	assert.GreaterOrEqual(t, ttl, 8, "TTL should be at least 8 seconds")
	assert.LessOrEqual(t, ttl, 10, "TTL should be at most 10 seconds")
}

func TestKVStore_Concurrency(t *testing.T) {
	s := NewKVStore()
	var wg sync.WaitGroup
	numGoroutines := 100
	numOps := 100

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				key := "key" + string(rune(id))
				s.Set(key, "value")
			}
		}(i)
	}

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				key := "key" + string(rune(id))
				s.Get(key)
			}
		}(i)
	}

	// Concurrent deletes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOps; j++ {
				key := "key" + string(rune(id))
				s.Del(key)
			}
		}(i)
	}

	wg.Wait()

	// No assertion needed - test passes if no race condition detected
}

func TestKVStore_ExpiredKeysNotInKeys(t *testing.T) {
	s := NewKVStore()

	s.Set("key1", "value1")
	s.Set("key2", "value2")
	s.Expire("key1", 1)

	// Wait for expiration
	time.Sleep(1100 * time.Millisecond)

	// Keys should not include expired key
	keys := s.Keys("*")
	assert.NotContains(t, keys, "key1")
	assert.Contains(t, keys, "key2")
}

func TestKVStore_UpdateValue(t *testing.T) {
	s := NewKVStore()

	// Set initial value
	s.Set("key1", "value1")
	val, _ := s.Get("key1")
	assert.Equal(t, "value1", val)

	// Update value
	s.Set("key1", "value2")
	val, _ = s.Get("key1")
	assert.Equal(t, "value2", val)
}

func TestKVStore_ExpireUpdatesExisting(t *testing.T) {
	s := NewKVStore()

	s.Set("key1", "value1")
	s.Expire("key1", 10)

	// Update expiration
	success := s.Expire("key1", 20)
	assert.True(t, success)

	ttl := s.TTL("key1")
	assert.GreaterOrEqual(t, ttl, 18)
	assert.LessOrEqual(t, ttl, 20)
}

func TestMatchPattern(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		pattern string
		match   bool
	}{
		{"wildcard all", "anything", "*", true},
		{"exact match", "key1", "key1", true},
		{"exact no match", "key1", "key2", false},
		{"prefix match", "user:123", "user:*", true},
		{"prefix no match", "post:123", "user:*", false},
		{"suffix match", "123:user", "*:user", true},
		{"suffix no match", "123:post", "*:user", false},
		{"both ends match", "pre:middle:suf", "pre:*:suf", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchPattern(tt.key, tt.pattern)
			assert.Equal(t, tt.match, result)
		})
	}
}
