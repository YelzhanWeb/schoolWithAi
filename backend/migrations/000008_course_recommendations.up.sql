-- ========================================
-- MIGRATION: Course Recommendations Table
-- Таблица для рекомендаций КУРСОВ (не ресурсов)
-- ========================================

CREATE TABLE IF NOT EXISTS course_recommendations (
    id SERIAL PRIMARY KEY,
    student_id INTEGER REFERENCES users (id) ON DELETE CASCADE,
    course_id INTEGER REFERENCES courses (id) ON DELETE CASCADE,
    score FLOAT,
    reason TEXT,
    algorithm_type VARCHAR(50) CHECK (
        algorithm_type IN (
            'collaborative',
            'content_based',
            'knowledge_based',
            'hybrid'
        )
    ),
    is_viewed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (student_id, course_id)
);

CREATE INDEX idx_course_recommendations_student ON course_recommendations (student_id);

CREATE INDEX idx_course_recommendations_score ON course_recommendations (score DESC);

CREATE INDEX idx_course_recommendations_created ON course_recommendations (created_at DESC);

COMMENT ON TABLE course_recommendations IS 'Рекомендации курсов для студентов (не отдельных ресурсов)';