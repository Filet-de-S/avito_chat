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

const (
	getChatsByUserIDQ = `SELECT * FROM get_chats($1)`
	getMSGsByChatIDQ  = `SELECT * FROM get_chat_msgs($1)`
	saveUserQ         = `SELECT save_user($1, $2, $3)`
	saveChatQ         = `SELECT save_chat($1, $2, $3);` +
		` INSERT INTO chat_users VALUES`
	saveMSGQ = `SELECT save_msg($1, $2, $3, $4, $5)`
)

// GetMSGs ...
func (s *Store) GetMSGsByChatID(ctx context.Context, chatID uuidgen.UUID) (
	db.Messages, error) {
	res, err := s.conn.Query(ctx, getMSGsByChatIDQ, chatID)
	if err != nil {
		return nil, handleError(getMSGs, err)
	}
	defer res.Close()

	msgs := db.Messages{}

	for res.Next() {
		m := db.MSG{}

		err := res.Scan(&m.ID, &m.AuthorID, &m.Text, &m.CreatedAt)
		if err != nil {
			return nil, handleError(getMSGs, err)
		}

		msgs = append(msgs, m)
	}

	if res.Err() != nil {
		return nil, handleError(getMSGs, res.Err())
	}

	return msgs, nil
}

// GetChats ...
func (s *Store) GetChatsByUserID(ctx context.Context, userID uuidgen.UUID) (
	db.UserChats, error) {
	res, err := s.conn.Query(ctx, getChatsByUserIDQ, userID)
	if err != nil {
		return nil, handleError(getChats, err)
	}
	defer res.Close()

	chats := db.UserChats{}

	for res.Next() {
		chat := db.Chat{}

		err = res.Scan(&chat.ID, &chat.Name, &chat.CreatedAt, &chat.Users)
		if err != nil {
			return nil, handleError(getChats, err)
		}

		chats = append(chats, chat)
	}

	if res.Err() != nil {
		return nil, handleError(getChats, res.Err())
	}

	return chats, nil
}

// SaveUser ...
func (s *Store) SaveUser(ctx context.Context, user users.User) error {
	_, err := s.conn.Exec(
		ctx,
		saveUserQ,
		user.ID,
		user.Name,
		user.CreatedAt,
	)
	if err != nil {
		return handleError(saveUser, err)
	}

	return nil
}

// SaveChat ...
func (s *Store) SaveChat(ctx context.Context, chat chats.Chat) error {
	query, args := getSaveChatQuery(chat)

	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return handleError(saveChat, err)
	}

	_, err = tx.Exec(ctx, query, args...)
	if err != nil {
		return handleError(saveChat, err)
	}

	err = tx.Commit(ctx)
	if err != nil && !errors.Is(err, pgx.ErrTxCommitRollback) {
		return handleError(saveChat, err)
	}

	return nil
}

func getSaveChatQuery(chat chats.Chat) (string, []interface{}) {
	query := saveChatQ

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

	return query[:len(query)-1], args
}

// SaveMSG ...
func (s *Store) SaveMSG(ctx context.Context, msg messages.MSG) error {
	_, err := s.conn.Exec(
		ctx,
		saveMSGQ,
		msg.ID,
		msg.ChatID,
		msg.AuthorID,
		msg.Text,
		msg.CreatedAt,
	)
	if err != nil {
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
		return "GetChatsByUserID"
	case getMSGs:
		return "GetMSGsByChatID"
	default:
		return "undefined"
	}
}
