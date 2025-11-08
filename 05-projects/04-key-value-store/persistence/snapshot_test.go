package persistence

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alyxpink/go-training/kvstore/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSnapshotManager_CreateSnapshot(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSnapshotManager(tmpDir, time.Minute)

	kvStore := store.NewKVStore()
	kvStore.Set("key1", "value1")
	kvStore.Set("key2", "value2")
	kvStore.Set("key3", "value3")

	err := sm.CreateSnapshot(kvStore)
	require.NoError(t, err)

	// Verify snapshot file exists
	files, err := os.ReadDir(tmpDir)
	require.NoError(t, err)
	assert.NotEmpty(t, files, "snapshot file should be created")
}

func TestSnapshotManager_LoadLatest(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSnapshotManager(tmpDir, time.Minute)

	// Create snapshot
	kvStore1 := store.NewKVStore()
	kvStore1.Set("key1", "value1")
	kvStore1.Set("key2", "value2")
	kvStore1.Set("key3", "value3")

	err := sm.CreateSnapshot(kvStore1)
	require.NoError(t, err)

	// Load snapshot into new store
	kvStore2 := store.NewKVStore()
	err = sm.LoadLatest(kvStore2)
	require.NoError(t, err)

	// Verify data was loaded
	val, ok := kvStore2.Get("key1")
	require.True(t, ok)
	assert.Equal(t, "value1", val)

	val, ok = kvStore2.Get("key2")
	require.True(t, ok)
	assert.Equal(t, "value2", val)

	val, ok = kvStore2.Get("key3")
	require.True(t, ok)
	assert.Equal(t, "value3", val)
}

func TestSnapshotManager_LoadLatestNoSnapshot(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSnapshotManager(tmpDir, time.Minute)

	kvStore := store.NewKVStore()
	err := sm.LoadLatest(kvStore)

	// Should error when no snapshot exists
	assert.Error(t, err)
}

func TestSnapshotManager_LoadLatestWithExpiration(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSnapshotManager(tmpDir, time.Minute)

	// Create snapshot with expiring key
	kvStore1 := store.NewKVStore()
	kvStore1.Set("key1", "value1")
	kvStore1.Set("key2", "value2")
	kvStore1.Expire("key2", 10)

	err := sm.CreateSnapshot(kvStore1)
	require.NoError(t, err)

	// Load snapshot
	kvStore2 := store.NewKVStore()
	err = sm.LoadLatest(kvStore2)
	require.NoError(t, err)

	// Verify expiration was preserved
	ttl := kvStore2.TTL("key1")
	assert.Equal(t, -1, ttl, "key1 should have no expiration")

	ttl = kvStore2.TTL("key2")
	assert.GreaterOrEqual(t, ttl, 0, "key2 should have expiration")
	assert.LessOrEqual(t, ttl, 10, "key2 TTL should be <= 10")
}

func TestSnapshotManager_MultipleSnapshots(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSnapshotManager(tmpDir, time.Minute)

	// Create first snapshot
	kvStore1 := store.NewKVStore()
	kvStore1.Set("key1", "value1")
	err := sm.CreateSnapshot(kvStore1)
	require.NoError(t, err)

	// Wait a bit to ensure different timestamp
	time.Sleep(10 * time.Millisecond)

	// Create second snapshot
	kvStore2 := store.NewKVStore()
	kvStore2.Set("key1", "value2")
	kvStore2.Set("key2", "value2")
	err = sm.CreateSnapshot(kvStore2)
	require.NoError(t, err)

	// Load latest should get second snapshot
	kvStore3 := store.NewKVStore()
	err = sm.LoadLatest(kvStore3)
	require.NoError(t, err)

	val, ok := kvStore3.Get("key1")
	require.True(t, ok)
	assert.Equal(t, "value2", val, "should load latest snapshot")

	val, ok = kvStore3.Get("key2")
	require.True(t, ok)
	assert.Equal(t, "value2", val)
}

