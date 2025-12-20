package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"stocky/internal/models"
)

type RewardRepository struct {
	db *sqlx.DB
}

func NewRewardRepository(db *sqlx.DB) *RewardRepository {
	return &RewardRepository{db: db}
}

func (r *RewardRepository) GetByEventID(ctx context.Context, eventID uuid.UUID) (*models.RewardEvent, error) {
	reward := &models.RewardEvent{}
	err := r.db.GetContext(ctx, reward, `
		SELECT id, event_id, user_id, stock_symbol, quantity, timestamp, created_at
		FROM reward_events WHERE event_id = $1
	`, eventID)
	return reward, err
}

func (r *RewardRepository) Create(ctx context.Context, tx *sqlx.Tx, reward *models.RewardEvent) error {
	query := `
		INSERT INTO reward_events (id, event_id, user_id, stock_symbol, quantity, timestamp, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := tx.ExecContext(ctx, query,
		reward.ID, reward.EventID, reward.UserID, reward.StockSymbol,
		reward.Quantity, reward.Timestamp, reward.CreatedAt)
	return err
}

func (r *RewardRepository) GetTodayRewards(ctx context.Context, userID uuid.UUID, istDate time.Time) ([]models.RewardEvent, error) {
	startOfDay := time.Date(istDate.Year(), istDate.Month(), istDate.Day(), 0, 0, 0, 0, istDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var rewards []models.RewardEvent
	err := r.db.SelectContext(ctx, &rewards, `
		SELECT id, event_id, user_id, stock_symbol, quantity, timestamp, created_at
		FROM reward_events
		WHERE user_id = $1 AND timestamp >= $2 AND timestamp < $3
		ORDER BY timestamp DESC
	`, userID, startOfDay, endOfDay)
	return rewards, err
}

func (r *RewardRepository) GetTotalSharesByStockToday(ctx context.Context, userID uuid.UUID, istDate time.Time) (map[string]decimal.Decimal, error) {
	startOfDay := time.Date(istDate.Year(), istDate.Month(), istDate.Day(), 0, 0, 0, 0, istDate.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	type result struct {
		StockSymbol string          `db:"stock_symbol"`
		TotalQty    decimal.Decimal `db:"total_quantity"`
	}

	var results []result
	err := r.db.SelectContext(ctx, &results, `
		SELECT stock_symbol, SUM(quantity) as total_quantity
		FROM reward_events
		WHERE user_id = $1 AND timestamp >= $2 AND timestamp < $3
		GROUP BY stock_symbol
	`, userID, startOfDay, endOfDay)

	if err != nil {
		return nil, err
	}

	totals := make(map[string]decimal.Decimal)
	for _, r := range results {
		totals[r.StockSymbol] = r.TotalQty
	}
	return totals, nil
}

func (r *RewardRepository) GetTotalSharesByStockUpToDate(ctx context.Context, userID uuid.UUID, endDate time.Time) (map[string]decimal.Decimal, error) {
	type result struct {
		StockSymbol string          `db:"stock_symbol"`
		TotalQty    decimal.Decimal `db:"total_quantity"`
	}

	var results []result
	err := r.db.SelectContext(ctx, &results, `
		SELECT stock_symbol, SUM(quantity) as total_quantity
		FROM reward_events
		WHERE user_id = $1 AND timestamp <= $2
		GROUP BY stock_symbol
	`, userID, endDate)

	if err != nil {
		return nil, err
	}

	totals := make(map[string]decimal.Decimal)
	for _, r := range results {
		totals[r.StockSymbol] = r.TotalQty
	}
	return totals, nil
}

