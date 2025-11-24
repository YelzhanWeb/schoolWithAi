-- +goose Up
-- +goose StatementBegin
CREATE TYPE progress_status AS ENUM ('not_started', 'in_progress', 'completed');

CREATE TABLE lesson_progress (
    user_id TEXT REFERENCES users (id) ON DELETE CASCADE,
    lesson_id TEXT REFERENCES lessons (id) ON DELETE CASCADE,
    status progress_status DEFAULT 'not_started',
    is_completed BOOLEAN DEFAULT FALSE,
    last_accessed_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, lesson_id)
);

CREATE TABLE course_progress (
    user_id TEXT REFERENCES users (id) ON DELETE CASCADE,
    course_id TEXT REFERENCES courses (id) ON DELETE CASCADE,
    completed_lessons_count INTEGER DEFAULT 0,
    total_lessons_count INTEGER DEFAULT 0,
    progress_percentage INTEGER DEFAULT 0, -- 0-100%
    is_completed BOOLEAN DEFAULT FALSE,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, course_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS course_progress CASCADE;

DROP TABLE IF EXISTS lesson_progress CASCADE;

DROP TYPE IF EXISTS progress_status;
-- +goose StatementEnd