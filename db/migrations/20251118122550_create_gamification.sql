-- +goose Up
-- +goose StatementBegin

CREATE TABLE leagues (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(50) NOT NULL,
    order_index INTEGER NOT NULL,
    icon_url TEXT
);

INSERT INTO
    leagues (
        slug,
        name,
        order_index,
        icon_url
    )
VALUES (
        'bronze',
        'Бронзовая лига',
        1,
        "http://localhost:9000/school-assets/leagues/bronze_league.png"
    ),
    (
        'silver',
        'Серебряная лига',
        2,
        "http://localhost:9000/school-assets/leagues/silver_league.png"
    ),
    (
        'gold',
        'Золотая лига',
        3,
        "http://localhost:9000/school-assets/leagues/gold_league.png"
    ),
    (
        'sapphire',
        'Сапфировая лига',
        4,
        "http://localhost:9000/school-assets/leagues/sapphire_league.png"
    ),
    (
        'ruby',
        'Рубиновая лига',
        5,
        "http://localhost:9000/school-assets/leagues/ruby_league.png"
    ),
    (
        'emerald',
        'Изумрудная лига',
        6,
        "http://localhost:9000/school-assets/leagues/emerald_league.png"
    ),
    (
        'amethyst',
        'Аметистовая лига',
        7,
        "http://localhost:9000/school-assets/leagues/amethyst_league.png"
    ),
    (
        'diamond',
        'Алмазная лига',
        8,
        "http://localhost:9000/school-assets/leagues/diamond_league.png"
    );

CREATE TABLE achievements (
    id TEXT PRIMARY KEY,
    slug VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    icon_url TEXT,
    xp_reward INTEGER DEFAULT 0 NOT NULL
);

CREATE TABLE user_achievements (
    user_id TEXT REFERENCES users (id) ON DELETE CASCADE,
    achievement_id TEXT REFERENCES achievements (id) ON DELETE CASCADE,
    earned_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (user_id, achievement_id)
);

CREATE TABLE leaderboard_history (
    id TEXT PRIMARY KEY,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    user_id TEXT REFERENCES users (id) ON DELETE CASCADE,
    league_id INTEGER REFERENCES leagues (id),
    rank INTEGER NOT NULL,
    total_xp BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_leaderboard_history_period ON leaderboard_history (period_start, period_end);

CREATE INDEX idx_leaderboard_history_user ON leaderboard_history (user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS leaderboard_history CASCADE;

DROP TABLE IF EXISTS user_achievements CASCADE;

DROP TABLE IF EXISTS achievements CASCADE;

DROP TABLE IF EXISTS leagues CASCADE;
-- +goose StatementEnd