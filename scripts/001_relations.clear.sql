BEGIN;

DELETE FROM users;
DELETE FROM chat;
DELETE FROM chat_users;
DELETE FROM msg;

COMMIT;