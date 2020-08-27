package psql

import (
	"avito-chat_service/internal/api/contexts"
	"avito-chat_service/internal/api/contexts/chats"
	"avito-chat_service/internal/api/contexts/messages"
	"avito-chat_service/internal/api/contexts/users"
	"avito-chat_service/internal/api/db"
	"avito-chat_service/internal/api/uuidgen"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

const (
	saveUser = iota
	saveChat
	saveMSG
	getChats
	getMSGs
)

// GetMSGs ...
func (s *Store) GetMSGs(ctx context.Context, chatID uuidgen.UUID) (
	db.Messages, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	query := s.rollback.String() + `SELECT * FROM get_chat_msgs($1);`

	res, err := s.conn.Query(ctx, query, chatID)
	if err != nil {
		s.rollback.needed = true
		return nil, handleError(getMSGs, err)
	}
	defer res.Close()

	msgs := db.Messages{}

	for res.Next() {
		m := db.MSG{}

		err := res.Scan(&m.ID, &m.AuthorID, &m.Text, &m.CreatedAt)
		if err != nil {
			s.rollback.needed = true
			return nil, handleError(getMSGs, err)
		}

		msgs = append(msgs, m)
	}

	err = res.Err()
	if err != nil {
		s.rollback.needed = true
		return nil, handleError(getMSGs, err)
	}

	return msgs, nil
}

// GetChats ...
func (s *Store) GetChats(ctx context.Context, userID uuidgen.UUID) (
	db.UserChats, error) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	chats, usersQuery, err := getChatsWithoutUsers(ctx, userID, s)
	if err != nil {
		s.rollback.needed = true
		return nil, handleError(getChats, err)
	}

	bRes := s.conn.SendBatch(ctx, usersQuery)
	defer func(batchRes pgx.BatchResults) {
		err = batchRes.Close()
		if err != nil {
			s.rollback.needed = true

			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				log.Println("ERROR: pgx.BatchResults.Close() fail in GetChats:", *pgErr)
			} else {
				log.Println("ERROR: pgx.BatchResults.Close() fail in GetChats:", err)
			}
		}
	}(bRes)

	for i := 0; i < usersQuery.Len(); i++ {
		res, err := bRes.Query()
		if err != nil {
			s.rollback.needed = true
			return nil, handleError(getChats, err)
		}
		defer res.Close()

		err = getUsersFromChats(res, &chats, i)
		if err != nil {
			s.rollback.needed = true
			return nil, handleError(getChats, err)
		}
	}

	return chats, nil
}

func getUsersFromChats(res pgx.Rows, chats *db.UserChats, i int) error {
	for res.Next() {
		var userFromChat uuidgen.UUID

		err := res.Scan(&userFromChat)
		if err != nil {
			return nil
		}

		(*chats)[i].Users = append((*chats)[i].Users, userFromChat)
	}

	err := res.Err()
	if err != nil {
		return err
	}

	return nil
}

func getChatsWithoutUsers(ctx context.Context, userID uuidgen.UUID, s *Store) (
	db.UserChats, *pgx.Batch, error) {
	query := s.rollback.String() + `SELECT * FROM get_chats_without_users($1);`

	res, err := s.conn.Query(ctx, query, userID)
	if err != nil {
		return nil, nil, err
	}
	defer res.Close()

	batch := pgx.Batch{}
	bQuery := `SELECT * FROM get_users_from_chat($1);`
	chatsArr := db.UserChats{}

	for res.Next() {
		ch := db.Chat{}

		err = res.Scan(&ch.ID, &ch.Name, &ch.CreatedAt)
		if err != nil {
			return nil, nil, err
		}

		chatsArr = append(chatsArr, ch)
		batch.Queue(bQuery, ch.ID)
	}

	err = res.Err()
	if err != nil {
		return nil, nil, err
	}

	return chatsArr, &batch, nil
}

// SaveUser ...
func (s *Store) SaveUser(ctx context.Context, user users.User) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	_, err := s.conn.Exec(
		ctx,
		s.rollback.String()+`SELECT save_user($1, $2, $3);`,
		user.ID,
		user.Name,
		user.CreatedAt,
	)
	if err != nil {
		s.rollback.needed = true
		return handleError(saveUser, err)
	}

	return nil
}

// SaveChat ...
func (s *Store) SaveChat(ctx context.Context, chat chats.Chat) error {
	query, args := prepareTRQuerySaveChat(chat)

	s.mtx.Lock()
	defer s.mtx.Unlock()

	_, err := s.conn.Exec(
		ctx,
		s.rollback.String()+query,
		args...,
	)
	if err != nil {
		s.rollback.needed = true
		return handleError(saveChat, err)
	}

	return nil
}

func prepareTRQuerySaveChat(chat chats.Chat) (string, []interface{}) {
	query := "BEGIN; " +
		"SELECT save_chat($1, $2, $3);" +
		"INSERT INTO chat_users " +
		"VALUES"
	args := make([]interface{}, 0, 7)
	args = append(args, chat.ID)
	args = append(args, chat.Name)
	args = append(args, chat.CreatedAt)

	i := 4
	for j := range chat.Users {
		args = append(args, chat.ID)
		args = append(args, chat.Users[j])
		query += fmt.Sprintf(" ($%d, $%d),", i, i+1)
		i += 2
	}

	return query[:len(query)-1] + "; COMMIT;", args
}

// SaveMSG ...
func (s *Store) SaveMSG(ctx context.Context, msg messages.MSG) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	_, err := s.conn.Exec(
		ctx,
		s.rollback.String()+`SELECT save_msg($1, $2, $3, $4, $5);`,
		msg.ID,
		msg.ChatID,
		msg.AuthorID,
		msg.Text,
		msg.CreatedAt,
	)
	if err != nil {
		s.rollback.needed = true
		return handleError(saveMSG, err)
	}

	return nil
}

func handleError(module uint8, err error) error {
	chErr := &contexts.Error{}
	chErr.Type = db.InternalErr{}
	chErr.Err = errors.New("ERR FROM from db with method " +
		getMethodName(module) + "; RAW: " + err.Error())

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			if module != saveMSG {
				chErr.Type = db.NotUniqueErr{}
				chErr.UserDescription = "Field '" + pgErr.ConstraintName +
					"' is already exists with that name: '" + pgErr.Detail + "'"
			}
		case pgerrcode.ForeignKeyViolation:
			chErr.Type = db.ForeignKeyViolation{}

			if module == saveChat && pgErr.ConstraintName == "user" {
				val := strings.Split(pgErr.Detail, ")=(")[1]
				value := strings.Split(val, ")")[0]
				chErr.UserDescription = "'" + pgErr.ConstraintName + "' with value '" +
					value + "' seems not to exist"
			} else {
				chErr.UserDescription = "'" + pgErr.ConstraintName + "' with value '" +
					pgErr.Detail + "' seems not to exist"
			}
		case pgerrcode.SyntaxError:
			chErr.Type = db.SyntaxErr{}
		}

		chErr.Err = fmt.Errorf("%s; MORE DETAILs: %+v", chErr.Err, *pgErr)
	}

	return chErr
}

func getMethodName(module uint8) string {
	switch module {
	case saveUser:
		return "SaveUser"
	case saveChat:
		return "SaveChat"
	case saveMSG:
		return "SaveMSG"
	case getChats:
		return "GetChats"
	case getMSGs:
		return "GetMSGs"
	default:
		return "undefined"
	}
}
