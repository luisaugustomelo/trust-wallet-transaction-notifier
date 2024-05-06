package mocks

import (
	"net/http"

	"github.com/luisaugustmelo/trust-wallet-transaction-notifier/entities"
	"github.com/stretchr/testify/mock"
)

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) StartBlockWatcher() {
	m.Called()
}

func (m *MockHTTPClient) GetCurrentBlock() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockHTTPClient) CleanUpTransactions(address string) {
	m.Called(address)
}

func (m *MockHTTPClient) Subscribe(address string) bool {
	args := m.Called(address)
	return args.Bool(0)
}

func (m *MockHTTPClient) GetTransactions(address string) ([]entities.Transaction, error) {
	args := m.Called(address)
	return args.Get(0).([]entities.Transaction), args.Error(1)
}

func (m *MockHTTPClient) GetTransactionsFromBlock(blockNumber int64, address string) ([]entities.Transaction, error) {
	args := m.Called(blockNumber, address)
	return args.Get(0).([]entities.Transaction), args.Error(1)
}

func (m *MockHTTPClient) MakeRPCRequest(data string) (*http.Response, error) {
	args := m.Called(data)
	return args.Get(0).(*http.Response), args.Error(1)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}
