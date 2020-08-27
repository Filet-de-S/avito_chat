package chats

import (
	"avito-chat_service/internal/api/contexts/chats/delivery"
	"avito-chat_service/internal/api/uuidgen"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// UseCase ...
type UseCase interface {
	Create(context.Context, *http.Request, delivery.NewChat) (uuidgen.UUID, error)
	Get(ctx *gin.Context, request *http.Request, userID uuidgen.UUID) (
		[]delivery.Chat, error)
}

// Chat ...
type Chat struct {
	ID        uuidgen.UUID
	Name      string
	Users     []uuidgen.UUID
	CreatedAt time.Time
}
