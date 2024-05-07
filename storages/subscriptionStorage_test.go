package storages

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSubscriptionStorageSave tests the Save method of SubscriptionStorage.
func TestSubscriptionStorageSave(t *testing.T) {
	storage := &SubscriptionStorage{
		subscriptions: make(map[string]int64),
	}

	// Test saving a new key-value pair
	storage.Save("user1", int64(1234567890))
	val, exists := storage.subscriptions["user1"]
	require.True(t, exists, "The key should exist after saving.")
	assert.Equal(t, int64(1234567890), val, "The value should match the saved value.")

	// Test saving with a non-int64 value (should not save)
	storage.Save("user2", "not-an-int64")
	_, exists = storage.subscriptions["user2"]
	assert.False(t, exists, "The key should not exist as the value was not an int64.")
}

// TestSubscriptionStorageDelete tests the Delete method.
func TestSubscriptionStorageDelete(t *testing.T) {
	storage := &SubscriptionStorage{
		subscriptions: map[string]int64{"user1": 1234567890},
	}

	// Test deleting an existing key
	storage.Delete("user1")
	_, exists := storage.subscriptions["user1"]
	assert.False(t, exists, "The key should no longer exist.")

	// Test deleting a non-existing key
	storage.Delete("user2")
	_, exists = storage.subscriptions["user2"]
	assert.False(t, exists, "The key should not exist as it was never added.")
}

// TestSubscriptionStorageFind tests the Find method.
func TestSubscriptionStorageFind(t *testing.T) {
	storage := &SubscriptionStorage{
		subscriptions: map[string]int64{"user1": 1234567890},
	}

	// Test finding an existing key
	value, exists := storage.Find("user1")
	assert.True(t, exists, "The key should exist.")
	assert.Equal(t, int64(1234567890), value.(int64), "The value should match.")

	// Test finding a non-existing key
	_, exists = storage.Find("nonexistent")
	assert.False(t, exists, "The key should not exist.")
}

// TestSubscriptionStorageUpdate tests the Update method.
func TestSubscriptionStorageUpdate(t *testing.T) {
	storage := &SubscriptionStorage{
		subscriptions: map[string]int64{"user1": 1234567890},
	}

	// Test updating an existing key
	storage.Update("user1", int64(987654321))
	val, exists := storage.subscriptions["user1"]
	assert.True(t, exists, "The key should still exist after update.")
	assert.Equal(t, int64(987654321), val, "The value should be updated.")

	// Test updating a non-existing key
	storage.Update("user2", int64(555555555))
	_, exists = storage.subscriptions["user2"]
	assert.False(t, exists, "Update should not create a new key.")
}

// TestSubscriptionStorageGetAll tests the GetAll method.
func TestSubscriptionStorageGetAll(t *testing.T) {
	storage := &SubscriptionStorage{
		subscriptions: map[string]int64{"user1": 1234567890, "user2": 987654321},
	}

	allSubs := storage.GetAll().(map[string]int64)
	assert.Equal(t, 2, len(allSubs), "There should be two entries in the map.")
	assert.Equal(t, int64(1234567890), allSubs["user1"], "The value for 'user1' should match.")
	assert.Equal(t, int64(987654321), allSubs["user2"], "The value for 'user2' should match.")
}
