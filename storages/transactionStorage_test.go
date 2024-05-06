package storages

import (
	"testing"

	"github.com/luisaugustmelo/trust-wallet-transaction-notifier/entities"
	"github.com/stretchr/testify/assert"
)

func TestTransactionStorageSave(t *testing.T) {
	storage := &TransactionStorage{
		transactions: make(map[string][]entities.Transaction),
	}

	// Test saving transactions
	transactions := []entities.Transaction{{From: "addr1", To: "addr2", Value: "100", Hash: "hash1"}}
	storage.Save("tx1", transactions)
	assert.Equal(t, transactions, storage.transactions["tx1"], "Transactions should match the saved transactions.")

	// Test saving with incorrect type
	storage.Save("tx2", "incorrect type")
	_, exists := storage.transactions["tx2"]
	assert.False(t, exists, "No transaction should be saved with incorrect type.")
}

func TestTransactionStorageDelete(t *testing.T) {
	storage := &TransactionStorage{
		transactions: map[string][]entities.Transaction{
			"tx1": {{From: "addr1", To: "addr2", Value: "100", Hash: "hash1"}},
		},
	}

	// Test deleting an existing transaction
	err := storage.Delete("tx1")
	assert.NoError(t, err, "Deleting an existing transaction should not produce an error.")
	_, exists := storage.transactions["tx1"]
	assert.False(t, exists, "The key should no longer exist.")

	// Test deleting a non-existing transaction
	err = storage.Delete("tx2")
	assert.Error(t, err, "Deleting a non-existing transaction should produce an error.")
}

func TestTransactionStorageFind(t *testing.T) {
	storage := &TransactionStorage{
		transactions: map[string][]entities.Transaction{
			"tx1": {{From: "addr1", To: "addr2", Value: "100", Hash: "hash1"}},
		},
	}

	// Test finding an existing transaction
	value, exists := storage.Find("tx1")
	assert.True(t, exists, "The key should exist.")
	assert.Equal(t, storage.transactions["tx1"], value, "The transactions should match.")

	// Test finding a non-existing transaction
	_, exists = storage.Find("tx2")
	assert.False(t, exists, "The key should not exist.")
}

func TestTransactionStorageUpdate(t *testing.T) {
	storage := &TransactionStorage{
		transactions: map[string][]entities.Transaction{
			"tx1": {{From: "addr1", To: "addr2", Value: "50", Hash: "hash1"}},
		},
	}

	// Test updating an existing transaction
	newTransactions := []entities.Transaction{{From: "addr1", To: "addr3", Value: "150", Hash: "hash2"}}
	storage.Update("tx1", newTransactions)
	assert.Equal(t, newTransactions, storage.transactions["tx1"], "Transactions should be updated.")

	// Test updating a non-existing transaction
	storage.Update("tx2", newTransactions)
	_, exists := storage.transactions["tx2"]
	assert.False(t, exists, "Update should not create a new key.")
}

func TestTransactionStorageGetAll(t *testing.T) {
	storage := &TransactionStorage{
		transactions: map[string][]entities.Transaction{
			"tx1": {{From: "addr1", To: "addr2", Value: "100", Hash: "hash1"}},
			"tx2": {{From: "addr2", To: "addr3", Value: "200", Hash: "hash2"}},
		},
	}

	allTransactions := storage.GetAll().(map[string][]entities.Transaction)
	assert.Equal(t, 2, len(allTransactions), "There should be two entries in the map.")
	assert.Equal(t, storage.transactions["tx1"], allTransactions["tx1"], "The transactions for 'tx1' should match.")
	assert.Equal(t, storage.transactions["tx2"], allTransactions["tx2"], "The transactions for 'tx2' should match.")
}
