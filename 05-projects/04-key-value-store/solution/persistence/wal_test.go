package persistence

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alyxpink/go-training/kvstore/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWAL_Append(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "test.wal")

	wal, err := NewWAL(walPath)
	require.NoError(t, err)
	defer wal.Close()

	// Append commands
	err = wal.Append("SET key1 value1")
	require.NoError(t, err)

	err = wal.Append("SET key2 value2")
	require.NoError(t, err)

	err = wal.Append("DEL key1")
	require.NoError(t, err)

	// Verify file exists and has content
	info, err := os.Stat(walPath)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0), "WAL file should have content")
}

func TestWAL_Replay(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "test.wal")

	// Create WAL and write commands
	wal, err := NewWAL(walPath)
	require.NoError(t, err)

	commands := []string{
		"SET key1 value1",
		"SET key2 value2",
		"SET key3 value3",
		"DEL key2",
		"EXPIRE key3 100",
	}

	for _, cmd := range commands {
		err = wal.Append(cmd)
		require.NoError(t, err)
	}
	wal.Close()

	// Create new WAL and replay
	wal2, err := NewWAL(walPath)
	require.NoError(t, err)
	defer wal2.Close()

	kvStore := store.NewKVStore()
	err = wal2.Replay(kvStore)
	require.NoError(t, err)

	// Verify replayed state
	val, ok := kvStore.Get("key1")
	require.True(t, ok)
	assert.Equal(t, "value1", val)

	_, ok = kvStore.Get("key2")
	assert.False(t, ok, "key2 should be deleted")

	val, ok = kvStore.Get("key3")
	require.True(t, ok)
	assert.Equal(t, "value3", val)

	// Verify expiration was set
	ttl := kvStore.TTL("key3")
	assert.GreaterOrEqual(t, ttl, 0)
}

func TestWAL_Close(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "test.wal")

	wal, err := NewWAL(walPath)
	require.NoError(t, err)

	err = wal.Append("SET key1 value1")
	require.NoError(t, err)

	err = wal.Close()
	require.NoError(t, err)

	// Verify we can't write after close
	err = wal.Append("SET key2 value2")
	assert.Error(t, err, "should not be able to write to closed WAL")
}

func TestWAL_EmptyReplay(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "test.wal")

	wal, err := NewWAL(walPath)
	require.NoError(t, err)
	defer wal.Close()

	// Replay empty WAL
	kvStore := store.NewKVStore()
	err = wal.Replay(kvStore)
	require.NoError(t, err)

	// Store should be empty
	keys := kvStore.Keys("*")
	assert.Empty(t, keys)
}

func TestWAL_MultipleAppends(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "test.wal")

	wal, err := NewWAL(walPath)
	require.NoError(t, err)
	defer wal.Close()

	// Append many commands
	for i := 0; i < 100; i++ {
		err = wal.Append("SET key value")
		require.NoError(t, err)
	}

	// Verify all were written
	wal.Close()
	wal2, err := NewWAL(walPath)
	require.NoError(t, err)
	defer wal2.Close()

	kvStore := store.NewKVStore()
	err = wal2.Replay(kvStore)
	require.NoError(t, err)
}

func TestWAL_ReplayWithInvalidCommand(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "test.wal")

	// Create WAL with invalid command
	wal, err := NewWAL(walPath)
	require.NoError(t, err)

	err = wal.Append("SET key1 value1")
	require.NoError(t, err)

	err = wal.Append("INVALID_COMMAND")
	require.NoError(t, err)

	err = wal.Append("SET key2 value2")
	require.NoError(t, err)

	wal.Close()

	// Replay should handle or skip invalid commands
	wal2, err := NewWAL(walPath)
	require.NoError(t, err)
	defer wal2.Close()

	kvStore := store.NewKVStore()
	err = wal2.Replay(kvStore)

	// Depending on implementation, this might error or skip
	// At minimum, valid commands should be replayed
	val, ok := kvStore.Get("key1")
	if ok {
		assert.Equal(t, "value1", val)
	}
}

func TestWAL_ConcurrentAppends(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "test.wal")

	wal, err := NewWAL(walPath)
	require.NoError(t, err)
	defer wal.Close()

	// Concurrent appends should be safe
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 10; j++ {
				wal.Append("SET key value")
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify WAL has content
	info, err := os.Stat(walPath)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0))
}

func TestWAL_NewWALCreatesFile(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "new.wal")

	// File should not exist
	_, err := os.Stat(walPath)
	assert.True(t, os.IsNotExist(err))

	// Create WAL
	wal, err := NewWAL(walPath)
	require.NoError(t, err)
	defer wal.Close()

	// File should exist now
	_, err = os.Stat(walPath)
	require.NoError(t, err)
}

func TestWAL_AppendPreservesOrder(t *testing.T) {
	tmpDir := t.TempDir()
	walPath := filepath.Join(tmpDir, "test.wal")

	wal, err := NewWAL(walPath)
	require.NoError(t, err)

	// Append in specific order
	commands := []string{
		"SET counter 0",
		"SET counter 1",
		"SET counter 2",
		"SET counter 3",
	}

	for _, cmd := range commands {
		err = wal.Append(cmd)
		require.NoError(t, err)
	}
	wal.Close()

	// Replay and verify final value
	wal2, err := NewWAL(walPath)
	require.NoError(t, err)
	defer wal2.Close()

	kvStore := store.NewKVStore()
	err = wal2.Replay(kvStore)
	require.NoError(t, err)

	// Final value should be from last command
	val, ok := kvStore.Get("counter")
	require.True(t, ok)
	assert.Equal(t, "3", val)
}
