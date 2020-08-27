BEGIN;

CREATE ROLE service_client WITH
    LOGIN
    NOSUPERUSER
    NOCREATEDB
    NOCREATEROLE
    NOINHERIT
    NOREPLICATION
    CONNECTION LIMIT 2
    PASSWORD '1234';

COMMIT;