package gofrsuuid

import (
	"avito-chat_service/internal/api/uuidgen"
	"errors"

	"github.com/gofrs/uuid"
)

// UUIDs ...
type UUIDs struct {
	User uuid.UUID
	Chat uuid.UUID
	Msg  uuid.UUID
}

// New ...
func New(uuids UUIDs) (uuidgen.Service, error) {
	switch {
	case uuids == UUIDs{}:
		return nil, errors.New("empty uuids")
	case uuids.User == uuid.UUID{}:
		return nil, errors.New("empty uuid user")
	case uuids.Chat == uuid.UUID{}:
		return nil, errors.New("empty uuid chat")
	case uuids.Msg == uuid.UUID{}:
		return nil, errors.New("empty uuid msg")
	}
	return &uuids, nil
}
