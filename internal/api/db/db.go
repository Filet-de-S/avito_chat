package db

import (
	"avito-chat_service/internal/api/contexts/chats"
	chatsDelivery "avito-chat_service/internal/api/contexts/chats/delivery"
	"avito-chat_service/internal/api/contexts/messages"
	msgDelivery "avito-chat_service/internal/api/contexts/messages/delivery"
	"avito-chat_service/internal/api/contexts/users"
	"avito-chat_service/internal/api/uuidgen"
	"context"
)

// Service ...
type Service interface {
	SaveUser(context.Context, users.User) error
	SaveChat(context.Context, chats.Chat) error
	SaveMSG(context.Context, messages.MSG) error
	GetChatsByUserID(ctx context.Context, userID uuidgen.UUID) (UserChats, error)
	GetMSGsByChatID(ctx context.Context, chatID uuidgen.UUID) (Messages, error)
}

// UserChats ...
type (
	UserChats []Chat
	Messages  []MSG
)

// Chat ...
type (
	Chat  = chatsDelivery.Chat
	MSG   = msgDelivery.MSG
	Query = string
)

// InternalErr ...
type (
	InternalErr         struct{}
	NotUniqueErr        struct{}
	SyntaxErr           struct{}
	ForeignKeyViolation struct{}
	NotFound            struct{}
)
