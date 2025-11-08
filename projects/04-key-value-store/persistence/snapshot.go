package persistence

import (
	"encoding/gob"
	"os"
	"path/filepath"
	"time"

	"github.com/alyxpink/go-training/kvstore/store"
)

type SnapshotManager struct {
	dataDir  string
	interval time.Duration
}

func NewSnapshotManager(dataDir string, interval time.Duration) *SnapshotManager {
	return &SnapshotManager{
		dataDir:  dataDir,
		interval: interval,
	}
}

func (sm *SnapshotManager) Run(store *store.KVStore) {
	ticker := time.NewTicker(sm.interval)
	defer ticker.Stop()

	for range ticker.C {
		sm.CreateSnapshot(store)
	}
}

func (sm *SnapshotManager) CreateSnapshot(store *store.KVStore) error {
	// TODO: Serialize store to disk
	path := filepath.Join(sm.dataDir, "snapshot.gob")
	tempPath := path + ".tmp"

	file, err := os.Create(tempPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(store); err != nil {
		return err
	}

	return os.Rename(tempPath, path)
}

func (sm *SnapshotManager) LoadLatest(store *store.KVStore) error {
	// TODO: Load snapshot from disk
	path := filepath.Join(sm.dataDir, "snapshot.gob")

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	return decoder.Decode(store)
}
