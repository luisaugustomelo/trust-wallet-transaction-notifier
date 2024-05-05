package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/luisaugustmelo/trust-wallet-transaction-notifier/interfaces"
)

func HandleCurrentBlock(w http.ResponseWriter, r *http.Request, rpc interfaces.Parser) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	block := rpc.GetCurrentBlock()
	fmt.Fprintf(w, "Current Block: %d", block)
}

func HandleSubscribe(w http.ResponseWriter, r *http.Request, rpc interfaces.Parser) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Address string `json:"address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if data.Address == "" {
		http.Error(w, "Address is required", http.StatusBadRequest)
		return
	}

	if rpc.Subscribe(data.Address) {
		fmt.Fprintf(w, "Subscribed to: %s", data.Address)
	} else {
		fmt.Fprintf(w, "Already subscribed to: %s", data.Address)
	}
}

func HandleTransactions(w http.ResponseWriter, r *http.Request, rpc interfaces.Parser) {
	address := r.URL.Query().Get("address")
	transactions, err := rpc.GetTransactions(address)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting transactions: %v", err), http.StatusBadRequest)
		return
	}

	if len(transactions) == 0 {
		http.Error(w, fmt.Sprintf("\"There are no new transactions at this time, please wait\""), http.StatusNoContent)
		return
	}

	js, err := json.Marshal(transactions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
