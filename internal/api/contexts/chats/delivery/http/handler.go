package http

import (
	"avito-chat_service/internal/api/contexts/chats"
	"errors"

	"github.com/gin-gonic/gin"
)

type engine struct {
	useCase chats.UseCase
}

// SetHandlers ...
func SetHandlers(g gin.IRouter, useCase chats.UseCase) error {
	if useCase == nil {
		return errors.New("empty chats usecase")
	}

	e := &engine{
		useCase: useCase,
	}

	g.POST("/add", e.create)
	g.POST("/get", e.get)

	return nil
}
