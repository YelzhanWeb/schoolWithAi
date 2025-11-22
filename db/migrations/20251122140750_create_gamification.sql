-- +goose Up
-- +goose StatementBegin
CREATE TABLE achievements (
    id TEXT PRIMARY KEY,
    slug TEXT UNIQUE NOT NULL, -- 'fast_learner'
    name TEXT NOT NULL,
    description TEXT,
    icon_url TEXT
);

CREATE TABLE user_achievements (
    user_id TEXT REFERENCES users (id) ON DELETE CASCADE,
    achievement_id TEXT REFERENCES achievements (id) ON DELETE CASCADE,
    earned_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, achievement_id)
);

CREATE TABLE leaderboard_history (
    id TEXT PRIMARY KEY,
    period_start DATE NOT NULL, -- Начало недели
    period_end DATE NOT NULL, -- Конец недели
    user_id TEXT REFERENCES users (id) ON DELETE CASCADE,
    rank INTEGER NOT NULL, -- Место (1, 2, 3...)
    total_xp BIGINT NOT NULL, -- Сколько набрал за ту неделю
    created_at TIMESTAMP DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS leaderboard_history CASCADE;

DROP TABLE IF EXISTS user_achievements CASCADE;

DROP TABLE IF EXISTS achievements CASCADE;

-- +goose StatementEnd