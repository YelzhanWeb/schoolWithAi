-- USERS
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    class_grade SMALLINT CHECK (
        class_grade BETWEEN 1 AND 11
    ),
    created_at TIMESTAMP DEFAULT now()
);
-- COURSES (внешние и внутренние)
CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    external_url TEXT,
    -- ссылка на внешний курс
    min_grade SMALLINT,
    max_grade SMALLINT,
    difficulty SMALLINT DEFAULT 1,
    -- 1..5
    is_external BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT now()
);
-- QUESTIONS для мини-тестов
CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    course_id INT REFERENCES courses(id) ON DELETE CASCADE,
    q_type TEXT NOT NULL,
    -- single | multiple | tf | short
    text TEXT NOT NULL,
    options JSONB,
    -- ["A","B","C"] или NULL
    correct_answer TEXT NOT NULL
);
-- USER_ANSWERS (ответы учеников)
CREATE TABLE user_answers (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    question_id INT REFERENCES questions(id),
    answer TEXT,
    is_correct BOOLEAN,
    created_at TIMESTAMP DEFAULT now()
);
-- PROGRESS (прогресс по курсам)
CREATE TABLE user_progress (
    user_id INT REFERENCES users(id),
    course_id INT REFERENCES courses(id),
    status TEXT DEFAULT 'not_started',
    -- not_started | in_progress | completed
    progress_percent SMALLINT DEFAULT 0,
    score INT DEFAULT 0,
    last_updated TIMESTAMP DEFAULT now(),
    PRIMARY KEY (user_id, course_id)
);
-- SCORES (общие очки для рейтинга)
CREATE TABLE user_scores (
    user_id INT PRIMARY KEY REFERENCES users(id),
    total_points INT DEFAULT 0
);
-- RECOMMENDATIONS (храним кеш рекомендаций)
CREATE TABLE recommendations (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    course_id INT REFERENCES courses(id),
    score NUMERIC,
    reason TEXT,
    updated_at TIMESTAMP DEFAULT now()
);
-- Пользователи
INSERT INTO users (name, email, password_hash, class_grade)
VALUES ('Айдана', 'aidana@example.com', 'hash1', 3),
    -- младшие классы
    ('Марат', 'marat@example.com', 'hash2', 7),
    -- средние классы
    ('Жансая', 'zhansaya@example.com', 'hash3', 10);
-- старшие классы
-- Курсы
INSERT INTO courses (
        title,
        description,
        external_url,
        min_grade,
        max_grade,
        difficulty,
        is_external
    )
VALUES (
        'Основы математики',
        'Базовые арифметические операции для младших классов',
        'https://www.khanacademy.org/math/arithmetic',
        1,
        4,
        1,
        TRUE
    ),
    (
        'Логика и мышление',
        'Развитие критического мышления через логические задачи',
        'https://stepik.org/course/LogicalThinking',
        5,
        9,
        2,
        TRUE
    ),
    (
        'Основы программирования Python',
        'Начало работы с Python для старших школьников',
        'https://stepik.org/course/BeginPython',
        9,
        11,
        3,
        TRUE
    );
-- Вопросы к курсу "Основы математики"
INSERT INTO questions (course_id, q_type, text, options, correct_answer)
VALUES (
        1,
        'single',
        'Сколько будет 2 + 2?',
        '["3", "4", "5"]',
        '4'
    ),
    (
        1,
        'single',
        'Какое число больше?',
        '["5", "8", "6"]',
        '8'
    );
-- Вопросы к курсу "Логика и мышление"
INSERT INTO questions (course_id, q_type, text, options, correct_answer)
VALUES (
        2,
        'single',
        'У Маши есть 3 яблока, у Пети — 5. Сколько всего яблок?',
        '["7", "8", "9"]',
        '8'
    ),
    (
        2,
        'tf',
        'Если сегодня понедельник, то через 2 дня будет среда.',
        NULL,
        'true'
    );
-- Вопросы к курсу "Основы Python"
INSERT INTO questions (course_id, q_type, text, options, correct_answer)
VALUES (
        3,
        'single',
        'Какой тип данных у числа 5 в Python?',
        '["int", "float", "string"]',
        'int'
    ),
    (
        3,
        'single',
        'Что выведет print(2 * 3)?',
        '["5", "6", "23"]',
        '6'
    );