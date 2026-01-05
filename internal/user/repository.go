package user

import (
	"context"
	"database/sql"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, u User) error {
	// SQLite specific syntax uses ? placeholders
	// postgres specific syntax uses $1, $2, $3, .... placeholders
	query := `INSERT INTO users (
		email,
		role,
		password_hash,
		refresh_token_hash,
		refresh_token_expiry,
		is_verified,
		full_name,
		profile_pic_name,
		created_at,
		updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
	)`

	if _, err := r.db.ExecContext(
		ctx,
		query,
		u.Email,
		u.Role,
		u.PasswordHash,
		u.RefreshTokenHash,
		u.RefreshTokenExpiry,
		u.IsVerified,
		u.FullName,
		u.ProfilePicName,
		u.UpdatedAt,
		u.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	query := `SELECT id, email, role, password_hash, refresh_token_hash, refresh_token_expiry, is_verified, full_name, profile_pic_name, created_at, updated_at
	
	FROM users
	WHERE email = $1`

	row := r.db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.Role,
		&user.PasswordHash,
		&user.RefreshTokenHash,
		&user.RefreshTokenExpiry,
		&user.IsVerified,
		&user.FullName,
		&user.ProfilePicName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
