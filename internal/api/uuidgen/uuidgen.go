package uuidgen

import (
	"time"
)

// UUID ...
type UUID = string

// Service ...
type Service interface {
	GenUser(name string) UUID
	GenChat(name string, usersID []UUID) UUID
	GenMsg(MSG) UUID
}

// MSG ...
type MSG struct {
	ChatID    UUID
	AuthorID  UUID
	Text      string
	CreatedAt time.Time
}
