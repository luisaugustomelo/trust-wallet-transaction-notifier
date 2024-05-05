package storages

import (
	"fmt"
	"sync"

	"github.com/luisaugustmelo/trust-wallet-transaction-notifier/entities"
	"github.com/luisaugustmelo/trust-wallet-transaction-notifier/interfaces"
)

// TransactionStorage manages transactions data.
type TransactionStorage struct {
	transactions map[string][]entities.Transaction
	mu           sync.RWMutex
}

// Ensures that TransactionStorage implements Storage
var _ interfaces.Storage = (*TransactionStorage)(nil)

// TransactionStorage methods

func (t *TransactionStorage) Save(key string, value interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Check and initialize the map if it has not been initialized yet
	if t.transactions == nil {
		t.transactions = make(map[string][]entities.Transaction)
	}

	// Saves the value to the map if it is of the correct type ([]entities.Transaction)
	if val, ok := value.([]entities.Transaction); ok {
		t.transactions[key] = val
	}
}

func (t *TransactionStorage) Delete(key string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, exists := t.transactions[key]; exists {
		delete(t.transactions, key)
		return nil
	}
	return fmt.Errorf("no transactions found for key %s", key)
}

func (t *TransactionStorage) Find(key string) (interface{}, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if value, exists := t.transactions[key]; exists {
		return value, true
	}
	return nil, false
}

func (t *TransactionStorage) Update(key string, value interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if val, ok := value.([]entities.Transaction); ok {
		if _, exists := t.transactions[key]; exists {
			t.transactions[key] = val
		}
	}
}

func (t *TransactionStorage) GetAll() interface{} {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// Copy to prevent external modifications
	c := make(map[string][]entities.Transaction)
	for k, v := range t.transactions {
		c[k] = v
	}
	return c
}
