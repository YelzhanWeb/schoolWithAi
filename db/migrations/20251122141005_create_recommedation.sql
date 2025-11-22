-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_activity_logs (
    id TEXT PRIMARY KEY,
    user_id TEXT REFERENCES users (id) ON DELETE CASCADE,
    course_id TEXT REFERENCES courses (id),
    action_type TEXT NOT NULL, -- 'view', 'complete', 'like', 'fail_test'
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE recommendations (
    id TEXT PRIMARY KEY,
    user_id TEXT REFERENCES users (id) ON DELETE CASCADE,
    course_id TEXT REFERENCES courses (id) ON DELETE CASCADE,
    score FLOAT NOT NULL, -- Вес рекомендации (0.95)
    reason TEXT, -- "Потому что вы прошли Python Basic"
    created_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS recommendations CASCADE;

DROP TABLE IF EXISTS user_activity_logs CASCADE;
-- +goose StatementEnd