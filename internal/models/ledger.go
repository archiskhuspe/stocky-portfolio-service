package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type LedgerEntryType string

const (
	LedgerEntryTypeStock LedgerEntryType = "STOCK"
	LedgerEntryTypeCash  LedgerEntryType = "CASH"
	LedgerEntryTypeFee   LedgerEntryType = "FEE"
)

type LedgerEntry struct {
	ID        uuid.UUID       `db:"id"`
	EventID   uuid.UUID       `db:"event_id"`
	EntryType LedgerEntryType `db:"entry_type"`
	Symbol    *string          `db:"symbol"`
	Debit     decimal.Decimal  `db:"debit"`
	Credit    decimal.Decimal  `db:"credit"`
	CreatedAt time.Time       `db:"created_at"`
}

