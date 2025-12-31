package user

import "time"

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type User struct {
	Id                 int64
	Email              string
	PasswordHash       string
	Role               UserRole
	FullName           string
	ProfilePicName     string
	RefreshToken       string
	RefreshTokenExpiry time.Time
	CreatedAt          time.Time
}
