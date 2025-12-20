package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"stocky/internal/models"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetOrCreate(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user := &models.User{}
	err := r.db.GetContext(ctx, user, `
		SELECT id, created_at FROM users WHERE id = $1
	`, userID)

	if err == nil {
		return user, nil
	}

	user.ID = userID
	user.CreatedAt = time.Now()
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO users (id, created_at) VALUES ($1, $2)
	`, userID, user.CreatedAt)

	if err != nil {
		return nil, err
	}

	logrus.WithField("user_id", userID).Info("Created new user")
	return user, nil
}

