package psql

import (
	"avito-chat_service/internal/api/db"
	"errors"
	"sync"

	"github.com/jackc/pgx/v4"
)

// Store ...
type Store struct {
	conn     *pgx.Conn
	mtx      *sync.RWMutex
	rollback rollback
}

type rollback struct {
	needed bool
}

// New ...
func New(conn *pgx.Conn) (db.Service, error) {
	if conn == nil {
		return nil, errors.New("empty conn")
	}

	return &Store{
		conn: conn,
		mtx:  &sync.RWMutex{},
	}, nil
}

// String ...
func (r *rollback) String() string {
	if r.needed {
		r.needed = false

		return "ROLLBACK;"
	}
	return ""
}
