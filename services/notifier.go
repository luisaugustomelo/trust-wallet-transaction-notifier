package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/luisaugustmelo/trust-wallet-transaction-notifier/entities"
	"github.com/luisaugustmelo/trust-wallet-transaction-notifier/interfaces"
	"github.com/luisaugustmelo/trust-wallet-transaction-notifier/storages"
)

type EthereumRPC struct {
	URL     string
	Storage *storages.MemoryStorage
	mu      sync.Mutex
	Client  interfaces.HTTPClient
	Methods interfaces.Parser
}

func NewEthereumRPC(url string, client interfaces.HTTPClient, storage *storages.MemoryStorage) interfaces.Parser {
	rpc := &EthereumRPC{
		URL:     url,
		Storage: storage,
		Client:  client,
	}

	var _ interfaces.Parser = rpc

	rpc.Methods = rpc
	go rpc.StartBlockWatcher()

	return rpc
}

func (rpc *EthereumRPC) StartBlockWatcher() {
	ticker := time.NewTicker(1 * time.Second) // Checks for new blocks every 1 second
	for {
		select {
		case <-ticker.C:
			currentBlock := rpc.Methods.GetCurrentBlock()
			for address, lastCheckedBlock := range rpc.Storage.Subscriptions.GetAll().(map[string]int64) {
				for block := lastCheckedBlock + 1; block <= int64(currentBlock); block++ {
					transactions, err := rpc.GetTransactionsFromBlock(block, address)
					if err != nil {
						fmt.Printf("Error fetching transactions for block %d and address %s: %v", block, address, err)
						continue
					}

					if transactions != nil {
						// Updates transactions and signatures using storage-specific methods
						rpc.Storage.Transactions.Save(address, transactions)
						rpc.Storage.Subscriptions.Update(address, block)
					}

				}
			}
		}
	}
}

func (rpc *EthereumRPC) GetCurrentBlock() int {
	// Json request data
	requestData := `{"jsonrpc":"2.0", "method":"eth_blockNumber", "params":[], "id":1}`

	// Make RPC call using makeRPCRequest method
	resp, err := rpc.Methods.MakeRPCRequest(requestData)
	if err != nil {
		fmt.Println("Error making RPC request:", err)
		return -1
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error body read closer:", err)
		}
	}(resp.Body) // close body to make resources free

	// Decode the JSON response
	var result struct {
		JSONRPC string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Result  string `json:"result"`
		Error   struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return -1
	}

	// Verify JSONRPC response error
	if result.Error.Code != 0 {
		fmt.Printf("Error in RPC response: %s (Code: %d)\n", result.Error.Message, result.Error.Code)
		return -1
	}

	// Hexadecimal to decimal
	blockNumber, err := strconv.ParseInt(result.Result, 0, 64)
	if err != nil {
		fmt.Println("Error parsing block number:", err)
		return -1
	}

	return int(blockNumber)
}

func (rpc *EthereumRPC) CleanUpTransactions(address string) {
	rpc.Storage.Transactions.Delete(address)
}

func (rpc *EthereumRPC) Subscribe(address string) bool {
	rpc.mu.Lock()
	defer rpc.mu.Unlock()
	startBlock := rpc.Methods.GetCurrentBlock()
	if _, exists := rpc.Storage.Subscriptions.Find(address); !exists {
		rpc.Storage.Subscriptions.Save(address, int64(startBlock))
		return true
	}
	return false
}

func (rpc *EthereumRPC) GetTransactionsFromBlock(blockNumber int64, address string) ([]entities.Transaction, error) {
	hexBlockNumber := fmt.Sprintf("0x%x", blockNumber)
	requestData := fmt.Sprintf(`{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["%s", true],"id":1}`, hexBlockNumber)

	resp, err := rpc.Methods.MakeRPCRequest(requestData)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error body read closer:", err)
		}
	}(resp.Body)

	var rpcResult struct {
		Result struct {
			Transactions []entities.Transaction `json:"transactions"`
		} `json:"result"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&rpcResult); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	// Filter transactions to only include those involving the specified address
	var filteredTransactions []entities.Transaction
	for _, tx := range rpcResult.Result.Transactions {
		if strings.ToLower(tx.From) == strings.ToLower(address) || strings.ToLower(tx.To) == strings.ToLower(address) {
			filteredTransactions = append(filteredTransactions, tx)
		}
	}

	return filteredTransactions, nil
}

func (rpc *EthereumRPC) GetTransactions(address string) ([]entities.Transaction, error) {
	rpc.mu.Lock()
	defer rpc.mu.Unlock()

	transactions, exists := rpc.Storage.Transactions.Find(address)
	if exists {
		return transactions.([]entities.Transaction), nil
	}

	return nil, fmt.Errorf("no transactions found for address %s", address)
}

func (rpc *EthereumRPC) MakeRPCRequest(data string) (*http.Response, error) {
	req, err := http.NewRequest("POST", rpc.URL, bytes.NewBufferString(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := rpc.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
