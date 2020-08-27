package delivery

import (
	"avito-chat_service/internal/api/uuidgen"
	"time"
)

// NewChat ...
type NewChat struct {
	Name  string   `json:"name" binding:"required"`
	Users []string `json:"users" binding:"min=2,unique,dive,uuid5"`
}

// ChatCreated ...
type ChatCreated struct {
	ID string `json:"id"`
}

// GetChats ...
type GetChats struct {
	User string `json:"user" binding:"required,uuid5"`
}

// Chat ...
type Chat struct {
	ID        uuidgen.UUID   `json:"id"`
	Name      string         `json:"name"`
	Users     []uuidgen.UUID `json:"users"`
	CreatedAt time.Time      `json:"created_at"`
}
