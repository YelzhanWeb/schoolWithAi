-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_activity_logs (
    id TEXT PRIMARY KEY,
    user_id TEXT REFERENCES users (id) ON DELETE CASCADE,
    -- Сделали NULLABLE, чтобы можно было логировать просто поисковые запросы
    course_id TEXT REFERENCES courses (id) ON DELETE SET NULL,
    -- Тип действия: 'view', 'complete', 'search', 'click_tag'  
    action_type VARCHAR(50) NOT NULL,
    -- Мета-данные (JSON). Например:
    -- {"duration": 120} (время просмотра)
    -- {"query": "python"} (что искал)
    meta_data JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_logs_user ON user_activity_logs (user_id);

CREATE INDEX idx_logs_course ON user_activity_logs (course_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS user_activity_logs CASCADE;
-- +goose StatementEnd