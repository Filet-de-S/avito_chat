package http

import (
	"avito-chat_service/internal/api/contexts"
	"avito-chat_service/internal/api/contexts/chats/delivery"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (e *engine) create(ctx *gin.Context) {
	chat := delivery.NewChat{}

	err := ctx.ShouldBindJSON(&chat)
	if err != nil {
		resp := contexts.ErrorHandler(
			contexts.Chats,
			contexts.Create,
			ctx.Request,
			contexts.ParseErr(err),
		)
		ctx.JSON(resp.Errors[0].Status, resp)

		return
	}

	id, err := e.useCase.Create(ctx, ctx.Request, chat)
	if err != nil {
		resp := contexts.ErrorHandler(
			contexts.Chats,
			contexts.Create,
			ctx.Request,
			[]error{err},
		)
		ctx.JSON(resp.Errors[0].Status, resp)

		return
	}

	ctx.JSON(http.StatusCreated, delivery.ChatCreated{ID: id})
}

func (e *engine) get(ctx *gin.Context) {
	gc := delivery.GetChats{}

	err := ctx.ShouldBindJSON(&gc)
	if err != nil {
		resp := contexts.ErrorHandler(
			contexts.Chats,
			contexts.Get,
			ctx.Request,
			contexts.ParseErr(err),
		)
		ctx.JSON(resp.Errors[0].Status, resp)

		return
	}

	chats, err := e.useCase.Get(ctx, ctx.Request, gc.User)
	if err != nil {
		resp := contexts.ErrorHandler(
			contexts.Chats,
			contexts.Get,
			ctx.Request,
			[]error{err},
		)
		ctx.JSON(resp.Errors[0].Status, resp)

		return
	}

	if len(chats) == 0 {
		ctx.Status(http.StatusNoContent)

		return
	}

	ctx.JSON(http.StatusOK, chats)
}
