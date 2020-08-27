package http

import (
	"avito-chat_service/internal/api/contexts/users"
	"errors"

	"github.com/gin-gonic/gin"
)

type engine struct {
	useCase users.UseCase
}

// SetHandlers ...
func SetHandlers(g gin.IRouter, useCase users.UseCase) error {
	if useCase == nil {
		return errors.New("empty users useCase")
	}

	e := &engine{
		useCase: useCase,
	}

	g.POST("/add", e.add)
	return nil
}
