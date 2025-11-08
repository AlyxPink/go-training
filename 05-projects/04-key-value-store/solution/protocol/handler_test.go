package protocol

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alyxpink/go-training/kvstore/persistence"
	"github.com/alyxpink/go-training/kvstore/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_SET(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "test.wal")

	kvStore := store.NewKVStore()
	wal, err := persistence.NewWAL(walPath)
	require.NoError(t, err)
	defer wal.Close()

	handler := NewHandler(kvStore, wal)

	// Test SET command
	resp := handler.Handle("SET key1 value1")
	assert.Equal(t, "+OK", resp)

	// Verify value was set
	val, ok := kvStore.Get("key1")
	require.True(t, ok)
	assert.Equal(t, "value1", val)
}

func TestHandler_GET(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "test.wal")

	kvStore := store.NewKVStore()
	wal, err := persistence.NewWAL(walPath)
	require.NoError(t, err)
	defer wal.Close()

	handler := NewHandler(kvStore, wal)

	// Set value first
	kvStore.Set("key1", "value1")

	// Test GET command
	resp := handler.Handle("GET key1")
	assert.Equal(t, "$6\r\nvalue1", resp)

	// Test GET non-existent key
	resp = handler.Handle("GET nonexistent")
	assert.Equal(t, "$-1", resp)
}

func TestHandler_DEL(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "test.wal")

	kvStore := store.NewKVStore()
	wal, err := persistence.NewWAL(walPath)
	require.NoError(t, err)
	defer wal.Close()

	handler := NewHandler(kvStore, wal)

	// Set value
	kvStore.Set("key1", "value1")

	// Test DEL command
	resp := handler.Handle("DEL key1")
	assert.Equal(t, ":1", resp)

	// Verify deletion
	_, ok := kvStore.Get("key1")
	assert.False(t, ok)

	// Test DEL non-existent key
	resp = handler.Handle("DEL nonexistent")
	assert.Equal(t, ":0", resp)
}

func TestHandler_EXISTS(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "test.wal")

	kvStore := store.NewKVStore()
	wal, err := persistence.NewWAL(walPath)
	require.NoError(t, err)
	defer wal.Close()

	handler := NewHandler(kvStore, wal)

	// Test EXISTS for non-existent key
	resp := handler.Handle("EXISTS key1")
	assert.Equal(t, ":0", resp)

	// Set value
	kvStore.Set("key1", "value1")

	// Test EXISTS for existing key
	resp = handler.Handle("EXISTS key1")
	assert.Equal(t, ":1", resp)
}

func TestHandler_KEYS(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "test.wal")

	kvStore := store.NewKVStore()
	wal, err := persistence.NewWAL(walPath)
	require.NoError(t, err)
	defer wal.Close()

	handler := NewHandler(kvStore, wal)

	// Set multiple keys
	kvStore.Set("user:1", "alice")
	kvStore.Set("user:2", "bob")
	kvStore.Set("post:1", "hello")

	// Test KEYS command
	resp := handler.Handle("KEYS user:*")
	assert.Contains(t, resp, "user:1")
	assert.Contains(t, resp, "user:2")
	assert.NotContains(t, resp, "post:1")
}

func TestHandler_EXPIRE(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "test.wal")

	kvStore := store.NewKVStore()
	wal, err := persistence.NewWAL(walPath)
	require.NoError(t, err)
	defer wal.Close()

	handler := NewHandler(kvStore, wal)

	// Set value
	kvStore.Set("key1", "value1")

	// Test EXPIRE command
	resp := handler.Handle("EXPIRE key1 10")
	assert.Equal(t, ":1", resp)

	// Test EXPIRE non-existent key
	resp = handler.Handle("EXPIRE nonexistent 10")
	assert.Equal(t, ":0", resp)
}

func TestHandler_TTL(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "test.wal")

	kvStore := store.NewKVStore()
	wal, err := persistence.NewWAL(walPath)
	require.NoError(t, err)
	defer wal.Close()

	handler := NewHandler(kvStore, wal)

	// Test TTL for non-existent key
	resp := handler.Handle("TTL nonexistent")
	assert.Equal(t, ":-2", resp)

	// Set value without expiration
	kvStore.Set("key1", "value1")
	resp = handler.Handle("TTL key1")
	assert.Equal(t, ":-1", resp)

	// Set expiration
	kvStore.Expire("key1", 10)
	resp = handler.Handle("TTL key1")
	assert.Contains(t, resp, ":")
	assert.NotEqual(t, ":-1", resp)
	assert.NotEqual(t, ":-2", resp)
}

func TestHandler_InvalidCommand(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "test.wal")

	kvStore := store.NewKVStore()
	wal, err := persistence.NewWAL(walPath)
	require.NoError(t, err)
	defer wal.Close()

	handler := NewHandler(kvStore, wal)

	resp := handler.Handle("INVALID command")
	assert.Contains(t, resp, "-ERR")
}

