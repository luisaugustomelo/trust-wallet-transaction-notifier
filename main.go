package main

import (
	"fmt"
	"net/http"

	"github.com/luisaugustmelo/trust-wallet-transaction-notifier/routes"
	"github.com/luisaugustmelo/trust-wallet-transaction-notifier/services"
	"github.com/luisaugustmelo/trust-wallet-transaction-notifier/storages"
)

func main() {
	client := &http.Client{}
	storage := storages.NewMemoryStorage(&storages.SubscriptionStorage{}, &storages.TransactionStorage{})
	rpc := services.NewEthereumRPC("https://cloudflare-eth.com", client, storage)

	router := http.NewServeMux()
	routes.RegisterRoutes(router, rpc)

	fmt.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Printf("Error starting HTTP server: %v\n", err)
		return
	}
}
