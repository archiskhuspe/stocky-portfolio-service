package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type RewardEvent struct {
	ID          uuid.UUID       `db:"id"`
	EventID     uuid.UUID       `db:"event_id"`
	UserID      uuid.UUID       `db:"user_id"`
	StockSymbol string          `db:"stock_symbol"`
	Quantity    decimal.Decimal `db:"quantity"`
	Timestamp   time.Time       `db:"timestamp"`
	CreatedAt   time.Time       `db:"created_at"`
}

