package gofrsuuid

import (
	"avito-chat_service/internal/api/uuidgen"

	"github.com/gofrs/uuid"
)

// GenUser ...
func (u *UUIDs) GenUser(name string) uuidgen.UUID {
	return uuid.NewV5(u.User, name).String()
}

// GenChat ...
func (u *UUIDs) GenChat(name string, usersID []uuidgen.UUID) uuidgen.UUID {
	for i := range usersID {
		name += usersID[i]
	}

	return uuid.NewV5(u.User, name).String()
}

// GenMsg ...
func (u *UUIDs) GenMsg(msg uuidgen.MSG) uuidgen.UUID {
	from := msg.ChatID + msg.AuthorID + msg.Text + msg.CreatedAt.String()

	return uuid.NewV5(u.User, from).String()
}
