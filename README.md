# 🏦 SimpleBank

A modern banking application built with Go, featuring multi-currency transfers, exchange rates, and clean architecture principles.

## 🚀 Features

### Core Banking Features
- **Account Management**: Create and manage bank accounts
- **Multi-Currency Support**: USD, EUR, GBP, NGN
- **Cross-Currency Transfers**: Real-time exchange rate conversion
- **Transaction History**: Complete audit trail of all transfers
- **Balance Tracking**: Real-time account balance updates

### Exchange Rate System
- **Dynamic Exchange Rates**: Real-time currency conversion
- **Rate Calculation**: Calculate amounts to send/receive
- **Transaction Validation**: Verify if currency pairs can be transacted

### Technical Features
- **Clean Architecture**: Separation of concerns with Controllers, Services, Repositories
- **Database Transactions**: Atomic operations for data consistency
- **Precise Financial Calculations**: Using `decimal.Decimal` for accuracy
- **Comprehensive Error Handling**: Custom error types with HTTP status codes
- **Structured Logging**: Consistent logging across all layers
- **Dependency Injection**: Proper service injection and interface-based design

## 🏗️ Architecture

### Clean Architecture Layers

```
┌─────────────────┐
│   Controllers   │  ← HTTP request/response handling
├─────────────────┤
│    Services     │  ← Business logic & validation
├─────────────────┤
│  Repositories   │  ← Data access & persistence
├─────────────────┤
│   Database      │  ← PostgreSQL with SQLC
└─────────────────┘
```

### Application Structure

```
internal/apps/
├── accounts/           # Account management
│   ├── controllers/    # HTTP handlers
│   ├── services/       # Business logic
│   ├── repositories/   # Data access
│   ├── requests/       # Request DTOs
│   ├── responses/      # Response DTOs
│   ├── errors/         # Custom errors
│   └── routes.go       # Route definitions
├── transfers/          # Transfer operations
│   ├── controllers/
│   ├── services/
│   ├── repositories/
│   ├── requests/
│   ├── responses/
│   ├── errors/
│   ├── validationMessages/
│   └── routes.go
└── exchangeRates/      # Exchange rate operations
    ├── controllers/
    ├── services/
    ├── repositories/
    ├── requests/
    ├── responses/
    ├── errors/
    ├── validationMessages/
    └── routes.go
```

## 🛠️ Technology Stack

- **Backend**: Go 1.21+
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL 15+
- **ORM**: SQLC (type-safe SQL)
- **Financial Precision**: `github.com/shopspring/decimal`
- **Validation**: `go-playground/validator`
- **Testing**: `testify` + `gomock`
- **Containerization**: Docker & Docker Compose
- **Migration**: `migrate` tool

## 📋 Prerequisites

- Go 1.21 or higher
- Docker & Docker Compose
- Make (optional, for convenience commands)

## 🚀 Quick Start

### 1. Clone the Repository
```bash
git clone <repository-url>
cd simplebank
```

### 2. Start the Application
```bash
# Start all services (PostgreSQL + API)
make up

# Or start services individually
make postgres
make server
```

### 3. Run Database Migrations
```bash
make migrateup
```

### 4. Generate SQLC Code
```bash
make sqlc
```

### 5. Run Tests
```bash
make test
```

## 📚 API Documentation

