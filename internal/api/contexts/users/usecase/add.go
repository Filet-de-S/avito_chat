package usecase

import (
	"avito-chat_service/internal/api/contexts"
	"avito-chat_service/internal/api/contexts/users"
	"avito-chat_service/internal/api/contexts/users/delivery"
	"avito-chat_service/internal/api/uuidgen"
	"context"
	"net/http"
	"time"
)

// Add ...
func (u *UseCase) Add(ctx context.Context, req *http.Request,
	nu delivery.NewUser) (uuidgen.UUID, error) {
	uid := u.uuidGen.GenUser(nu.Username)

	err := u.dataStore.SaveUser(ctx, users.User{
		ID:        uid,
		Name:      nu.Username,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return "", contexts.ErrHandlerUseCase(
			contexts.Users, contexts.Add, req, err)
	}

	return uid, nil
}
