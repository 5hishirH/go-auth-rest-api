package user

import (
	"context"
	"database/sql"
	"time"
)

type Repository interface {
	Create(ctx context.Context, u User) error
	FindByEmail(email string) (*User, error)
}

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) Create(u User) error {
	// SQLite specific syntax (uses ? placeholders)
	query := `INSERT INTO users (
		email,
		password_hash,
		full_name,
		profile_pic_name,
		updated_at
	) VALUES (
		$1, $2, $3, $4, $5 
	)`

	if _, err := r.db.Exec(query, u.Email, u.PasswordHash, u.FullName, u.ProfilePicName, time.Now()); err != nil {
		return err
	}

	return nil
}
