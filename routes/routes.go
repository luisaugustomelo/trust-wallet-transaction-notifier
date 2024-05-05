package routes

import (
	"net/http"

	"github.com/luisaugustmelo/trust-wallet-transaction-notifier/handlers"
	"github.com/luisaugustmelo/trust-wallet-transaction-notifier/interfaces"
)

func RegisterRoutes(router *http.ServeMux, rpc interfaces.Parser) {
	router.HandleFunc("/currentBlock", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleCurrentBlock(w, r, rpc)
	})

	router.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleSubscribe(w, r, rpc)
	})

	router.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleTransactions(w, r, rpc)
	})
}
