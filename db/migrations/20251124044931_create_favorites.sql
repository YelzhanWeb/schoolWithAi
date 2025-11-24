-- +goose Up
-- +goose StatementBegin
CREATE TABLE course_favorites (
    user_id TEXT REFERENCES users (id) ON DELETE CASCADE,
    course_id TEXT REFERENCES courses (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, course_id)
);

CREATE INDEX idx_favorites_course ON course_favorites (course_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS course_favorites;
-- +goose StatementEnd