package interfaces

import (
	"net/http"

	"github.com/luisaugustmelo/trust-wallet-transaction-notifier/entities"
)

// Storage defines the CRUD operations.
type Storage interface {
	Save(key string, value interface{})
	Delete(key string)
	Find(key string) (interface{}, bool)
	Update(key string, value interface{})
	GetAll() interface{}
}

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) ([]entities.Transaction, error)
	GetTransactionsFromBlock(blockNumber int64, address string) ([]entities.Transaction, error)
	MakeRPCRequest(data string) (*http.Response, error)
	StartBlockWatcher()
	CleanUpTransactions(address string)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}
