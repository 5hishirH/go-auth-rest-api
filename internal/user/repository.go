package user

import (
	"context"
	"database/sql"
)

type Repository interface {
	Create(ctx context.Context, u User) error
	FindByEmail(email string) (*User, error)
	SaveRefreshToken(userId int64, token string) (int, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, u User) error {
	// SQLite specific syntax (uses ? placeholders)
	query := `INSERT INTO users (email, password_hash, full_name, profile_pic_name) VALUES (?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query, u.Email, u.PasswordHash, u.FullName, u.ProfilePicName)
	return err
}
