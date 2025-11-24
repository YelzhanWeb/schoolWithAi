-- +goose Up
-- +goose StatementBegin
CREATE TABLE subjects (
    id TEXT PRIMARY KEY,
    slug VARCHAR(50) UNIQUE NOT NULL,
    name_ru VARCHAR(100) NOT NULL,
    name_kz VARCHAR(100) NOT NULL
);

INSERT INTO
    subjects (id, slug, name_ru, name_kz)
VALUES (
        'math-12345',
        'math',
        'Математика',
        'Математика'
    ),
    (
        'kaz_lang-12345',
        'kaz_lang',
        'Казахский язык',
        'Қазақ тілі'
    ),
    (
        'history_kz-12345',
        'history_kz',
        'История Казахстана',
        'Қазақстан тарихы'
    ),
    (
        'physics-12345',
        'physics',
        'Физика',
        'Физика'
    );

CREATE TABLE student_profiles (
    id TEXT PRIMARY KEY,
    user_id TEXT UNIQUE REFERENCES users (id) ON DELETE CASCADE,
    grade INTEGER CHECK (
        grade >= 1
        AND grade <= 11
    ),
    xp BIGINT DEFAULT 0 NOT NULL,
    level INTEGER DEFAULT 1 NOT NULL,
    --Weekly Leaderboard
    current_league_id INTEGER REFERENCES leagues (id) DEFAULT 1,
    weekly_xp BIGINT DEFAULT 0 NOT NULL,
    --Activity
    current_streak INTEGER DEFAULT 0,
    max_streak INTEGER DEFAULT 0,
    last_activity_date DATE DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE student_interests (
    user_id TEXT REFERENCES users (id) ON DELETE CASCADE,
    subject_id TEXT REFERENCES subjects (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, subject_id)
);

CREATE INDEX idx_league_ranking ON student_profiles (
    current_league_id,
    weekly_xp DESC
);

CREATE INDEX idx_student_profiles_grade ON student_profiles (grade);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS student_interests;

DROP TABLE IF EXISTS student_profiles;

DROP TABLE IF EXISTS subjects;
-- +goose StatementEnd