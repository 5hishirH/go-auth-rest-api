package user

import "time"

type User struct {
	Id                      int64     // `json:"id"`
	Email                   string    // `json:"email" validate:"required"`
	PasswordHash            string    // `json:"password_hash" validate:"required"`
	FullName                string    // `json:"full_name" validate:"required"`
	ProfilePicName          string    // `json:"profile_pic_path" validate:"required"`
	RefreshToken            string    // `json:"-"`
	refreshTokenExipryInSec time.Time // `json:"-"`
	CreatedAt               time.Time //`json:"-"`
}

type RefreshToken struct {
	Id          int64
	UserId      int64
	HashedToken string
	TokenExpiry time.Time
	CreatedAt   time.Time
}
