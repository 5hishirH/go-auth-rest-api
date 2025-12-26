package auth

import (
	"context"
	"database/sql"
	"time"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, userId int64, tokenHash string, duration string) (*time.Time, error) {
	query := `INSERT INTO refresh_tokens (
		user_id,
		token_hash,
		token_expiry,
		updated_at
	) VALUES (
		$1, $2, $3, $4 
	)`

	parsedDuration, err := time.ParseDuration(duration)
	if err != nil {
		return nil, err
	}

	expiry := time.Now().Add(parsedDuration)

	if _, err := r.db.ExecContext(ctx, query, userId, tokenHash, expiry, time.Now()); err != nil {
		return nil, err
	}

	return &expiry, nil
}
