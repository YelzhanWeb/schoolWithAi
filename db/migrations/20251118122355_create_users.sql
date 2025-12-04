-- +goose Up
-- +goose StatementBegin
CREATE TYPE user_role AS ENUM ('student', 'teacher', 'admin');

CREATE TABLE users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role user_role NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    avatar_url TEXT NOT NULL DEFAULT 'default_avatar.jpg',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users (email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users CASCADE;

DROP TYPE IF EXISTS user_role;
-- +goose StatementEnd