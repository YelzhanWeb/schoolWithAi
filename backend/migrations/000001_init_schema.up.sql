
OqysAI
/

Scrum ÑÐºÐ¾Ð»Ð° (4).docx
310 lines

docx

Вот наш как бы описание проекта. Можешь по дольше подумать и сказать как это реализовать и потом помочь в реализации. Я вообще не создавал по микросервисной архитектуре приложение можем по монолитке сделать? потом JS я фреймворки не знаю. база данных Postgresql будет.

для API и Python Для рекомендательной системе

Отлично! Давай я проанализирую проект и предложу реалистичную архитектуру с учетом твоих технологий.

🎯 Анализ проекта
Это образовательная платформа с персонализацией через ИИ. Ключевые компоненты:

Персонализация: индивидуальные траектории обучения
ИИ-рекомендации: анализ успеваемости и подбор материалов
Геймификация: игровые элементы для мотивации
Образовательный контент: упражнения на критическое мышление
💡 Рекомендуемая архитектура
Да, начнем с монолита! Это правильное решение для MVP. Микросервисы можно добавить потом, если понадобится масштабирование.

Предлагаю гибридный монолит:
┌─────────────────────────────────────────┐
│         Frontend (Vanilla JS)           │
│    HTML + CSS + чистый JavaScript       │
└──────────────┬──────────────────────────┘
               │ REST API
┌──────────────▼──────────────────────────┐
│

Backend (монолит)            │
│  - REST API (Gin/Echo)                  │
│  - Auth & Sessions                      │
│  - Business Logic                       │
│  - Static Files                         │
└──────────────┬──────────────────────────┘
               │
    ┌──────────┼──────────┐
    ▼          ▼          ▼
┌────────┐ ┌─────────┐ ┌──────────────┐
│ Postgres│ │ Python  │ │ File Storage │
│   DB    │ │ ML/AI   │ │   (local)    │
└─────────┘ └─────────┘ └──────────────┘
Почему такая архитектура?
1.

Backend (монолит)

✅ Простая разработка и деплой
✅ Легко отлаживать
✅ Быстрая работа
✅ Один репозиторий
2. Python как отдельный сервис (не микросервис, а модуль)

Интегрируем через HTTP API или напрямую
Python будет обрабатывать ML-задачи
3. Vanilla JS (без фреймворков)

Не нужно учить React/Vue
Полный контроль
Быстрая загрузка
🏗️ Структура проекта
education-platform/
├── backend/              # Go приложение
│   ├── cmd/
│   │   └── main.
# Точка входа
│   ├── internal/
│   │   ├── api/         # HTTP handlers
│   │   ├── models/      # Структуры данных
│   │   ├── db/          # Работа с БД
│   │   ├── auth/        # Аутентификация
│   │   └── services/    # Бизнес-логика
│   ├── migrations/      # SQL миграции
│   └──

.mod
│
├── ml-service/          # Python ML сервис
│   ├── app.py          # Flask/FastAPI
│   ├── models/         # ML модели
│   ├── recommender.py  # Рекомендательная система
│   └── requirements.txt
│
├── frontend/           # Vanilla JS
│   ├── index.html
│   ├── css/
│   │   └── style.css
│   ├── js/
│   │   ├── app.js
│   │   ├── api.js      # API клиент
│   │   └── components/ # UI компоненты
│   └── assets/
│
└── docker-compose.yml  # Для локальной разработки
📊 Структура базы данных (основные таблицы)
sql
-- Пользователи
users (id, email, password_hash, role, created_at)

-- Профили учеников
student_profiles ( id, user_id, grade, age_group, interests )

-- Образовательные ресурсы
resources ( id, title, description, type, difficulty, content_url )

-- Прогресс обучения
progress ( id, student_id, resource_id, status, score, time_spent )

-- Рекомендации
recommendations ( id, student_id, resource_id, reason, created_at )

-- Достижения (геймификация)
achievements (id, student_id, badge_type, earned_at)
🚀 План реализации (по неделям)
Недели 1-2: Фундамент
Настроить

