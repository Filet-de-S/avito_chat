package http

import (
	"avito-chat_service/internal/api/contexts/messages"
	"errors"

	"github.com/gin-gonic/gin"
)

type engine struct {
	useCase messages.UseCase
}

// SetHandlers ...
func SetHandlers(g gin.IRouter, useCase messages.UseCase) error {
	if useCase == nil {
		return errors.New("empty messages usecase")
	}
	e := &engine{
		useCase: useCase,
	}

	g.POST("/add", e.send)
	g.POST("/get", e.get)

	return nil
}
