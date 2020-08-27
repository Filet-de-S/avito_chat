package messages

import (
	"avito-chat_service/internal/api/contexts/messages/delivery"
	"avito-chat_service/internal/api/uuidgen"
	"context"
	"net/http"
	"time"
)

// UseCase ...
type UseCase interface {
	Send(context.Context, *http.Request, delivery.NewMSG) (uuidgen.UUID, error)
	Get(context.Context, *http.Request, uuidgen.UUID) ([]delivery.MSG, error)
}

// MSG ...
type MSG struct {
	ID        uuidgen.UUID
	ChatID    uuidgen.UUID
	AuthorID  uuidgen.UUID
	Text      string
	CreatedAt time.Time
}
