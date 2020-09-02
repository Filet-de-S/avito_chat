package psql

import (
	"avito-chat_service/internal/api/db"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Store ...
type Store struct {
	conn *pgxpool.Pool
}

// New ...
func New(pool *pgxpool.Pool) (db.Service, error) {
	if pool == nil {
		return nil, errors.New("empty pool")
	}

	return &Store{
		conn: pool,
	}, nil
}
