package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"stocky/internal/models"
)

type LedgerRepository struct {
	db *sqlx.DB
}

func NewLedgerRepository(db *sqlx.DB) *LedgerRepository {
	return &LedgerRepository{db: db}
}

func (r *LedgerRepository) Create(ctx context.Context, tx *sqlx.Tx, entry *models.LedgerEntry) error {
	query := `
		INSERT INTO ledger_entries (id, event_id, entry_type, symbol, debit, credit, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := tx.ExecContext(ctx, query,
		entry.ID, entry.EventID, entry.EntryType, entry.Symbol,
		entry.Debit, entry.Credit, entry.CreatedAt)
	return err
}

func (r *LedgerRepository) VerifyBalance(ctx context.Context) (bool, error) {
	var result struct {
		TotalDebit  string `db:"total_debit"`
		TotalCredit string `db:"total_credit"`
	}

	err := r.db.GetContext(ctx, &result, `
		SELECT 
			SUM(debit) as total_debit,
			SUM(credit) as total_credit
		FROM ledger_entries
	`)

	if err != nil {
		return false, err
	}

	return result.TotalDebit == result.TotalCredit, nil
}