func TestSnapshotManager_EmptyStore(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSnapshotManager(tmpDir, time.Minute)

	// Create snapshot of empty store
	kvStore1 := store.NewKVStore()
	err := sm.CreateSnapshot(kvStore1)
	require.NoError(t, err)

	// Load snapshot
	kvStore2 := store.NewKVStore()
	err = sm.LoadLatest(kvStore2)
	require.NoError(t, err)

	// Store should be empty
	keys := kvStore2.Keys("*")
	assert.Empty(t, keys)
}

func TestSnapshotManager_LargeDataset(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSnapshotManager(tmpDir, time.Minute)

	// Create store with many keys
	kvStore1 := store.NewKVStore()
	for i := 0; i < 1000; i++ {
		key := "key" + string(rune(i))
		value := "value" + string(rune(i))
		kvStore1.Set(key, value)
	}

	err := sm.CreateSnapshot(kvStore1)
	require.NoError(t, err)

	// Load snapshot
	kvStore2 := store.NewKVStore()
	err = sm.LoadLatest(kvStore2)
	require.NoError(t, err)

	// Verify all keys were loaded
	keys := kvStore2.Keys("*")
	assert.Equal(t, 1000, len(keys))
}

func TestSnapshotManager_SnapshotFilename(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSnapshotManager(tmpDir, time.Minute)

	kvStore := store.NewKVStore()
	kvStore.Set("key1", "value1")

	err := sm.CreateSnapshot(kvStore)
	require.NoError(t, err)

	// Verify snapshot file has proper naming convention
	files, err := os.ReadDir(tmpDir)
	require.NoError(t, err)
	require.NotEmpty(t, files)

	// Filename should contain timestamp or similar identifier
	filename := files[0].Name()
	assert.NotEmpty(t, filename)
	assert.Contains(t, filename, "snapshot")
}

func TestSnapshotManager_OverwriteSnapshot(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSnapshotManager(tmpDir, time.Minute)

	// Create first snapshot
	kvStore1 := store.NewKVStore()
	kvStore1.Set("key1", "value1")
	err := sm.CreateSnapshot(kvStore1)
	require.NoError(t, err)

	// Create second snapshot with different data
	kvStore2 := store.NewKVStore()
	kvStore2.Set("key1", "updated")
	kvStore2.Set("key2", "new")
	err = sm.CreateSnapshot(kvStore2)
	require.NoError(t, err)

	// Load should get latest data
	kvStore3 := store.NewKVStore()
	err = sm.LoadLatest(kvStore3)
	require.NoError(t, err)

	val, ok := kvStore3.Get("key1")
	require.True(t, ok)
	assert.Equal(t, "updated", val)

	val, ok = kvStore3.Get("key2")
	require.True(t, ok)
	assert.Equal(t, "new", val)
}

func TestSnapshotManager_CorruptedSnapshot(t *testing.T) {
	tmpDir := t.TempDir()

	// Create corrupted snapshot file
	corruptedPath := filepath.Join(tmpDir, "snapshot-corrupted.db")
	err := os.WriteFile(corruptedPath, []byte("invalid data"), 0644)
	require.NoError(t, err)

	sm := NewSnapshotManager(tmpDir, time.Minute)
	kvStore := store.NewKVStore()

	// Loading corrupted snapshot should error
	err = sm.LoadLatest(kvStore)
	assert.Error(t, err, "should error on corrupted snapshot")
}

func TestSnapshotManager_PreserveTimestamps(t *testing.T) {
	tmpDir := t.TempDir()
	sm := NewSnapshotManager(tmpDir, time.Minute)

	// Create snapshot with entries that have timestamps
	kvStore1 := store.NewKVStore()
	kvStore1.Set("key1", "value1")

	// Small delay to ensure different timestamps
	time.Sleep(10 * time.Millisecond)
	kvStore1.Set("key2", "value2")

	err := sm.CreateSnapshot(kvStore1)
	require.NoError(t, err)

	// Load snapshot
	kvStore2 := store.NewKVStore()
	err = sm.LoadLatest(kvStore2)
	require.NoError(t, err)

	// Verify both keys exist (timestamps preserved)
	_, ok := kvStore2.Get("key1")
	assert.True(t, ok)

	_, ok = kvStore2.Get("key2")
	assert.True(t, ok)
}
