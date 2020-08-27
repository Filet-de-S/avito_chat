package delivery

import (
	"avito-chat_service/internal/api/uuidgen"
	"time"
)

// NewMSG ...
type NewMSG struct {
	Chat   uuidgen.UUID `json:"chat" binding:"required,uuid5"`
	Author uuidgen.UUID `json:"author" binding:"required,uuid5"`
	Text   string       `json:"text" binding:"required"`
}

// MSGSend ...
type MSGSend struct {
	ID string `json:"id"`
}

// GetMSG ...
type GetMSG struct {
	Chat uuidgen.UUID `json:"chat" binding:"required,uuid5"`
}

// MSG ...
type MSG struct {
	ID        uuidgen.UUID `json:"id"`
	AuthorID  uuidgen.UUID `json:"author"`
	Text      string       `json:"text"`
	CreatedAt time.Time    `json:"created_at"`
}
