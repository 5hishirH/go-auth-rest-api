package types

import "encoding/gob"

type UserSession struct {
	UserID int64
	Role   string
}

// RegisterTypes ensures 'gob' knows how to encode this struct.
// Call this once in your app startup.
func RegisterTypes() {
	gob.Register(UserSession{})
}
