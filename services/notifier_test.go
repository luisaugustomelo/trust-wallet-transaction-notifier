package services

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/luisaugustmelo/trust-wallet-transaction-notifier/entities"
	"github.com/luisaugustmelo/trust-wallet-transaction-notifier/services/mocks"
	"github.com/luisaugustmelo/trust-wallet-transaction-notifier/storages"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetCurrentBlock(t *testing.T) {
	mockClient := new(mocks.MockHTTPClient)
	service := EthereumRPC{Methods: mockClient}

	responseBody := `{"jsonrpc":"2.0","id":1,"result":"0x5ba"}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(responseBody)))

	mockClient.On("MakeRPCRequest", mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil)

	blockNumber := service.GetCurrentBlock()
	assert.Equal(t, 1466, blockNumber) // 0x5ba in decimal
	mockClient.AssertExpectations(t)
}

func TestSubscribe(t *testing.T) {
	mockClient := new(mocks.MockHTTPClient)
	mockSubStorage := new(mocks.MockSubscriptionStorage)
	mockTransStorage := new(mocks.MockTransactionStorage)

	// Configuring mock for MemoryStorage
	mockSubStorage.On("Find", "0x123").Return(nil, false)          // First call returns not found
	mockSubStorage.On("Save", "0x123", int64(100000)).Return(nil)  // Simulates successful save
	mockSubStorage.On("Find", "0x123").Return(int64(100000), true) // Second call finds the subscription

	mockStorage := storages.NewMemoryStorage(mockSubStorage, mockTransStorage)

	service := EthereumRPC{
		URL:     "http://example.com",
		Client:  mockClient,
		Storage: mockStorage,
		Methods: mockClient,
	}

	// Configuring the mocked response for GetCurrentBlock
	mockClient.On("GetCurrentBlock").Return(100000, nil)

	// First attempt at subscription
	success := service.Subscribe("0x123")
	assert.True(t, success, "Subscription should succeed on first attempt")

	// Subsequent attempt to subscribe with the same address
	alreadySubscribed := service.Subscribe("0x123")
	assert.True(t, alreadySubscribed, "Subscription should fail on second attempt with the same address")
}

// Test for GetTransactionsFromBlock
func TestGetTransactionsFromBlock(t *testing.T) {
	mockClient := new(mocks.MockHTTPClient)
	service := EthereumRPC{Methods: mockClient}

	// Simulating RPC response
	mockResponse := `{"jsonrpc":"2.0","result":{"transactions":[{"from":"0x123","to":"0x456","value":"100"},{"from":"0x789","to":"0x123","value":"200"}]}}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(mockResponse)))

	mockClient.On("MakeRPCRequest", mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body:       r,
	}, nil)

	transactions, err := service.GetTransactionsFromBlock(123456, "0x123")
	assert.NoError(t, err)
	assert.Len(t, transactions, 2) // Only one transaction involves the address "0x123"
	assert.Equal(t, "0x123", transactions[0].From)
	assert.Equal(t, "0x123", transactions[1].To)
}

func TestGetTransactions(t *testing.T) {
	mockClient := new(mocks.MockHTTPClient)
	mockSubStorage := new(mocks.MockSubscriptionStorage)  // Mock for subscriptions
	mockTransStorage := new(mocks.MockTransactionStorage) // Mock for transactions

	mockStorage := storages.NewMemoryStorage(mockSubStorage, mockTransStorage)

	// Configuring mocks for transaction storage
	transactions := []entities.Transaction{
		{From: "0x789", To: "0x123", Value: "111", Hash: "xxx"},
		{From: "0x123", To: "0x789", Value: "3", Hash: "yyy"},
		{From: "0x123", To: "0x456", Value: "100", Hash: "zzz"},
	}

	mockTransStorage.On("Find", "0x123").Return(transactions, true)
	mockTransStorage.On("Find", "0x999").Return(nil, false)

	service := EthereumRPC{
		URL:     "http://example.com",
		Client:  mockClient,
		Storage: mockStorage,
		Methods: mockClient,
	}

	// Test to retrieve transactions for the address "0x123"
	transactionsResult, err := service.GetTransactions("0x123")
	assert.NoError(t, err)
	assert.Len(t, transactionsResult, 3, "Expected one transaction from each block, totaling 3 transactions")

	// Test to try to retrieve transactions for a non-existent address "0x999"
	_, err = service.GetTransactions("0x999")
	assert.Error(t, err, "should return an error for an unsubscribed address")
}
