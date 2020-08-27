\c msg_service;

BEGIN;

-- get_chats_without_users(user_ID)
CREATE FUNCTION get_chats_without_users(user_id_ TEXT)
    RETURNS TABLE (_id TEXT, _name TEXT, m_created_at TIMESTAMPTZ)
                  AS
$BODY$
BEGIN
IF NOT EXISTS(SELECT id FROM users WHERE id = $1) THEN
    RAISE FOREIGN_KEY_VIOLATION USING CONSTRAINT = 'user', DETAIL = $1;
END IF;

RETURN QUERY
SELECT chat.id, chat.name, chat.created_at
FROM chat
         JOIN chat_users
              ON chat.id = chat_users.chat_id
         LEFT JOIN msg
                   ON chat.id = msg.chat_id
WHERE chat_users.user_id = $1
GROUP BY chat.id
ORDER BY COALESCE(max(msg.created_at), to_timestamp(0)) DESC,
         chat.created_at DESC;
END
$BODY$
    LANGUAGE plpgsql;


-- get_users_from_chat(chatID)
CREATE FUNCTION get_users_from_chat(chat_id_ TEXT)
    RETURNS TABLE (user_id_ TEXT)
                  AS
$BODY$
BEGIN
IF NOT EXISTS(SELECT id FROM chat WHERE id = $1) THEN
    RAISE FOREIGN_KEY_VIOLATION USING CONSTRAINT = 'chat', DETAIL = $1;
END IF;

RETURN QUERY
SELECT user_id
FROM chat_users
WHERE chat_id = $1;
END
$BODY$
    LANGUAGE plpgsql;


-- get_chat_msgs(chatID)
CREATE FUNCTION get_chat_msgs(chat_id_ TEXT)
    RETURNS TABLE (m_id TEXT, m_author TEXT, m_text TEXT, m_created_at TIMESTAMPTZ)
                  AS
$BODY$
BEGIN
IF NOT EXISTS(SELECT id FROM chat WHERE id = $1) THEN
    RAISE FOREIGN_KEY_VIOLATION USING CONSTRAINT = 'chat', DETAIL = $1;
END IF;

RETURN QUERY
    SELECT id, author, text, created_at
    FROM msg
    WHERE chat_id = chat_id_
    ORDER BY created_at ASC;
END
$BODY$
    LANGUAGE plpgsql;


-- save_user(id, name, created_at)
CREATE FUNCTION save_user(user_id_ TEXT, name_ TEXT, created_at_ TIMESTAMPTZ)
    RETURNS VOID
AS
$BODY$
BEGIN
IF EXISTS(SELECT name FROM users WHERE name = name_) THEN
    RAISE UNIQUE_VIOLATION USING CONSTRAINT = 'name', DETAIL = $2;
END IF;

    INSERT INTO users
    VALUES (user_id_, name_, created_at_);
END
$BODY$
    LANGUAGE plpgsql;


-- save_chat(id, name, created_at)
CREATE FUNCTION save_chat(chat_id_ TEXT, name_ TEXT, created_at_ TIMESTAMPTZ)
    RETURNS VOID
AS
$BODY$
BEGIN
IF EXISTS(SELECT name FROM chat WHERE name = name_) THEN
        RAISE UNIQUE_VIOLATION USING CONSTRAINT = 'name', DETAIL = $2;
END IF;

INSERT INTO chat
VALUES (chat_id_, name_, created_at_);
END
$BODY$
    LANGUAGE plpgsql;


-- save_msg(msg_id, chat_id, author_id, text, created_at)
CREATE FUNCTION save_msg(msg_id_ TEXT, chat_id_ TEXT, author_id_ TEXT, text_ TEXT, created_at_ TIMESTAMPTZ)
    RETURNS VOID
AS
$BODY$
BEGIN
IF NOT EXISTS(SELECT chat.id FROM chat WHERE chat.id = chat_id_) THEN
    RAISE FOREIGN_KEY_VIOLATION USING CONSTRAINT = 'chat', DETAIL = chat_id_;
ELSIF NOT EXISTS(SELECT users.id FROM users WHERE users.id = author_id_) THEN
    RAISE FOREIGN_KEY_VIOLATION USING CONSTRAINT = 'author', DETAIL = author_id_;
ELSEIF NOT EXISTS(SELECT user_id FROM chat_users WHERE chat_id = chat_id_ AND user_id = author_id_) THEN
    RAISE FOREIGN_KEY_VIOLATION USING CONSTRAINT = 'author in chat', DETAIL = author_id_;
END IF;

INSERT INTO msg
VALUES ($1, $2, $3, $4, $5);
END
$BODY$
    LANGUAGE plpgsql;

COMMIT;