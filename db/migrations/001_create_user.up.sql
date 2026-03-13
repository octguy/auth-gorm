-- 001_create_users.up.sql
 
CREATE TABLE IF NOT EXISTS users (
    id         UUID    PRIMARY KEY,
    username   VARCHAR(32)  NOT NULL,
    email      VARCHAR(255) NOT NULL UNIQUE,
    password   TEXT         NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);