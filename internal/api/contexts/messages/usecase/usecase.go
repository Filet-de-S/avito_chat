package usecase

import (
	"avito-chat_service/internal/api/contexts/messages"
	"avito-chat_service/internal/api/db"
	"avito-chat_service/internal/api/uuidgen"
	"errors"
)

// UseCase ...
type UseCase struct {
	dataStore db.Service
	uuidGen   uuidgen.Service
}

// New ...
func New(ds db.Service, uuidsGen uuidgen.Service) (messages.UseCase, error) {
	switch {
	case ds == nil:
		return nil, errors.New("dataStore is nil")
	case uuidsGen == nil:
		return nil, errors.New("uuidsGen is nil")
	}

	return &UseCase{
		dataStore: ds,
		uuidGen:   uuidsGen,
	}, nil
}
