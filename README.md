# Stocky Portfolio Service by Archis Khuspe

A production-grade backend service for managing stock rewards, built with Go, Gin, and PostgreSQL.

## Architecture

### Layered Architecture

- **Handler Layer**: HTTP request/response handling, validation
- **Service Layer**: Business logic, orchestration
- **Repository Layer**: Database access abstraction
- **Model Layer**: Domain entities

### Key Components

1. **Double-Entry Accounting**: All transactions are recorded in a ledger with balanced debits and credits
2. **Idempotency**: Reward events use `event_id` for idempotent processing
3. **Price Service**: Hourly background job fetches and stores stock prices
4. **Decimal Precision**: Uses `shopspring/decimal` for accurate financial calculations

## Database Schema

### Tables

- `users`: User records
- `reward_events`: Reward transactions with idempotency
- `ledger_entries`: Double-entry accounting records
- `stock_prices`: Latest stock prices with timestamps

### Ledger Logic

Every reward creates three ledger entries:

1. **Credit STOCK**: Increases stock inventory asset
2. **Debit CASH**: Decreases cash asset
3. **Debit FEE**: Records transaction fees

The ledger always balances: Total Debit = Total Credit

## API Endpoints

### 1. POST /api/v1/reward

Record a stock reward for a user.

**Request:**

```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "stock_symbol": "RELIANCE",
  "quantity": 1.25,
  "timestamp": "2025-01-15T10:30:00Z",
  "event_id": "660e8400-e29b-41d4-a716-446655440000"
}
```

**Response:** 201 Created

```json
{
  "message": "Reward processed successfully",
  "event_id": "660e8400-e29b-41d4-a716-446655440000"
}
```

### 2. GET /api/v1/today-stocks/{userId}

Get all stock rewards for today (IST).

**Response:** 200 OK

```json
[
  {
    "stock_symbol": "RELIANCE",
    "quantity": 1.25,
    "timestamp": "2025-01-15T10:30:00Z",
    "event_id": "660e8400-e29b-41d4-a716-446655440000"
  }
]
```

### 3. GET /api/v1/historical-inr/{userId}

Get INR valuation for each past day (up to yesterday).

**Response:** 200 OK

```json
[
  {
    "date": "2025-01-14",
    "inr_value": 15432.25
  },
  {
    "date": "2025-01-15",
    "inr_value": 16211.9
  }
]
```

### 4. GET /api/v1/stats/{userId}

Get today's shares and current portfolio value.

**Response:** 200 OK

```json
{
  "today_shares_by_stock": {
    "RELIANCE": 1.25,
    "TCS": 0.5
  },
  "current_portfolio_value": 4375.0
}
```

### 5. GET /api/v1/portfolio/{userId}

Get detailed portfolio holdings.

**Response:** 200 OK

```json
[
  {
    "stock_symbol": "RELIANCE",
    "total_quantity": 1.25,
    "current_price": 2500.0,
    "current_value": 3125.0
  }
]
```

## Setup

### Prerequisites

- Go 1.21+
- PostgreSQL 12+
- Make (optional)

### Installation

1. Clone the repository

```bash
git clone <repo-url>
cd project
```

2. Install dependencies

```bash
go mod download
```

3. Setup database

```bash
createdb assignment
psql assignment < migrations/001_initial_schema.sql
```

4. Configure environment

```bash
cp .env.example .env
# Edit .env with your database credentials
```

5. Run the server

```bash
go run cmd/server/main.go
```

## Edge Cases Handled

### 1. Idempotency

- Duplicate `event_id` requests are ignored
- Returns success without creating duplicate records

### 2. Price API Downtime

- Uses last known price with warning log
- Background job retries on next interval

### 3. Rounding Errors

- Uses banker's rounding (Round method from decimal library)
- All calculations use NUMERIC type in PostgreSQL

### 4. Stock Splits/Mergers

- Schema supports symbol updates
- Historical prices preserved for audit

### 5. Delisted Stocks

- Can mark stocks inactive (extensible design)
- Frozen prices prevent valuation errors

### 6. Reward Reversal

- Can create negative ledger entries
- Maintains double-entry balance

## Fee Calculation

Fees include:

- Brokerage: 0.03% (min 20 INR)
- STT: 0.025% on delivery
- GST: 18% on brokerage
- Exchange charges: 0.00325%
- SEBI charges: 0.0001%
- Stamp duty: 0.003%

Total fees are calculated internally and debited separately in the ledger.

## Background Jobs

### Price Fetcher

- Runs hourly (configurable via `PRICE_FETCH_INTERVAL`)
- Fetches prices for all tracked stocks
- Updates `stock_prices` table
- Handles API failures gracefully

## Testing

```bash
# Run tests
go test ./...

# With coverage
go test -cover ./...
```

## Logging

Uses `logrus` with JSON formatting. Logs include:

- Request/response details
- Business events (rewards, price updates)
- Errors with context

## Production Considerations

1. **Database Connection Pooling**: Configured via sqlx
2. **Transaction Management**: All reward processing uses transactions
3. **Error Handling**: Comprehensive error messages with proper HTTP status codes
4. **Graceful Shutdown**: Handles SIGINT/SIGTERM
5. **Health Checks**: `/health` endpoint for monitoring

## License

MIT