### 📖 **Interactive API Documentation**
**📋 [Postman Collection](https://bold-equinox-997479.postman.co/documentation/9123773-c1b33bf1-3b88-4085-800e-16ecd096f5e5/publish?workspaceId=33f35143-e434-4aad-86d4-9151679ef5b6)**

Explore the complete API with interactive examples, request/response samples, and testing capabilities.

### Base URL
```
http://localhost:8080/api/v1
```

### Account Endpoints

#### Create Account
```http
POST /accounts
Content-Type: application/json

{
  "owner": "John Doe",
  "currency": "USD"
}
```

#### Get Account
```http
GET /accounts/{id}
```

### Transfer Endpoints

#### Make Transfer
```http
POST /transfers
Content-Type: application/json

{
  "from_account_id": 1,
  "to_account_id": 2,
  "amount": "100.50",
  "from_currency": "USD",
  "to_currency": "EUR"
}
```

**Sample Responses:**

**Same Currency Transfer:**
```json
{
  "transfer": {
    "id": 1,
    "from_account_id": 1,
    "to_account_id": 2,
    "amount": "100.50",
    "converted_amount": "100.50",
    "from_currency": "USD",
    "to_currency": "USD",
    "exchange_rate": "1.00000000",
    "created_at": "2025-08-02T10:00:00Z"
  },
  "from_account": { ... },
  "to_account": { ... },
  "from_entry": { ... },
  "to_entry": { ... },
  "message": "Transfer completed successfully"
}
```

**Cross-Currency Transfer:**
```json
{
  "transfer": {
    "id": 2,
    "from_account_id": 1,
    "to_account_id": 3,
    "amount": "100.00",
    "converted_amount": "85.00",
    "from_currency": "USD",
    "to_currency": "EUR",
    "exchange_rate": "0.85000000",
    "created_at": "2025-08-02T10:00:00Z"
  },
  "from_account": { ... },
  "to_account": { ... },
  "from_entry": { ... },
  "to_entry": { ... },
  "message": "Transfer completed successfully"
}
```

### Exchange Rate Endpoints

#### List All Exchange Rates
```http
GET /exchange-rates
```

#### Get Exchange Rate for Currency Pair
```http
POST /exchange-rates/calculate
Content-Type: application/json

{
  "from_currency": "USD",
  "to_currency": "EUR",
  "amount": "100.00"
}
```

**Response:**
```json
{
  "exchange_rate": {
    "id": 1,
    "from_currency": "USD",
    "to_currency": "EUR",
    "rate": "0.85000000",
    "created_at": "2025-08-02T10:00:00Z"
  },
  "amount_to_send": "100.00",
  "amount_to_receive": "85.00",
  "can_transact": true,
  "message": "Exchange rate available for transaction"
}
```

## 🗄️ Database Schema

### Tables

#### Accounts
```sql
CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "balance" DECIMAL(20,2) NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
```

#### Transfers
```sql
CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "amount" DECIMAL(20,2) NOT NULL,
  "converted_amount" DECIMAL(20,2),
  "exchange_rate" DECIMAL(20,8),
  "from_currency" VARCHAR(3),
  "to_currency" VARCHAR(3),
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
```

#### Entries
```sql
CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "amount" DECIMAL(20,2) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
```

#### Exchange Rates
```sql
CREATE TABLE exchange_rates (
    id BIGSERIAL PRIMARY KEY,
    from_currency VARCHAR(3) NOT NULL,
    to_currency VARCHAR(3) NOT NULL,
    rate DECIMAL(20,8) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(from_currency, to_currency)
);
```

## 🔧 Development Commands

```bash
# Start services
make up              # Start all services
make down            # Stop all services
make postgres        # Start PostgreSQL only
make server          # Start API server only

# Database operations
make migrateup       # Run migrations
make migratedown     # Rollback migrations
make sqlc            # Generate SQLC code
make sqlcgen         # Generate SQLC code and mocks

# Testing
make test            # Run all tests
make testdb          # Run database tests
make mockgen         # Generate mocks

# Build
make build           # Build the application
make clean           # Clean build artifacts
```

## 🧪 Testing

### Running Tests
```bash
# Run all tests
make test

# Run specific test packages
go test ./internal/apps/accounts/...
go test ./internal/apps/transfers/...
go test ./internal/apps/exchangeRates/...

# Run database tests
make testdb
```

### Test Structure
- **Unit Tests**: Test individual functions and methods
- **Integration Tests**: Test service layer with mocked repositories
- **Database Tests**: Test repository layer with real database
- **Mock Tests**: Test controllers with mocked services

## 🔒 Error Handling

The application uses a comprehensive error handling system:

### Error Types
- **ClientError**: User input validation errors (HTTP 400)
- **ServerError**: Internal server errors (HTTP 500)
- **Custom Errors**: Domain-specific errors with status codes

### Error Response Format
```json
{
  "error": "Error message",
  "status_code": 400,
  "details": "Additional error details"
}
```

## 📊 Logging

The application uses structured logging with `log/slog`:

### Log Levels
- **INFO**: General application flow
- **ERROR**: Error conditions
- **DEBUG**: Detailed debugging information

### Log Format
```json
{
  "time": "2025-08-02T10:00:00Z",
  "level": "INFO",
  "msg": "Processing transfer request",
  "from_account_id": 1,
  "to_account_id": 2,
  "amount": "100.50"
}
```

## 🚀 Deployment

### Docker Deployment
```bash
# Build and run with Docker Compose
docker-compose up -d

# Build production image
docker build -t simplebank:latest .
```

### Environment Variables
```bash
DB_DRIVER=postgres
DB_SOURCE=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable
SERVER_ADDRESS=0.0.0.0:8080
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [SQLC](https://sqlc.dev/) for type-safe SQL
- [Gin](https://gin-gonic.com/) for the HTTP framework
- [ShopSpring Decimal](https://github.com/shopspring/decimal) for precise financial calculations
- [Testify](https://github.com/stretchr/testify) for testing utilities

---

**SimpleBank** - Modern banking with clean architecture! 🏦✨ 