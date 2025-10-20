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