package usecase

import (
	"avito-chat_service/internal/api/contexts"
	"avito-chat_service/internal/api/contexts/chats"
	"avito-chat_service/internal/api/contexts/chats/delivery"
	"avito-chat_service/internal/api/uuidgen"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Create ...
func (u *UseCase) Create(ctx context.Context, req *http.Request,
	nc delivery.NewChat) (uuidgen.UUID, error) {
	uid := u.uuidGen.GenChat(nc.Name, nc.Users)

	err := u.dataStore.SaveChat(ctx, chats.Chat{
		ID:        uid,
		Name:      nc.Name,
		Users:     nc.Users,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return "", contexts.ErrHandlerUseCase(
			contexts.Chats, contexts.Create, req, err)
	}

	return uid, nil
}

// Get ...
func (u *UseCase) Get(ctx *gin.Context, req *http.Request,
	userID uuidgen.UUID) ([]delivery.Chat, error) {
	chats, err := u.dataStore.GetChatsByUserID(ctx, userID)
	if err != nil {
		return nil, contexts.ErrHandlerUseCase(
			contexts.Chats, contexts.Get, req, err)
	}

	return chats, nil
}
