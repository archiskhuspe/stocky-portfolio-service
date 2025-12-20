package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"stocky/internal/models"
)

type StockPriceRepository struct {
	db *sqlx.DB
}

func NewStockPriceRepository(db *sqlx.DB) *StockPriceRepository {
	return &StockPriceRepository{db: db}
}

func (r *StockPriceRepository) Upsert(ctx context.Context, price *models.StockPrice) error {
	query := `
		INSERT INTO stock_prices (symbol, price, fetched_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (symbol) DO UPDATE SET
			price = EXCLUDED.price,
			fetched_at = EXCLUDED.fetched_at
	`
	_, err := r.db.ExecContext(ctx, query, price.Symbol, price.Price, price.FetchedAt)
	return err
}

func (r *StockPriceRepository) GetLatest(ctx context.Context, symbol string) (*models.StockPrice, error) {
	price := &models.StockPrice{}
	err := r.db.GetContext(ctx, price, `
		SELECT symbol, price, fetched_at
		FROM stock_prices
		WHERE symbol = $1
		ORDER BY fetched_at DESC
		LIMIT 1
	`, symbol)
	return price, err
}

func (r *StockPriceRepository) GetAllLatest(ctx context.Context) (map[string]decimal.Decimal, error) {
	type result struct {
		Symbol string          `db:"symbol"`
		Price  decimal.Decimal `db:"price"`
	}

	var results []result
	err := r.db.SelectContext(ctx, &results, `
		SELECT DISTINCT ON (symbol) symbol, price
		FROM stock_prices
		ORDER BY symbol, fetched_at DESC
	`)

	if err != nil {
		return nil, err
	}

	prices := make(map[string]decimal.Decimal)
	for _, r := range results {
		prices[r.Symbol] = r.Price
	}
	return prices, nil
}

func (r *StockPriceRepository) GetHistoricalPrices(ctx context.Context, symbol string, date time.Time) (*models.StockPrice, error) {
	price := &models.StockPrice{}
	err := r.db.GetContext(ctx, price, `
		SELECT symbol, price, fetched_at
		FROM stock_prices
		WHERE symbol = $1 AND fetched_at <= $2
		ORDER BY fetched_at DESC
		LIMIT 1
	`, symbol, date)
	return price, err
}

