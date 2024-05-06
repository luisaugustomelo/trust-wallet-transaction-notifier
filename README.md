# Transaction Monitoring Application

## Overview

This application is designed to monitor transactions on specific blockchain addresses using Ethereum's JSON-RPC interface. It checks for new transactions at subscribed addresses and updates the subscribers with the latest transaction data.

## Commands

### Clone the Repository

First, you'll need to clone the Trust Wallet Transaction Notifier repository from GitHub. Open your terminal and run the following commands:

```
git clone https://github.com/luisaugustomelo/trust-wallet-transaction-notifier.git
cd trust-wallet-transaction-notifier
```

This will download the repository and change your current directory to the repository's directory.

### Install Dependencies

If the project has any dependencies, you can install them using Go's package manager. Run the following command in the project directory:

```
go mod tidy
```

This command cleans up the dependencies, making sure all necessary packages are downloaded and installed.

### Run the Application

To start the application, execute the following command in the terminal within the project directory:
```
go run main.go
```

This command compiles and runs the main.go file, starting your server. By default, the server should be accessible through http://localhost:8080, unless specified otherwise in the code.

## Architecture

### Components

- **EthereumRPC**: Core service that interfaces with the Ethereum blockchain using JSON-RPC requests. It provides functionalities to subscribe to addresses, fetch current block numbers, and retrieve transactions from specific blocks.

- **MemoryStorage**: Implements the `Storage` interface for in-memory data management, allowing quick access and updates to subscription and transaction data. It's designed to be easily replaceable with database storage systems if persistence or distributed storage is needed.

- **Subscription and Transaction Management**: Uses a combination of in-memory storage mechanisms and lock-based concurrency controls to manage subscriptions and transactions effectively. Subscriptions are monitored, and transactions are stored per address basis.

### Functions

- **Subscribe**: Allows a user to subscribe to a specific address. This function registers the address in the system, and any transactions involving this address will be tracked and stored.

- **GetTransactions**: Retrieves the list of transactions for a subscribed address. It returns transactions that have occurred since the last check, ensuring subscribers receive up-to-date information.

- **GetCurrentBlock**: Fetches the current block number from the blockchain. This function is crucial for tracking the latest block and ensuring that the application checks transactions up to the most recent block.

### Continuous Monitoring

The application employs a Go routine that runs every second, checking the latest block on the blockchain. If new blocks have been mined, it checks each subscribed address for new transactions from the last checked block to the current block. New transactions are appended to the respective address's transaction list in the `MemoryStorage`.

### Project Structure

The project is organized into several directories reflecting different aspects of the application:

- **entities/**: Defines data models used throughout the application.
- **handlers/**: Contains the HTTP handlers that manage web requests and responses.
- **insomnia/**: Includes pre-configured requests for use with Insomnia REST client to test API endpoints.
- **interfaces/**: Holds interface definitions for consistent interaction across components.
- **routes/**: Manages the API routes setup.
- **services/**: Core business logic and service layer implementation.
  - **mocks/**: Mock implementations for testing.
- **storages/**: Implementation of storage mechanisms for managing persistent data.
- **main.go**: Entry point of the application.

## Scalability and Storage

The current in-memory storage solution serves the purpose of demonstrating the application's capabilities and provides fast access to data. For production environments, especially those requiring long-term data retention and higher scalability, integrating a database system would be advisable.

## Future Enhancements

- **Database Integration**: Implementing a persistent database to replace or augment in-memory storage.
- **Optimization**: Enhancing the Go routine to manage higher loads and more addresses efficiently.
- **API Security**: Adding authentication and secure communication channels for API access.

This application serves as a foundational platform for building more complex transaction monitoring systems tailored to specific needs, including real-time analytics and automated alerting systems.
