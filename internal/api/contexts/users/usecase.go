package users

import (
	"avito-chat_service/internal/api/contexts/users/delivery"
	"avito-chat_service/internal/api/uuidgen"
	"context"
	"net/http"
	"time"
)

// UseCase ...
type UseCase interface {
	Add(context.Context, *http.Request, delivery.NewUser) (uuidgen.UUID, error)
}

// User ...
type User struct {
	ID        uuidgen.UUID
	Name      string
	CreatedAt time.Time
}

// InternalErr ...
type (
	InternalErr  struct{}
	NotUniqueErr struct{}
)