проект (Gin framework) Создать базовую структуру БД Реализовать простую аутентификацию Сделать регистрацию / вход Недели 3 -4:Базовый функционал CRUD для образовательных ресурсов Профили учеников Простой frontend на Vanilla JS Недели 5 -6:Python ML сервис FastAPI для ML API Базовая рекомендательная система Интеграция с

backend
Недели 7-8: Геймификация
Система достижений
Прогресс-бары
Лидерборды
🔧 Стек технологий (финальный)
Backend:

Go
1.21+
Gin Web Framework
GORM (ORM для PostgreSQL)
JWT для авторизации
ML Service:

Python 3.10+
FastAPI
scikit-learn / TensorFlow
pandas, numpy
Frontend:

Vanilla JavaScript (ES6+)
HTML5 + CSS3
Fetch API для запросов
Database:

PostgreSQL 15+
DevOps:

Docker & Docker Compose
Git
❓ Следующие шаги
Я готов помочь с:

Создать starter code - базовую структуру проекта
Написать Docker Compose для запуска всего стека
Создать первые API endpoints
Показать примеры Vanilla JS для работы с API
Что делаем первым? Предлагаю начать с базовой структуры

backend + настройка PostgreSQL. Согласен?

Давай сперва определимся с базой данных. Я не понял по проекту нужен быть разные роли да. И оброзавательный контент это какие и откуда брать самим создавать или это если будет роль учителя создают типа курс тесты и повествование тогда база должно быть большым нет.

Отличные вопросы! Давай разберемся по порядку.

🎭 Роли в системе
Судя по проекту, нужны 3 основные роли:

1. STUDENT (ученик) - основная целевая аудитория
   - Проходит обучение
   - Получает рекомендации от ИИ
   - Зарабатывает достижения

2. TEACHER (учитель) - создает и управляет контентом
   - Создает упражнения и материалы
   - Видит прогресс учеников
   - Дает обратную связь

3. ADMIN (администратор)
   - Управляет пользователями
   - Настройка системы
   - Аналитика
📚 Образовательный контент - откуда брать?
Есть 2 подхода:

Вариант 1: Учителя создают контент (рекомендую для старта)
Teacher создает:
├── Курсы/Модули
├── Упражнения на критическое мышление
├── Тесты/Квизы
└── Текстовые материалы
Плюсы:

✅ Не нужно заранее наполнять БД
✅ Контент адаптирован под конкретную школу/группу
✅ Учителя контролируют качество
✅ Проще начать разработку
Минусы:

❌ Нужен удобный редактор контента
❌ Сначала будет пусто
Вариант 2: Предзаполненная база + возможность расширения
Стартовый контент от разработчиков:
├── 50-100 базовых упражнений
├── 10-15 шаблонных курсов
└── Примеры по разным предметам

+ Учителя могут добавлять свое
💾 Размер базы данных
Давай посчитаем для среднего проекта:

Если учителя создают контент:
Контент:
- 1 курс = ~5-10 KB текста
- 1 упражнение = ~2-5 KB
- Изображения/файлы = хранятся отдельно (file storage)

Пример для школы на 500 учеников:
├── 500 учеников
├── 50 учителей
├── 200 курсов/модулей
├── 1000 упражнений
├── 10000 записей прогресса
└── ~50-100 MB БД + файлы

Это ОЧЕНЬ мало! PostgreSQL легко справится.
База НЕ будет большой, потому что:
Текстовый контент занимает мало места
Медиафайлы (видео, картинки) хранятся отдельно, не в БД
Даже 10,000 учеников = ~1-2 GB БД максимум
🗄️ Предлагаемая структура БД
sql
-- ПОЛЬЗОВАТЕЛИ И РОЛИ
CREATE TYPE user_role AS ENUM ('student', 'teacher', 'admin');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL,
    full_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ПРОФИЛИ УЧЕНИКОВ
