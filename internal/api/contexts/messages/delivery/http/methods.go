package http

import (
	"avito-chat_service/internal/api/contexts"
	"avito-chat_service/internal/api/contexts/messages/delivery"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (e *engine) send(ctx *gin.Context) {
	msg := delivery.NewMSG{}

	err := ctx.ShouldBindJSON(&msg)
	if err != nil {
		resp := contexts.ErrorHandler(
			contexts.Messages,
			contexts.Send,
			ctx.Request,
			contexts.ParseErr(err),
		)
		ctx.JSON(resp.Errors[0].Status, resp)

		return
	}

	id, err := e.useCase.Send(ctx, ctx.Request, msg)
	if err != nil {
		resp := contexts.ErrorHandler(
			contexts.Messages,
			contexts.Send,
			ctx.Request,
			[]error{err},
		)
		ctx.JSON(resp.Errors[0].Status, resp)

		return
	}

	ctx.JSON(http.StatusCreated, delivery.MSGSend{ID: id})
}

func (e *engine) get(ctx *gin.Context) {
	msg := delivery.GetMSG{}

	err := ctx.ShouldBindJSON(&msg)
	if err != nil {
		resp := contexts.ErrorHandler(
			contexts.Messages,
			contexts.Get,
			ctx.Request,
			contexts.ParseErr(err),
		)
		ctx.JSON(resp.Errors[0].Status, resp)

		return
	}

	msgs, err := e.useCase.Get(ctx, ctx.Request, msg.Chat)
	if err != nil {
		resp := contexts.ErrorHandler(
			contexts.Messages,
			contexts.Get,
			ctx.Request,
			[]error{err},
		)
		ctx.JSON(resp.Errors[0].Status, resp)

		return
	}

	if len(msgs) == 0 {
		ctx.Status(http.StatusNoContent)

		return
	}

	ctx.JSON(http.StatusOK, msgs)
}
