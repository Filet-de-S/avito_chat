package router

import (
	"avito-chat_service/internal/api/contexts/chats"
	chatsDelivery "avito-chat_service/internal/api/contexts/chats/delivery/http"
	"avito-chat_service/internal/api/contexts/messages"
	msgDelivery "avito-chat_service/internal/api/contexts/messages/delivery/http"
	"avito-chat_service/internal/api/contexts/users"
	usersDelivery "avito-chat_service/internal/api/contexts/users/delivery/http"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Engine ...
type Engine struct {
	gin *gin.Engine
}

// UseCases ...
type UseCases struct {
	Users users.UseCase
	Chats chats.UseCase
	MSG   messages.UseCase
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// New ...
func New(useCases *UseCases) (*Engine, error) {
	router := Engine{
		gin: gin.Default(),
	}

	if gin.Mode() == "debug" {
		router.gin.Use(logger())
	}

	router.gin.GET("/status", router.status)

	users := router.gin.Group("/users")
	if err := usersDelivery.SetHandlers(users, useCases.Users); err != nil {
		return nil, fmt.Errorf("set users handlers error: %w", err)
	}
	chats := router.gin.Group("/chats")
	if err := chatsDelivery.SetHandlers(chats, useCases.Chats); err != nil {
		return nil, fmt.Errorf("set chats handlers error: %w", err)
	}
	msg := router.gin.Group("/messages")
	if err := msgDelivery.SetHandlers(msg, useCases.MSG); err != nil {
		return nil, fmt.Errorf("set msg handlers error: %w", err)
	}

	return &router, nil
}

// Handler ...
func (e *Engine) Handler() http.Handler {
	return e.gin
}

func (e *Engine) status(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

func logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		b := ctx.Request.Body
		req, err := ioutil.ReadAll(b)
		if err != nil {
			log.Println("DEBUG: logger/ERROR: can't read body", err, "\nREQ:",
				ctx.Request, "\nEND DEBUG MSG")
			return
		}
		defer func() {
			if err := b.Close(); err != nil {
				log.Println("DEBUG: logger/ERROR: can't close REQ.BODY:",
					err, "\nREQ:", ctx.Request, "\nEND DEBUG MSG")
				return
			}
		}()

		resp := fmt.Sprintln("DEBUG: logger/REQ.BODY:\n"+string(req),
			"\nEND OF REQ.BODY, other fields:\n", ctx.Request)

		ctx.Request.Body = ioutil.NopCloser(bytes.NewReader(req))
		respBW := &responseBodyWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: ctx.Writer,
		}
		ctx.Writer = respBW
		ctx.Next()

		fmt.Print(resp, "\nRESP.B:\n"+respBW.body.String(),
			"\nEND OF RESP.B AND DEBUG MSG\n")
	}
}

// Write ...
func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteString ...
func (w responseBodyWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

//func contentType() gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//		//contentTypeH := ctx.GetHeader("Content-Type")
//		//acceptH := ctx.GetHeader("Accept")
//		//if contentTypeH != "application/vnd.api+json" {
//		//	resp := jsonapi.ResponseObject{
//		//		Errors: []jsonapi.ErrorObject{{
//		//			Status: http.StatusUnsupportedMediaType,
//		//			Title:  "Invalid 'Content-Type' header",
//		//			Detail: "want 'application/vnd.api+json'",
//		//		}},
//		//	}
//		//	ctx.JSONAPI(http.StatusUnsupportedMediaType, resp)
//		//	log.Println("ERROR FROM JSONAPI HEADER MIDDLEWARE:", resp, ctx.Request)
//		//} else if acceptH != "" && acceptH != "*/*" &&
//				acceptH != "application/vnd.api+json" {
//		//	resp := jsonapi.ResponseObject{
//		//		Errors: []jsonapi.ErrorObject{{
//		//			Status: http.StatusNotAcceptable,
//		//			Title:  "Service doesn't support this type of 'Accept' header",
//		//			Detail: "accept 'application/vnd.api+json' or nothing",
//		//		}},
//		//	}
//		//	ctx.JSONAPI(http.StatusNotAcceptable, resp)
//		//	log.Println("ERROR FROM JSONAPI HEADER MIDDLEWARE:", resp, ctx.Request)
//		//}
//	}
//}

// to routerInit
//if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
//	err := v.RegisterValidation("jsonapi", JsonAPIFormatValidation, true)
//	if err != nil {
//		return nil, fmt.Errorf("jsonapi register validation error]: %w", err)
//	}
//}