CREATE TABLE student_profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
    grade INTEGER, -- класс обучения
    age_group VARCHAR(50), -- младшие/средние/старшие
    interests TEXT [], -- массив интересов
    learning_style VARCHAR(50), -- визуал/аудиал/кинестетик
    created_at TIMESTAMP DEFAULT NOW()
);

-- КУРСЫ (создает учитель)
CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_by INTEGER REFERENCES users (id), -- ID учителя
    difficulty_level INTEGER, -- 1-5
    age_group VARCHAR(50),
    subject VARCHAR(100), -- математика, литература и т.д.
    is_published BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- МОДУЛИ внутри курса
CREATE TABLE modules (
    id SERIAL PRIMARY KEY,
    course_id INTEGER REFERENCES courses (id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    order_index INTEGER, -- порядок в курсе
    created_at TIMESTAMP DEFAULT NOW()
);

-- ОБРАЗОВАТЕЛЬНЫЕ РЕСУРСЫ (упражнения, тесты)
CREATE TYPE resource_type AS ENUM (
    'exercise',      -- упражнение
    'quiz',          -- тест
    'reading',       -- текст для чтения
    'video',         -- видео
    'interactive'    -- интерактив
);

CREATE TABLE resources (
    id SERIAL PRIMARY KEY,
    module_id INTEGER REFERENCES modules (id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT, -- основной контент (JSON для гибкости)
    resource_type resource_type NOT NULL,
    difficulty INTEGER, -- 1-5
    estimated_time INTEGER, -- минуты
    file_url VARCHAR(500), -- если есть файл
    created_at TIMESTAMP DEFAULT NOW()
);

-- ВОПРОСЫ для тестов/упражнений
CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    resource_id INTEGER REFERENCES resources (id) ON DELETE CASCADE,
    question_text TEXT NOT NULL,
    question_type VARCHAR(50), -- multiple_choice, open_ended, etc.
    correct_answer TEXT, -- для автопроверки
    options JSONB, -- варианты ответов
    points INTEGER DEFAULT 1,
    order_index INTEGER
);

-- ПРОГРЕСС УЧЕНИКА
CREATE TABLE student_progress (
    id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES users (id),
    resource_id INTEGER REFERENCES resources (id),
    status VARCHAR(50), -- not_started, in_progress, completed
    score INTEGER, -- баллы
    time_spent INTEGER, -- минуты
    attempts INTEGER DEFAULT 0,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ОТВЕТЫ УЧЕНИКОВ
CREATE TABLE student_answers (
    id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES users (id),
    question_id INTEGER REFERENCES questions (id),
    answer TEXT,
    is_correct BOOLEAN,
    answered_at TIMESTAMP DEFAULT NOW()
);

-- РЕКОМЕНДАЦИИ ОТ ИИ
CREATE TABLE recommendations (
    id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES users (id),
    resource_id INTEGER REFERENCES resources (id),
    reason TEXT, -- почему рекомендовано
    priority INTEGER, -- приоритет рекомендации
    is_viewed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ГЕЙМИФИКАЦИЯ: Достижения
CREATE TABLE achievements (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    badge_icon VARCHAR(255), -- URL иконки
    criteria JSONB, -- условия получения
    points INTEGER DEFAULT 0
);

-- Достижения учеников
CREATE TABLE student_achievements (
    id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES users (id),
    achievement_id INTEGER REFERENCES achievements (id),
    earned_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (student_id, achievement_id)
);

-- ОЧКИ/БАЛЛЫ
CREATE TABLE student_points (
    id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES users (id) UNIQUE,
    total_points INTEGER DEFAULT 0,
    level INTEGER DEFAULT 1,
    updated_at TIMESTAMP DEFAULT NOW()
);


🎯 Мои рекомендации
Для MVP (первая версия):

Роли: Student + Teacher + Admin
Контент: Учителя создают через веб-интерфейс
Стартовый набор: 10-20 примеров упражнений для демонстрации
Хранение файлов: Локально на сервере (позже можно S3)
Этапы наполнения:

Этап 1: Вы создаете 20-30 примеров упражнений
         ↓
Этап 2: Запускаете для 1-2 учителей (бета)
         ↓
Этап 3: Учителя добавляют свой контент
         ↓
Этап 4: База растет естественным образом
❓ Твое мнение?
Роли: Согласен с Student/Teacher/Admin?
Контент: Учителя создают сами или нужна предзаполненная база?
Хочешь, я создам SQL миграции для этой структуры?
Что думаешь? Может что-то упростить или наоборот добавить?

Согласен
Давай сразу предзаполненная дальше потом учителя же могут же добавлять
Давай но база точно будет и останется таким есть вопросики. Потому что рекомендательной системе что нужно для рекомендации это нужно понять. Например мы уже реализовавали но там мы как делали типа ученик ставит оценки и сравнивается с другим учеником вкуси коллабративная фильтрация короче + вкус самого ученика тоже учитовалось. И тогда гибрид получается. Или нам не надо так делать что посоветуешь. Если мое идея нравиться то не надо ли бд менять ?
Отличный вопрос! Твоя идея с гибридной рекомендательной системой очень правильная для образовательной платформы. Давай разберем.

🤖 Рекомендательная система для обучения
Твой подход (гибридный) - ИДЕАЛЬНЫЙ для этого проекта!
Гибридная система =
├── Collaborative Filtering (коллаборативная)
│   └── "Ученики с похожим прогрессом любят эти курсы"
│
├── Content-Based (контентная)
│   └── "Тебе нравится математика → вот еще математика"
│
└── Knowledge-Based (на основе знаний)
    └── "Ты в 7 классе, слабо по алгебре → вот материал"
Почему гибрид лучше для образования?
Проблема чистой коллаборативной фильтрации:

❌ "Холодный старт" - новый ученик, нет данных
❌ Может рекомендовать слишком сложное/легкое
❌ Не учитывает педагогическую логику
Преимущества гибрида:

✅ Работает даже для новых учеников
✅ Учитывает уровень сложности
✅ Персонализация + педагогическая правильность
✅ Адаптируется со временем
📊 Что нужно добавить в БД для рекомендаций
1. Таблица для оценок/рейтингов
sql
-- ОЦЕНКИ РЕСУРСОВ учениками
CREATE TABLE resource_ratings (
    id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES users(id),
    resource_id INTEGER REFERENCES resources(id),
    rating INTEGER CHECK (rating >= 1 AND rating <= 5), -- 1-5 звезд
    review TEXT, -- опциональный отзыв
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(student_id, resource_id) -- один ученик = одна оценка
);

-- Индекс для быстрого поиска похожих учеников
CREATE INDEX idx_ratings_student ON resource_ratings (student_id);

CREATE INDEX idx_ratings_resource ON resource_ratings (resource_id);

2. Улучшенная таблица интересов sql
-- Расширяем student_profiles
ALTER TABLE student_profiles
ADD COLUMN IF NOT EXISTS preferences JSONB;

-- Пример структуры preferences:
-- {
--   "subjects": ["math", "physics", "literature"],
--   "difficulty_preference": "medium",
--   "learning_pace": "fast",
--   "resource_types": ["video", "interactive"]
-- }
3. Таблица взаимодействий (для ML) sql
-- ИСТОРИЯ ВЗАИМОДЕЙСТВИЙ для обучения модели
CREATE TABLE student_interactions (
    id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES users (id),
    resource_id INTEGER REFERENCES resources (id),
    interaction_type VARCHAR(50), -- viewed, started, completed, skipped, bookmarked
    duration INTEGER, -- секунды
    timestamp TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_interactions_student ON student_interactions (student_id);

CREATE INDEX idx_interactions_resource ON student_interactions (resource_id);

4. Таблица тегов / навыков sql
-- ТЕГИ для контента (для content-based filtering)
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    category VARCHAR(50) -- subject, skill, topic
);

-- Связь ресурсов с тегами
CREATE TABLE resource_tags (
    resource_id INTEGER REFERENCES resources (id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags (id) ON DELETE CASCADE,
    PRIMARY KEY (resource_id, tag_id)
);

-- НАВЫКИ ученика (отслеживание прогресса по скиллам)
CREATE TABLE student_skills (
    id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES users (id),
    skill_name VARCHAR(100), -- "algebra", "critical_thinking", "reading_comprehension"
    proficiency_level FLOAT, -- 0.0 - 1.0
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (student_id, skill_name)
);

5. Таблица для хранения векторов (опционально, для advanced ML) sql
-- Для более сложных ML моделей (embeddings)
CREATE TABLE resource_embeddings (
    resource_id INTEGER PRIMARY KEY REFERENCES resources (id),
    embedding FLOAT[], -- вектор представления ресурса
    model_version VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE student_embeddings (
    student_id INTEGER PRIMARY KEY REFERENCES users (id),
    embedding FLOAT[], -- вектор предпочтений ученика
    model_version VARCHAR(50),
    updated_at TIMESTAMP DEFAULT NOW()
);

🎯 Логика работы гибридной системы
python
# Псевдокод для Python ML сервиса

def get_recommendations(student_id, limit=10):

# 1. COLLABORATIVE FILTERING (40% веса)
similar_students = find_similar_students(student_id)
    collab_recommendations = get_resources_liked_by_similar(similar_students)

# 2. CONTENT-BASED (30% веса)
student_preferences = get_student_preferences(student_id)
    content_recommendations = find_similar_resources(student_preferences)

# 3. KNOWLEDGE-BASED (30% веса)
student_level = get_student_skill_level(student_id)
    student_gaps = identify_knowledge_gaps(student_id)
    knowledge_recommendations = get_resources_for_gaps(student_gaps)

# 4. ГИБРИДНЫЙ SCORING
recommendations = combine_scores({
        'collaborative': (collab_recommendations, 0.4),
        'content': (content_recommendations, 0.3),
        'knowledge': (knowledge_recommendations, 0.3)
    })

# 5. ФИЛЬТРЫ
recommendations = filter_already_completed(student_id, recommendations)
    recommendations = filter_by_difficulty(student_level, recommendations)
    
    return recommendations[:limit]
🗄️ Обновленная структура БД (финальная)
Database Schema - Education Platform
Code 
📈 Что мы добавили для ML:
1. Collaborative Filtering
resource_ratings - явные оценки учеников
student_interactions - неявные сигналы (просмотры, время)
2. Content-Based Filtering
tags + resource_tags - характеристики контента
student_profiles.preferences - предпочтения ученика
3. Knowledge-Based
student_skills - уровень владения навыками
difficulty в ресурсах - адаптация сложности
4. Hybrid Scoring
recommendations - хранит результаты всех алгоритмов
score - итоговый взвешенный скор
🎯 Следующие шаги
Создать миграции для

(я могу это сделать) Seed data - предзаполнить базу примерами (50 -100 ресурсов) Python ML сервис - реализовать гибридные рекомендации Что делаем дальше ? Создать

миграции?
Подготовить seed data (примеры курсов)?
Начать с Python ML API?

давай постепенно сперва 1 потом 2 потом 3

Отлично! Начинаем с миграций.

🚀 Шаг 1: Создание

миграций Для миграций будем использовать библиотеку golang - migrate.Это стандарт в

-сообществе. Структура проекта backend/ ├── cmd/ │ └── main.

├── migrations/
│   ├── 000001_init_schema.up.sql
│   └── 000001_init_schema.down.sql
├── internal/
│   ├── config/
│   ├── db/
│   └── models/
├──

.mod └──

.sum
Давай создам файлы миграций
000001_init_schema.up.sql
Code 
Теперь создадим файл для отката миграции (down):

000001_init_schema.down.sql
Code 
Теперь создадим

код для работы с миграциями и базой данных: database.

- Database Connection & Migrations
Code 
Теперь создадим основной main.

файл: main.

- Application Entry Point Code Теперь создадим файлы для настройки проекта:

go
.mod -

Dependencies Code Теперь создадим.env файл для конфигурации и docker - compose для удобной разработки:docker - compose.yml - Development Environment Code Теперь создадим.env файл и README с инструкциями:README.md - Setup Instructions Document

-- ========================================
-- MIGRATION: Initial Schema
-- Version: 000001
-- Description: Create all tables for education platform with ML recommendation system
-- ========================================

-- ENUMS
CREATE TYPE user_role AS ENUM ('student', 'teacher', 'admin');

CREATE TYPE resource_type AS ENUM ('exercise', 'quiz', 'reading', 'video', 'interactive');

CREATE TYPE progress_status AS ENUM ('not_started', 'in_progress', 'completed');

CREATE TYPE interaction_type AS ENUM ('viewed', 'started', 'completed', 'skipped', 'bookmarked', 'rated');

-- ========================================
-- USERS & AUTHENTICATION
-- ========================================

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users (email);

CREATE INDEX idx_users_role ON users (role);

-- ========================================
-- STUDENT PROFILES
-- ========================================

CREATE TABLE student_profiles (
    id SERIAL PRIMARY KEY,
    user_id INTEGER UNIQUE REFERENCES users (id) ON DELETE CASCADE,
    grade INTEGER CHECK (
        grade >= 1
        AND grade <= 11
    ),
    age_group VARCHAR(50) CHECK (
        age_group IN ('junior', 'middle', 'senior')
    ),
    interests TEXT [],
    learning_style VARCHAR(50) CHECK (
        learning_style IN (
            'visual',
            'auditory',
            'kinesthetic',
            'mixed'
        )
    ),
    preferences JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_student_profiles_user ON student_profiles (user_id);

CREATE INDEX idx_student_profiles_grade ON student_profiles (grade);

-- ========================================
-- CONTENT STRUCTURE
-- ========================================

CREATE TABLE courses (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_by INTEGER REFERENCES users (id) ON DELETE SET NULL,
    difficulty_level INTEGER CHECK (
        difficulty_level >= 1
        AND difficulty_level <= 5
    ),
    age_group VARCHAR(50),
    subject VARCHAR(100),
    is_published BOOLEAN DEFAULT FALSE,
    thumbnail_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_courses_subject ON courses (subject);

CREATE INDEX idx_courses_published ON courses (is_published);

CREATE INDEX idx_courses_difficulty ON courses (difficulty_level);

CREATE TABLE modules (
    id SERIAL PRIMARY KEY,
    course_id INTEGER REFERENCES courses (id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    order_index INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_modules_course ON modules (course_id);

CREATE INDEX idx_modules_order ON modules (course_id, order_index);

CREATE TABLE resources (
    id SERIAL PRIMARY KEY,
    module_id INTEGER REFERENCES modules (id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    resource_type resource_type NOT NULL,
    difficulty INTEGER CHECK (
        difficulty >= 1
        AND difficulty <= 5
    ),
    estimated_time INTEGER,
    file_url VARCHAR(500),
    thumbnail_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_resources_module ON resources (module_id);

CREATE INDEX idx_resources_type ON resources (resource_type);

CREATE INDEX idx_resources_difficulty ON resources (difficulty);

-- ========================================
-- QUESTIONS & ANSWERS
-- ========================================

CREATE TABLE questions (
    id SERIAL PRIMARY KEY,
    resource_id INTEGER REFERENCES resources (id) ON DELETE CASCADE,
    question_text TEXT NOT NULL,
    question_type VARCHAR(50) CHECK (
        question_type IN (
            'multiple_choice',
            'true_false',
            'open_ended',
            'matching'
        )
    ),
    correct_answer TEXT,
    options JSONB,
    explanation TEXT,
    points INTEGER DEFAULT 1,
    order_index INTEGER
);

CREATE INDEX idx_questions_resource ON questions (resource_id);

CREATE TABLE student_answers (
    id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
    question_id INTEGER REFERENCES questions (id) ON DELETE CASCADE,
    answer TEXT,
    is_correct BOOLEAN,
    answered_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_student_answers_student ON student_answers (student_id);

CREATE INDEX idx_student_answers_question ON student_answers (question_id);

-- ========================================
-- PROGRESS TRACKING
-- ========================================

CREATE TABLE student_progress (
    id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
    resource_id INTEGER REFERENCES resources (id) ON DELETE CASCADE,
    status progress_status DEFAULT 'not_started',
    score INTEGER,
    time_spent INTEGER DEFAULT 0,
    attempts INTEGER DEFAULT 0,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (student_id, resource_id)
);

CREATE INDEX idx_progress_student ON student_progress (student_id);

CREATE INDEX idx_progress_resource ON student_progress (resource_id);

CREATE INDEX idx_progress_status ON student_progress (status);

CREATE INDEX idx_progress_completed ON student_progress (completed_at);

-- ========================================
-- RECOMMENDATION SYSTEM
-- ========================================

-- Ratings for collaborative filtering
CREATE TABLE resource_ratings (
    id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
    resource_id INTEGER REFERENCES resources (id) ON DELETE CASCADE,
    rating INTEGER CHECK (
        rating >= 1
        AND rating <= 5
    ),
    review TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (student_id, resource_id)
);

CREATE INDEX idx_ratings_student ON resource_ratings (student_id);

CREATE INDEX idx_ratings_resource ON resource_ratings (resource_id);

CREATE INDEX idx_ratings_rating ON resource_ratings (rating);

-- Interaction history for ML
CREATE TABLE student_interactions (
    id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
    resource_id INTEGER REFERENCES resources (id) ON DELETE CASCADE,
    interaction_type interaction_type NOT NULL,
    duration INTEGER DEFAULT 0,
    timestamp TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_interactions_student ON student_interactions (student_id);

CREATE INDEX idx_interactions_resource ON student_interactions (resource_id);

CREATE INDEX idx_interactions_type ON student_interactions (interaction_type);

CREATE INDEX idx_interactions_timestamp ON student_interactions (timestamp);

-- Tags for content-based filtering
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    category VARCHAR(50) CHECK (
        category IN (
            'subject',
            'skill',
            'topic',
            'difficulty'
        )
    )
);

CREATE INDEX idx_tags_category ON tags (category);

CREATE TABLE resource_tags (
    resource_id INTEGER REFERENCES resources (id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags (id) ON DELETE CASCADE,
    weight FLOAT DEFAULT 1.0 CHECK (
        weight >= 0
        AND weight <= 1
    ),
    PRIMARY KEY (resource_id, tag_id)
);

CREATE INDEX idx_resource_tags_resource ON resource_tags (resource_id);

CREATE INDEX idx_resource_tags_tag ON resource_tags (tag_id);

-- Student skills for knowledge-based filtering
CREATE TABLE student_skills (
    id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
    skill_name VARCHAR(100) NOT NULL,
    proficiency_level FLOAT CHECK (
        proficiency_level >= 0
        AND proficiency_level <= 1
    ),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (student_id, skill_name)
);

CREATE INDEX idx_student_skills_student ON student_skills (student_id);

CREATE INDEX idx_student_skills_name ON student_skills (skill_name);

-- Saved recommendations
CREATE TABLE recommendations (
    id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
    resource_id INTEGER REFERENCES resources (id) ON DELETE CASCADE,
    score FLOAT,
    reason TEXT,
    algorithm_type VARCHAR(50) CHECK (
        algorithm_type IN (
            'collaborative',
            'content',
            'knowledge',
            'hybrid'
        )
    ),
    is_viewed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_recommendations_student ON recommendations (student_id);

CREATE INDEX idx_recommendations_score ON recommendations (score DESC);

CREATE INDEX idx_recommendations_created ON recommendations (created_at DESC);

-- Embeddings for advanced ML (optional)
CREATE TABLE resource_embeddings (
    resource_id INTEGER PRIMARY KEY REFERENCES resources (id) ON DELETE CASCADE,
    embedding FLOAT[],
    model_version VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE student_embeddings (
    student_id INTEGER PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    embedding FLOAT[],
    model_version VARCHAR(50),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- ========================================
-- GAMIFICATION
-- ========================================

CREATE TABLE achievements (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    badge_icon VARCHAR(255),
    criteria JSONB,
    points INTEGER DEFAULT 0
);

CREATE TABLE student_achievements (
    id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
    achievement_id INTEGER REFERENCES achievements (id) ON DELETE CASCADE,
    earned_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (student_id, achievement_id)
);

CREATE INDEX idx_student_achievements_student ON student_achievements (student_id);

CREATE TABLE student_points (
    id SERIAL PRIMARY KEY,
    student_id INTEGER UNIQUE REFERENCES users (id) ON DELETE CASCADE,
    total_points INTEGER DEFAULT 0,
    level INTEGER DEFAULT 1,
    experience INTEGER DEFAULT 0,
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_points_total ON student_points (total_points DESC);

CREATE TABLE leaderboard (
    id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
    period VARCHAR(50) CHECK (
        period IN (
            'daily',
            'weekly',
            'monthly',
            'all_time'
        )
    ),
    rank INTEGER,
    points INTEGER,
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_leaderboard_period_rank ON leaderboard (period, rank);

-- ========================================
-- INITIAL SEED DATA
-- ========================================

-- Insert default tags
INSERT INTO
    tags (name, category)
VALUES ('mathematics', 'subject'),
    ('science', 'subject'),
    ('literature', 'subject'),
    ('history', 'subject'),
    ('physics', 'subject'),
    ('chemistry', 'subject'),
    ('biology', 'subject'),
    ('geography', 'subject'),
    ('algebra', 'skill'),
    ('geometry', 'skill'),
    ('critical_thinking', 'skill'),
    ('problem_solving', 'skill'),
    (
        'reading_comprehension',
        'skill'
    ),
    ('writing', 'skill'),
    ('logic', 'skill'),
    ('analysis', 'skill'),
    ('beginner', 'difficulty'),
    ('intermediate', 'difficulty'),
    ('advanced', 'difficulty'),
    ('expert', 'difficulty');

-- Insert default achievements
INSERT INTO
    achievements (
        title,
        description,
        badge_icon,
        criteria,
        points
    )
VALUES (
        'First Steps',
        'Complete your first lesson',
        'badge_first.png',
        '{"completed_resources": 1}'::jsonb,
        10
    ),
    (
        'Fast Learner',
        'Complete 5 lessons in one day',
        'badge_fast.png',
        '{"daily_completions": 5}'::jsonb,
        50
    ),
    (
        'Perfect Score',
        'Get 100% on a quiz',
        'badge_perfect.png',
        '{"perfect_score": true}'::jsonb,
        30
    ),
    (
        'Week Warrior',
        'Study every day for a week',
        'badge_streak.png',
        '{"streak_days": 7}'::jsonb,
        100
    ),
    (
        'Knowledge Seeker',
        'Complete 10 different courses',
        'badge_courses.png',
        '{"completed_courses": 10}'::jsonb,
        200
    ),
    (
        'Math Master',
        'Complete all math courses',
        'badge_math.png',
        '{"subject": "mathematics", "completion": 100}'::jsonb,
        150
    ),
    (
        'Early Bird',
        'Study before 8 AM',
        'badge_early.png',
        '{"early_study": true}'::jsonb,
        25
    ),
    (
        'Night Owl',
        'Study after 10 PM',
        'badge_night.png',
        '{"late_study": true}'::jsonb,
        25
    ),
    (
        'Social Learner',
        'Help 5 other students',
        'badge_social.png',
        '{"helped_students": 5}'::jsonb,
        75
    ),
    (
        'Perfectionist',
        'Get 100% on 5 quizzes',
        'badge_perfectionist.png',
        '{"perfect_scores": 5}'::jsonb,
        100
    );