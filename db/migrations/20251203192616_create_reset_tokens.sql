-- +goose Up
-- +goose StatementBegin
CREATE TABLE password_reset_tokens (
    email TEXT NOT NULL,
    token VARCHAR(6) NOT NULL, -- 6-значный код
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (email) -- Один активный код на email
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS password_reset_tokens;
-- +goose StatementEnd