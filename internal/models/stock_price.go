package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type StockPrice struct {
	Symbol    string          `db:"symbol"`
	Price     decimal.Decimal `db:"price"`
	FetchedAt time.Time       `db:"fetched_at"`
}

