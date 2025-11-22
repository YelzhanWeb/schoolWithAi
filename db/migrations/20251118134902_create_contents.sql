-- +goose Up
-- +goose StatementBegin

CREATE TABLE courses (
    id TEXT PRIMARY KEY,
    author_id TEXT REFERENCES users (id),
    subject_id TEXT REFERENCES subjects (id),
    title TEXT NOT NULL,
    description TEXT,
    difficulty_level INTEGER CHECK (
        difficulty_level >= 1
        AND difficulty_level <= 5
    ),
    tags TEXT [], -- Теги для Python ['python', 'backend', 'loops']
    cover_image_url VARCHAR(255),
    is_published BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE modules (
    id TEXT PRIMARY KEY,
    course_id TEXT REFERENCES courses (id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    order_index INTEGER NOT NULL
);

CREATE TABLE lessons (
    id TEXT PRIMARY KEY,
    module_id TEXT REFERENCES modules (id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    content_text TEXT, -- Markdown текст урока
    video_url TEXT, -- Ссылка на видео (MinIO)
    file_attachment_url TEXT, -- Ссылка на PDF/ZIP (MinIO)
    xp_reward INTEGER DEFAULT 10, -- Награда за просмотр
    order_index INTEGER NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS lessons CASCADE;

DROP TABLE IF EXISTS modules CASCADE;

DROP TABLE IF EXISTS courses CASCADE;
-- +goose StatementEnd