package protocol

import (
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
	panic("not implemented")
}

func (h *Handler) handleSet(args []string) string {
	// TODO: Implement SET command
	panic("not implemented")
}

func (h *Handler) handleGet(args []string) string {
	// TODO: Implement GET command
	panic("not implemented")
}

func (h *Handler) handleDel(args []string) string {
	// TODO: Implement DEL command
	panic("not implemented")
}

func (h *Handler) handleExists(args []string) string {
	// TODO: Implement EXISTS command
	panic("not implemented")
}

func (h *Handler) handleKeys(args []string) string {
	// TODO: Implement KEYS command
	panic("not implemented")
}

func (h *Handler) handleExpire(args []string) string {
	// TODO: Implement EXPIRE command
	panic("not implemented")
}

func (h *Handler) handleTTL(args []string) string {
	// TODO: Implement TTL command
	panic("not implemented")
}
