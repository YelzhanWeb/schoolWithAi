-- ========================================
-- MIGRATION DOWN: Rollback Initial Schema
-- Version: 000001
-- ========================================

-- Drop tables in reverse order (respecting foreign keys)

-- Gamification
DROP TABLE IF EXISTS leaderboard CASCADE;

DROP TABLE IF EXISTS student_points CASCADE;

DROP TABLE IF EXISTS student_achievements CASCADE;

DROP TABLE IF EXISTS achievements CASCADE;

-- Recommendation System
DROP TABLE IF EXISTS student_embeddings CASCADE;

DROP TABLE IF EXISTS resource_embeddings CASCADE;

DROP TABLE IF EXISTS recommendations CASCADE;

DROP TABLE IF EXISTS student_skills CASCADE;

DROP TABLE IF EXISTS resource_tags CASCADE;

DROP TABLE IF EXISTS tags CASCADE;

DROP TABLE IF EXISTS student_interactions CASCADE;

DROP TABLE IF EXISTS resource_ratings CASCADE;

-- Progress & Answers
DROP TABLE IF EXISTS student_progress CASCADE;

DROP TABLE IF EXISTS student_answers CASCADE;

DROP TABLE IF EXISTS questions CASCADE;

-- Content
DROP TABLE IF EXISTS resources CASCADE;

DROP TABLE IF EXISTS modules CASCADE;

DROP TABLE IF EXISTS courses CASCADE;

-- Profiles
DROP TABLE IF EXISTS student_profiles CASCADE;

-- Users
DROP TABLE IF EXISTS users CASCADE;

-- Drop ENUMs
DROP TYPE IF EXISTS interaction_type;

DROP TYPE IF EXISTS progress_status;

DROP TYPE IF EXISTS resource_type;

DROP TYPE IF EXISTS user_role;