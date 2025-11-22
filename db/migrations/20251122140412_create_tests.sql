-- +goose Up
-- +goose StatementBegin
CREATE TABLE tests (
    id TEXT PRIMARY KEY,
    module_id TEXT UNIQUE REFERENCES modules (id) ON DELETE CASCADE, -- 1 модуль = 1 тест
    title TEXT NOT NULL,
    passing_score INTEGER DEFAULT 70 -- Порог прохождения (70%)
);

CREATE TABLE questions (
    id TEXT PRIMARY KEY,
    test_id TEXT REFERENCES tests (id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    question_type VARCHAR(50) DEFAULT 'single_choice' -- 'single_choice', 'multiple'
);

CREATE TABLE answers (
    id TEXT PRIMARY KEY,
    question_id TEXT REFERENCES questions (id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    is_correct BOOLEAN DEFAULT FALSE
);

CREATE TABLE test_results (
    id TEXT PRIMARY KEY,
    user_id TEXT REFERENCES users (id) ON DELETE CASCADE,
    test_id TEXT REFERENCES tests (id) ON DELETE CASCADE,
    score INTEGER NOT NULL, -- Набрано баллов
    is_passed BOOLEAN NOT NULL, -- Сдал/Не сдал
    attempt_date TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS test_results CASCADE;

DROP TABLE IF EXISTS answers CASCADE;

DROP TABLE IF EXISTS questions CASCADE;

DROP TABLE IF EXISTS tests CASCADE;

-- +goose StatementEnd