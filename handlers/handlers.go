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
	_, err := fmt.Fprintf(w, "Current Block: %d", block)
	if err != nil {
		return
	}
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
		_, err := fmt.Fprintf(w, "Subscribed to: %s", data.Address)
		if err != nil {
			return
		}
	} else {
		_, err := fmt.Fprintf(w, "Already subscribed to: %s", data.Address)
		if err != nil {
			return
		}
	}
}

func HandleTransactions(w http.ResponseWriter, r *http.Request, rpc interfaces.Parser) {
	address := r.URL.Query().Get("address")
	transactions, err := rpc.GetTransactions(address)
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		// Use json.Marshal to create a JSON-formatted error response
		errorResponse, _ := json.Marshal(map[string]string{
			"error": "Error getting transactions: " + err.Error(),
		})
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write(errorResponse)
		if err != nil {
			return
		}
		return
	}

	if len(transactions) == 0 {
		// Respond with a predefined message when no transactions are found
		noTransactionsResponse, _ := json.Marshal(map[string]string{
			"message": "There are no new transactions at this time, please wait!",
		})
		w.WriteHeader(http.StatusNotFound) // StatusNotFound (404) might be more appropriate here
		_, err := w.Write(noTransactionsResponse)
		if err != nil {
			return
		}
		return
	}

	// If transactions are found, marshal them into JSON and send the response
	js, err := json.Marshal(transactions)
	if err != nil {
		errorResponse, _ := json.Marshal(map[string]string{
			"error": "Failed to serialize transactions: " + err.Error(),
		})
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write(errorResponse)
		if err != nil {
			return
		}
		return
	}

	_, err = w.Write(js)
	if err != nil {
		return
	}
}