func TestHandler_SETWithSpaces(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "test.wal")

	kvStore := store.NewKVStore()
	wal, err := persistence.NewWAL(walPath)
	require.NoError(t, err)
	defer wal.Close()

	handler := NewHandler(kvStore, wal)

	// Test SET with value containing spaces
	resp := handler.Handle("SET key1 hello world")
	assert.Equal(t, "+OK", resp)

	val, ok := kvStore.Get("key1")
	require.True(t, ok)
	// Value should contain spaces
	assert.Contains(t, val, " ")
}

func TestHandler_CaseInsensitiveCommands(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "test.wal")

	kvStore := store.NewKVStore()
	wal, err := persistence.NewWAL(walPath)
	require.NoError(t, err)
	defer wal.Close()

	handler := NewHandler(kvStore, wal)

	// Commands should be case insensitive
	tests := []string{
		"set key1 value1",
		"SET key2 value2",
		"Set key3 value3",
	}

	for _, cmd := range tests {
		resp := handler.Handle(cmd)
		assert.Equal(t, "+OK", resp)
	}
}

func TestHandler_EmptyCommand(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "test.wal")

	kvStore := store.NewKVStore()
	wal, err := persistence.NewWAL(walPath)
	require.NoError(t, err)
	defer wal.Close()

	handler := NewHandler(kvStore, wal)

	resp := handler.Handle("")
	assert.Contains(t, resp, "-ERR")
}

func TestHandler_MissingArguments(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "test.wal")

	kvStore := store.NewKVStore()
	wal, err := persistence.NewWAL(walPath)
	require.NoError(t, err)
	defer wal.Close()

	handler := NewHandler(kvStore, wal)

	tests := []string{
		"SET",
		"SET key1",
		"GET",
		"DEL",
		"EXISTS",
		"EXPIRE",
		"EXPIRE key1",
	}

	for _, cmd := range tests {
		resp := handler.Handle(cmd)
		assert.Contains(t, resp, "-ERR", "command '%s' should return error", cmd)
	}
}

func TestHandler_WALPersistence(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "test.wal")

	kvStore := store.NewKVStore()
	wal, err := persistence.NewWAL(walPath)
	require.NoError(t, err)

	handler := NewHandler(kvStore, wal)

	// Execute commands
	handler.Handle("SET key1 value1")
	handler.Handle("SET key2 value2")
	handler.Handle("DEL key1")

	wal.Close()

	// Verify WAL contains commands
	info, err := os.Stat(walPath)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0), "WAL should contain data")

	// Replay WAL
	newStore := store.NewKVStore()
	newWAL, err := persistence.NewWAL(walPath)
	require.NoError(t, err)
	defer newWAL.Close()

	err = newWAL.Replay(newStore)
	require.NoError(t, err)

	// Verify replayed state
	_, ok := newStore.Get("key1")
	assert.False(t, ok, "key1 should be deleted")

	val, ok := newStore.Get("key2")
	require.True(t, ok)
	assert.Equal(t, "value2", val)
}

func TestHandler_MultipleCommands(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "test.wal")

	kvStore := store.NewKVStore()
	wal, err := persistence.NewWAL(walPath)
	require.NoError(t, err)
	defer wal.Close()

	handler := NewHandler(kvStore, wal)

	commands := []struct {
		cmd      string
		expected string
	}{
		{"SET key1 value1", "+OK"},
		{"GET key1", "$6\r\nvalue1"},
		{"EXISTS key1", ":1"},
		{"SET key2 value2", "+OK"},
		{"KEYS *", "*2"},
		{"DEL key1", ":1"},
		{"EXISTS key1", ":0"},
	}

	for _, tc := range commands {
		resp := handler.Handle(tc.cmd)
		assert.Contains(t, resp, tc.expected, "command: %s", tc.cmd)
	}
}

func TestHandler_KEYSPattern(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "test.wal")

	kvStore := store.NewKVStore()
	wal, err := persistence.NewWAL(walPath)
	require.NoError(t, err)
	defer wal.Close()

	handler := NewHandler(kvStore, wal)

	// Setup data
	handler.Handle("SET user:1 alice")
	handler.Handle("SET user:2 bob")
	handler.Handle("SET post:1 hello")

	// Test different patterns
	tests := []struct {
		pattern  string
		contains []string
		excludes []string
	}{
		{"*", []string{"user:1", "user:2", "post:1"}, nil},
		{"user:*", []string{"user:1", "user:2"}, []string{"post:1"}},
		{"post:*", []string{"post:1"}, []string{"user:1", "user:2"}},
	}

	for _, tc := range tests {
		resp := handler.Handle("KEYS " + tc.pattern)
		for _, key := range tc.contains {
			assert.Contains(t, resp, key, "pattern: %s", tc.pattern)
		}
		for _, key := range tc.excludes {
			assert.NotContains(t, resp, key, "pattern: %s", tc.pattern)
		}
	}
}
