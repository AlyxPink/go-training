package persistence

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alyxpink/go-training/kvstore/store"
)

type SnapshotManager struct {
	dataDir  string
	interval time.Duration
}

type SnapshotData struct {
	Data map[string]*store.Entry
}

func NewSnapshotManager(dataDir string, interval time.Duration) *SnapshotManager {
	return &SnapshotManager{
		dataDir:  dataDir,
		interval: interval,
	}
}

func (sm *SnapshotManager) Run(kvStore *store.KVStore) {
	ticker := time.NewTicker(sm.interval)
	defer ticker.Stop()

	for range ticker.C {
		sm.CreateSnapshot(kvStore)
	}
}

func (sm *SnapshotManager) CreateSnapshot(kvStore *store.KVStore) error {
	// Create snapshot filename with timestamp
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("snapshot-%d.db", timestamp)
	tempPath := filepath.Join(sm.dataDir, filename+".tmp")
	finalPath := filepath.Join(sm.dataDir, filename)

	// Create temporary file
	file, err := os.Create(tempPath)
	if err != nil {
		return err
	}

	// Get all keys and their entries
	keys := kvStore.Keys("*")
	snapshotData := SnapshotData{
		Data: make(map[string]*store.Entry),
	}

	// Copy data from store
	for _, key := range keys {
		if value, ok := kvStore.Get(key); ok {
			// We need to reconstruct the entry
			entry := &store.Entry{
				Value:     value,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			// Check if key has TTL
			ttl := kvStore.TTL(key)
			if ttl > 0 {
				expiresAt := time.Now().Add(time.Duration(ttl) * time.Second)
				entry.ExpiresAt = &expiresAt
			}

			snapshotData.Data[key] = entry
		}
	}

	// Encode to gob
	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(snapshotData); err != nil {
		file.Close()
		os.Remove(tempPath)
		return err
	}

	file.Close()

	// Atomic rename
	if err := os.Rename(tempPath, finalPath); err != nil {
		return err
	}

	// Clean up old snapshots (keep only the latest 3)
	sm.cleanupOldSnapshots()

	return nil
}

func (sm *SnapshotManager) LoadLatest(kvStore *store.KVStore) error {
	// Find latest snapshot file
	files, err := os.ReadDir(sm.dataDir)
	if err != nil {
		return err
	}

	var snapshots []string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "snapshot-") && strings.HasSuffix(file.Name(), ".db") {
			snapshots = append(snapshots, file.Name())
		}
	}

	if len(snapshots) == 0 {
		return fmt.Errorf("no snapshot found")
	}

	// Sort to get latest
	sort.Strings(snapshots)
	latestSnapshot := snapshots[len(snapshots)-1]

	// Load snapshot
	snapshotPath := filepath.Join(sm.dataDir, latestSnapshot)
	file, err := os.Open(snapshotPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var snapshotData SnapshotData
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&snapshotData); err != nil {
		return err
	}

	// Load data into store
	for key, entry := range snapshotData.Data {
		kvStore.Set(key, entry.Value)

		// Restore expiration if it exists
		if entry.ExpiresAt != nil {
			ttl := time.Until(*entry.ExpiresAt).Seconds()
			if ttl > 0 {
				kvStore.Expire(key, int(ttl))
			}
		}
	}

	return nil
}

func (sm *SnapshotManager) cleanupOldSnapshots() {
	files, err := os.ReadDir(sm.dataDir)
	if err != nil {
		return
	}

	var snapshots []string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "snapshot-") && strings.HasSuffix(file.Name(), ".db") {
			snapshots = append(snapshots, file.Name())
		}
	}

	if len(snapshots) <= 3 {
		return
	}

	// Sort and remove old snapshots
	sort.Strings(snapshots)
	for i := 0; i < len(snapshots)-3; i++ {
		os.Remove(filepath.Join(sm.dataDir, snapshots[i]))
	}
}
