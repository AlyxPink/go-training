package protocol

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alyxpink/go-training/kvstore/persistence"
	"github.com/alyxpink/go-training/kvstore/store"
)

type Handler struct {
	store *store.KVStore
	wal   *persistence.WAL
}

func NewHandler(store *store.KVStore, wal *persistence.WAL) *Handler {
	return &Handler{store: store, wal: wal}
}

func (h *Handler) Handle(line string) string {
	// Parse command
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return "-ERR empty command"
	}

	command := strings.ToUpper(parts[0])
	args := parts[1:]

	// Route to appropriate handler
	switch command {
	case "SET":
		return h.handleSet(args)
	case "GET":
		return h.handleGet(args)
	case "DEL":
		return h.handleDel(args)
	case "EXISTS":
		return h.handleExists(args)
	case "KEYS":
		return h.handleKeys(args)
	case "EXPIRE":
		return h.handleExpire(args)
	case "TTL":
		return h.handleTTL(args)
	default:
		return "-ERR unknown command '" + command + "'"
	}
}

func (h *Handler) handleSet(args []string) string {
	if len(args) < 2 {
		return "-ERR wrong number of arguments for 'set' command"
	}

	key := args[0]
	value := strings.Join(args[1:], " ")

	h.store.Set(key, value)

	// Log to WAL
	h.wal.Append(fmt.Sprintf("SET %s %s", key, value))

	return "+OK"
}

func (h *Handler) handleGet(args []string) string {
	if len(args) < 1 {
		return "-ERR wrong number of arguments for 'get' command"
	}

	key := args[0]
	value, ok := h.store.Get(key)
	if !ok {
		return "$-1"
	}

	// Redis bulk string format: $<length>\r\n<data>
	return fmt.Sprintf("$%d\r\n%s", len(value), value)
}

func (h *Handler) handleDel(args []string) string {
	if len(args) < 1 {
		return "-ERR wrong number of arguments for 'del' command"
	}

	key := args[0]
	deleted := h.store.Del(key)

	// Log to WAL
	if deleted {
		h.wal.Append(fmt.Sprintf("DEL %s", key))
		return ":1"
	}

	return ":0"
}

func (h *Handler) handleExists(args []string) string {
	if len(args) < 1 {
		return "-ERR wrong number of arguments for 'exists' command"
	}

	key := args[0]
	exists := h.store.Exists(key)

	if exists {
		return ":1"
	}
	return ":0"
}

func (h *Handler) handleKeys(args []string) string {
	if len(args) < 1 {
		return "-ERR wrong number of arguments for 'keys' command"
	}

	pattern := args[0]
	keys := h.store.Keys(pattern)

	// Redis array format: *<count>\r\n$<len>\r\n<key>\r\n...
	result := fmt.Sprintf("*%d", len(keys))
	for _, key := range keys {
		result += fmt.Sprintf("\r\n$%d\r\n%s", len(key), key)
	}

	return result
}

func (h *Handler) handleExpire(args []string) string {
	if len(args) < 2 {
		return "-ERR wrong number of arguments for 'expire' command"
	}

	key := args[0]
	seconds, err := strconv.Atoi(args[1])
	if err != nil {
		return "-ERR value is not an integer or out of range"
	}

	success := h.store.Expire(key, seconds)

	// Log to WAL
	if success {
		h.wal.Append(fmt.Sprintf("EXPIRE %s %d", key, seconds))
		return ":1"
	}

	return ":0"
}

func (h *Handler) handleTTL(args []string) string {
	if len(args) < 1 {
		return "-ERR wrong number of arguments for 'ttl' command"
	}

	key := args[0]
	ttl := h.store.TTL(key)

	return fmt.Sprintf(":%d", ttl)
}
