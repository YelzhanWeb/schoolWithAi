-- ========================================
-- SEED: Student Skills (Massive Data)
-- Критически важно для Knowledge-Based фильтрации!
-- ========================================

-- Полная очистка таблицы перед вставкой
TRUNCATE TABLE student_skills RESTART IDENTITY CASCADE;

-- Навыки должны совпадать с тегами в таблице tags!
-- Смотрим теги: algebra, geometry, critical_thinking, problem_solving,
--               reading_comprehension, writing, logic, analysis

-- ========================================
-- Ученик 5 (Иван)
-- ========================================
INSERT INTO
    student_skills (
        student_id,
        skill_name,
        proficiency_level
    )
VALUES (5, 'geometry', 0.55),
    (5, 'problem_solving', 0.78),
    (5, 'critical_thinking', 0.62),
    (5, 'logic', 0.70),
    (
        5,
        'reading_comprehension',
        0.35
    ),
    (5, 'writing', 0.30),
    (5, 'analysis', 0.48)
ON CONFLICT (student_id, skill_name) DO NOTHING;

-- ========================================
-- Ученик 6 (Алия)
-- ========================================
INSERT INTO
    student_skills (
        student_id,
        skill_name,
        proficiency_level
    )
VALUES (6, 'algebra', 0.25),
    (6, 'geometry', 0.30),
    (6, 'problem_solving', 0.40),
    (6, 'critical_thinking', 0.82),
    (6, 'logic', 0.65),
    (
        6,
        'reading_comprehension',
        0.92
    ),
    (6, 'writing', 0.88),
    (6, 'analysis', 0.85)
ON CONFLICT (student_id, skill_name) DO NOTHING;

-- ========================================
-- Ученик 7 (Дмитрий)
-- ========================================
INSERT INTO
    student_skills (
        student_id,
        skill_name,
        proficiency_level
    )
VALUES (7, 'algebra', 0.35),
    (7, 'geometry', 0.45),
    (7, 'problem_solving', 0.42),
    (7, 'critical_thinking', 0.38),
    (7, 'logic', 0.48),
    (
        7,
        'reading_comprehension',
        0.52
    ),
    (7, 'writing', 0.40),
    (7, 'analysis', 0.35)
ON CONFLICT (student_id, skill_name) DO NOTHING;

-- ========================================
-- Ученик 8 (Айгерим)
-- ========================================
INSERT INTO
    student_skills (
        student_id,
        skill_name,
        proficiency_level
    )
VALUES (8, 'algebra', 0.95),
    (8, 'geometry', 0.88),
    (8, 'problem_solving', 0.92),
    (8, 'critical_thinking', 0.85),
    (8, 'logic', 0.90),
    (
        8,
        'reading_comprehension',
        0.78
    ),
    (8, 'writing', 0.72),
    (8, 'analysis', 0.82)
ON CONFLICT (student_id, skill_name) DO NOTHING;

-- ========================================
-- Ученик 9 (Максим)
-- ========================================
INSERT INTO
    student_skills (
        student_id,
        skill_name,
        proficiency_level
    )
VALUES (9, 'algebra', 0.28),
    (9, 'geometry', 0.32),
    (9, 'problem_solving', 0.35),
    (9, 'critical_thinking', 0.30),
    (9, 'logic', 0.38),
    (
        9,
        'reading_comprehension',
        0.45
    ),
    (9, 'writing', 0.35),
    (9, 'analysis', 0.28)
ON CONFLICT (student_id, skill_name) DO NOTHING;

-- ========================================
-- Ученик 10 (Камила)
-- ========================================
INSERT INTO
    student_skills (
        student_id,
        skill_name,
        proficiency_level
    )
VALUES (10, 'algebra', 0.38),
    (10, 'geometry', 0.42),
    (10, 'problem_solving', 0.48),
    (10, 'critical_thinking', 0.88),
    (10, 'logic', 0.75),
    (
        10,
        'reading_comprehension',
        0.90
    ),
    (10, 'writing', 0.85),
    (10, 'analysis', 0.82)
ON CONFLICT (student_id, skill_name) DO NOTHING;

-- ========================================
-- Дополнительные предметные навыки
-- ========================================
INSERT INTO
    student_skills (
        student_id,
        skill_name,
        proficiency_level
    )
VALUES
    -- Иван
    (5, 'physics', 0.75),
    (5, 'chemistry', 0.65),
    (5, 'biology', 0.45),

-- Алия
(6, 'physics', 0.30),
(6, 'chemistry', 0.28),
(6, 'biology', 0.35),
(6, 'history', 0.85),
(6, 'literature', 0.92),

-- Дмитрий
(7, 'physics', 0.40),
(7, 'chemistry', 0.38),
(7, 'biology', 0.55),

-- Айгерим
(8, 'physics', 0.88),
(8, 'chemistry', 0.92),
(8, 'biology', 0.80),
(8, 'history', 0.75),
(8, 'literature', 0.70),

-- Максим
(9, 'physics', 0.25),
(9, 'chemistry', 0.22),
(9, 'biology', 0.35),
(9, 'history', 0.40),
(9, 'literature', 0.38),

-- Камила
(10, 'physics', 0.35),
(10, 'chemistry', 0.30),
(10, 'biology', 0.45),
(10, 'history', 0.88),
(10, 'literature', 0.90)
ON CONFLICT (student_id, skill_name) DO NOTHING;

COMMIT;