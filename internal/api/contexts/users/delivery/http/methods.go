package http

import (
	"avito-chat_service/internal/api/contexts"
	"avito-chat_service/internal/api/contexts/users/delivery"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (e *engine) add(ctx *gin.Context) {
	user := delivery.NewUser{}

	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		resp := contexts.ErrorHandler(
			contexts.Users,
			contexts.Add,
			ctx.Request,
			contexts.ParseErr(err),
		)
		ctx.JSON(resp.Errors[0].Status, resp)

		return
	}

	id, err := e.useCase.Add(ctx, ctx.Request, user)
	if err != nil {
		resp := contexts.ErrorHandler(
			contexts.Users,
			contexts.Add,
			ctx.Request,
			[]error{err},
		)
		ctx.JSON(resp.Errors[0].Status, resp)

		return
	}

	ctx.JSON(http.StatusCreated, delivery.UserCreated{ID: id})
}
