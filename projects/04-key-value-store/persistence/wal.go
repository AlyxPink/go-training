package persistence

import (
	"bufio"
	"os"
	"strings"
	"sync"

	"github.com/alyxpink/go-training/kvstore/store"
)

type WAL struct {
	file *os.File
	mu   sync.Mutex
}

func NewWAL(path string) (*WAL, error) {
	// TODO: Open WAL file for append
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &WAL{file: file}, nil
}

func (w *WAL) Append(command string) error {
	// TODO: Write command to WAL
	w.mu.Lock()
	defer w.mu.Unlock()

	_, err := w.file.WriteString(command + "\n")
	if err != nil {
		return err
	}

	return w.file.Sync()
}

func (w *WAL) Replay(store *store.KVStore) error {
	// TODO: Read WAL and replay commands
	file, err := os.Open(w.file.Name())
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		cmd := strings.ToUpper(parts[0])
		args := parts[1:]

		switch cmd {
		case "SET":
			if len(args) >= 2 {
				store.Set(args[0], strings.Join(args[1:], " "))
			}
		case "DEL":
			if len(args) >= 1 {
				store.Del(args[0])
			}
		}
	}

	return scanner.Err()
}

func (w *WAL) Close() error {
	return w.file.Close()
}
