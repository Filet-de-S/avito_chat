package usecase

import (
	"avito-chat_service/internal/api/contexts"
	"avito-chat_service/internal/api/contexts/messages"
	"avito-chat_service/internal/api/contexts/messages/delivery"
	"avito-chat_service/internal/api/uuidgen"
	"context"
	"net/http"
	"time"
)

// Send ...
func (u *UseCase) Send(ctx context.Context, req *http.Request,
	nm delivery.NewMSG) (uuidgen.UUID, error) {
	msg := messages.MSG{
		ChatID:    nm.Chat,
		AuthorID:  nm.Author,
		Text:      nm.Text,
		CreatedAt: time.Now(),
	}

	msg.ID = u.uuidGen.GenMsg(uuidgen.MSG{
		ChatID:    msg.ChatID,
		AuthorID:  msg.AuthorID,
		Text:      msg.Text,
		CreatedAt: msg.CreatedAt,
	})

	err := u.dataStore.SaveMSG(ctx, msg)
	if err != nil {
		return "", contexts.ErrHandlerUseCase(
			contexts.Messages, contexts.Send, req, err)
	}

	return msg.ID, nil
}

// Get ...
func (u *UseCase) Get(ctx context.Context, req *http.Request,
	chatID uuidgen.UUID) ([]delivery.MSG, error) {
	msgs, err := u.dataStore.GetMSGs(ctx, chatID)
	if err != nil {
		return nil, contexts.ErrHandlerUseCase(
			contexts.Messages, contexts.Get, req, err)
	}

	return msgs, nil
}
