package persistence

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/alyxpink/go-training/kvstore/store"
)

type WAL struct {
	file   *os.File
	mu     sync.Mutex
	closed bool
}

func NewWAL(path string) (*WAL, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	return &WAL{
		file:   file,
		closed: false,
	}, nil
}

func (w *WAL) Append(command string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed {
		return fmt.Errorf("WAL is closed")
	}

	// Write command with newline
	if _, err := w.file.WriteString(command + "\n"); err != nil {
		return err
	}

	// Sync to disk for durability
	return w.file.Sync()
}

func (w *WAL) Replay(kvStore *store.KVStore) error {
	// Open the file for reading from the beginning
	file, err := os.Open(w.file.Name())
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		// Parse and execute command
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}

		command := strings.ToUpper(parts[0])
		args := parts[1:]

		switch command {
		case "SET":
			if len(args) >= 2 {
				// Join all parts after the key as the value
				value := strings.Join(args[1:], " ")
				kvStore.Set(args[0], value)
			}
		case "DEL":
			if len(args) >= 1 {
				kvStore.Del(args[0])
			}
		case "EXPIRE":
			if len(args) >= 2 {
				seconds, err := strconv.Atoi(args[1])
				if err == nil {
					kvStore.Expire(args[0], seconds)
				}
			}
		}
	}

	return scanner.Err()
}

func (w *WAL) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.closed {
		return nil
	}

	w.closed = true
	return w.file.Close()
}
