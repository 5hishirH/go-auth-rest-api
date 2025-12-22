package session

import "time"

type Session struct {
	Id     string
	UserId int64
	Expiry time.Time
}
