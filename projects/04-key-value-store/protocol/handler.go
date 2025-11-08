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
	// TODO: Parse command and execute
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return "-ERR empty command"
	}

	cmd := strings.ToUpper(parts[0])
	args := parts[1:]

	// Log to WAL before executing (for SET, DEL, EXPIRE)
	if cmd == "SET" || cmd == "DEL" || cmd == "EXPIRE" {
		h.wal.Append(line)
	}

	switch cmd {
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
		return "-ERR unknown command"
	}
}

func (h *Handler) handleSet(args []string) string {
	if len(args) < 2 {
		return "-ERR wrong number of arguments"
	}
	key, value := args[0], strings.Join(args[1:], " ")
	h.store.Set(key, value)
	return "+OK"
}

func (h *Handler) handleGet(args []string) string {
	if len(args) != 1 {
		return "-ERR wrong number of arguments"
	}
	value, exists := h.store.Get(args[0])
	if !exists {
		return "$-1"
	}
	return fmt.Sprintf("$%d\r\n%s", len(value), value)
}

func (h *Handler) handleDel(args []string) string {
	if len(args) != 1 {
		return "-ERR wrong number of arguments"
	}
	deleted := h.store.Del(args[0])
	if deleted {
		return ":1"
	}
	return ":0"
}

func (h *Handler) handleExists(args []string) string {
	if len(args) != 1 {
		return "-ERR wrong number of arguments"
	}
	if h.store.Exists(args[0]) {
		return ":1"
	}
	return ":0"
}

func (h *Handler) handleKeys(args []string) string {
	pattern := "*"
	if len(args) > 0 {
		pattern = args[0]
	}
	keys := h.store.Keys(pattern)
	return fmt.Sprintf("*%d\r\n%s", len(keys), strings.Join(keys, "\r\n"))
}

func (h *Handler) handleExpire(args []string) string {
	if len(args) != 2 {
		return "-ERR wrong number of arguments"
	}
	seconds, err := strconv.Atoi(args[1])
	if err != nil {
		return "-ERR invalid expire time"
	}
	if h.store.Expire(args[0], seconds) {
		return ":1"
	}
	return ":0"
}

func (h *Handler) handleTTL(args []string) string {
	if len(args) != 1 {
		return "-ERR wrong number of arguments"
	}
	ttl := h.store.TTL(args[0])
	return fmt.Sprintf(":%d", ttl)
}
