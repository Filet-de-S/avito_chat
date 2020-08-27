\c msg_service;

BEGIN;

-- User
-- Пользователь приложения. Имеет следующие свойства:
--
-- id - уникальный идентификатор пользователя (может быть как числом, так и строковым – как удобнее)
-- username - уникальное имя пользователя
-- created_at - время создания пользователя
CREATE TABLE IF NOT EXISTS public.users
(
    id TEXT CONSTRAINT "name" PRIMARY KEY,
    name TEXT CONSTRAINT "username" UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

-- Chat
-- Отдельный чат. Имеет следующие свойства:
--
-- id - уникальный идентификатор чата
-- name - уникальное имя чата
-- users - список пользователей в чате, отношение многие-ко-многим
-- created_at - время создания
CREATE TABLE IF NOT EXISTS public.chat
(
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

-- chat users
CREATE TABLE IF NOT EXISTS public.chat_users
(
    chat_id TEXT REFERENCES public.chat(id),
    user_id TEXT CONSTRAINT "user"
                 REFERENCES public.users(id)
                 ON DELETE CASCADE,
    PRIMARY KEY (chat_id, user_id)
);

-- Message
-- Сообщение в чате. Имеет следующие свойства:
--
-- id - уникальный идентификатор сообщения
-- chat - ссылка на идентификатор чата, в который было отправлено сообщение
-- author - ссылка на идентификатор отправителя сообщения, отношение многие-к-одному
-- text - текст отправленного сообщения
-- created_at - время создания
CREATE TABLE IF NOT EXISTS public.msg
(
    id TEXT PRIMARY KEY,
    chat_id TEXT CONSTRAINT "chat"
                 REFERENCES public.chat(id)
                 ON DELETE CASCADE,
    author TEXT CONSTRAINT "author"
                REFERENCES public.users(id)
                ON DELETE CASCADE,
    text TEXT,
    created_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT "author in chat"
        FOREIGN KEY(chat_id, author)
        REFERENCES public.chat_users(chat_id, user_id)
        ON DELETE CASCADE
);

COMMIT;