-- +goose Up
-- +goose StatementBegin
CREATE TABLE subjects (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(50) UNIQUE NOT NULL,
    name_ru VARCHAR(100) NOT NULL,
    name_kz VARCHAR(100) NOT NULL
);

INSERT INTO
    subjects (slug, name_ru, name_kz)
VALUES (
        'math',
        'Математика',
        'Математика'
    ),
    (
        'kaz_lang',
        'Казахский язык',
        'Қазақ тілі'
    ),
    (
        'history_kz',
        'История Казахстана',
        'Қазақстан тарихы'
    ),
    ('physics', 'Физика', 'Физика');

CREATE TABLE student_profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE REFERENCES users (id) ON DELETE CASCADE,
    grade INTEGER CHECK (
        grade >= 1
        AND grade <= 11
    ),
    xp BIGINT DEFAULT 0 NOT NULL,
    level INTEGER DEFAULT 1 NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE student_interests (
    user_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
    subject_id INTEGER REFERENCES subjects (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, subject_id)
);

CREATE INDEX idx_student_profiles_grade ON student_profiles (grade);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS student_interests;

DROP TABLE IF EXISTS student_profiles;

DROP TABLE IF EXISTS subjects;
-- +goose StatementEnd