package user

import "time"

// type UserRole string

// const (
// 	RoleUser  UserRole = "user"
// 	RoleAdmin UserRole = "admin"
// )

type User struct {
	Id                 int64
	Email              string
	Role               string
	PasswordHash       string
	RefreshTokenHash   string
	RefreshTokenExpiry time.Time
	IsVerified         bool
	FullName           string
	ProfilePicName     string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
