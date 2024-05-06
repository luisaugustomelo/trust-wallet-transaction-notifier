package storages

import (
	"sync"

	"github.com/luisaugustmelo/trust-wallet-transaction-notifier/interfaces"
)

// SubscriptionStorage manages subscription data.
type SubscriptionStorage struct {
	subscriptions map[string]int64
	mu            sync.RWMutex
}

// Ensures that SubscriptionStorage implements Storage
var _ interfaces.Storage = (*SubscriptionStorage)(nil)

func NewSubscriptionStorage() *SubscriptionStorage {
	return &SubscriptionStorage{
		subscriptions: make(map[string]int64),
	}
}

// Save SubscriptionStorage methods
func (s *SubscriptionStorage) Save(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if val, ok := value.(int64); ok {
		s.subscriptions[key] = val
	}
}

func (s *SubscriptionStorage) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.subscriptions[key]; exists {
		delete(s.subscriptions, key)
	}
}

func (s *SubscriptionStorage) Find(key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if value, exists := s.subscriptions[key]; exists {
		return value, true
	}
	return nil, false
}

func (s *SubscriptionStorage) Update(key string, value interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if val, ok := value.(int64); ok {
		if _, exists := s.subscriptions[key]; exists {
			s.subscriptions[key] = val
		}
	}
}

func (s *SubscriptionStorage) GetAll() interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Copy to prevent external modifications
	c := make(map[string]int64)
	for k, v := range s.subscriptions {
		c[k] = v
	}
	return c
}
