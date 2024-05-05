package storages

import (
	"github.com/luisaugustmelo/trust-wallet-transaction-notifier/interfaces"
)

// MemoryStorage aggregates different types of storages.
type MemoryStorage struct {
	Subscriptions interfaces.Storage
	Transactions  interfaces.Storage
}

// NewMemoryStorage creates a new MemoryStorage instance with initialized sub-storages.
func NewMemoryStorage(subs interfaces.Storage, trans interfaces.Storage) *MemoryStorage {
	return &MemoryStorage{
		Subscriptions: subs,
		Transactions:  trans,
	}
}
