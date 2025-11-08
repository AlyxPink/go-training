package persistence

import (
	"os"
	"sync"

	"github.com/alyxpink/go-training/kvstore/store"
)

type WAL struct {
	file *os.File
	mu   sync.Mutex
}

func NewWAL(path string) (*WAL, error) {
	// TODO: Open WAL file for append
	panic("not implemented")
}

func (w *WAL) Append(command string) error {
	// TODO: Write command to WAL
	panic("not implemented")
}

func (w *WAL) Replay(store *store.KVStore) error {
	// TODO: Read WAL and replay commands
	panic("not implemented")
}

func (w *WAL) Close() error {
	// TODO: Close WAL file
	panic("not implemented")
}
