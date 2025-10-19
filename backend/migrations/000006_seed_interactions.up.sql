-- ========================================
-- SEED: Student Interactions (Massive Data)
-- История взаимодействий для аналитики
-- ========================================

-- Ученик 5 (Иван) - активный математик
INSERT INTO
    student_interactions (
        student_id,
        resource_id,
        interaction_type,
        duration,
        timestamp
    )
VALUES
    -- День 1: Начал алгебру
    (
        5,
        1,
        'viewed',
        45,
        NOW() - INTERVAL '7 days'
    ),
    (
        5,
        1,
        'started',
        600,
        NOW() - INTERVAL '7 days'
    ),
    (
        5,
        1,
        'completed',
        600,
        NOW() - INTERVAL '7 days'
    ),
    (
        5,
        2,
        'viewed',
        30,
        NOW() - INTERVAL '7 days'
    ),
    (
        5,
        2,
        'started',
        480,
        NOW() - INTERVAL '7 days'
    ),
    (
        5,
        2,
        'completed',
        480,
        NOW() - INTERVAL '7 days'
    ),

-- День 2: Продолжил алгебру
(
    5,
    3,
    'viewed',
    20,
    NOW() - INTERVAL '6 days'
),
(
    5,
    3,
    'started',
    900,
    NOW() - INTERVAL '6 days'
),
(
    5,
    3,
    'completed',
    900,
    NOW() - INTERVAL '6 days'
),
(
    5,
    4,
    'viewed',
    40,
    NOW() - INTERVAL '6 days'
),
(
    5,
    4,
    'started',
    720,
    NOW() - INTERVAL '6 days'
),
(
    5,
    4,
    'completed',
    720,
    NOW() - INTERVAL '6 days'
),

-- День 3: Видео и упражнения
(
    5,
    5,
    'viewed',
    30,
    NOW() - INTERVAL '5 days'
),
(
    5,
    5,
    'started',
    900,
    NOW() - INTERVAL '5 days'
),
(
    5,
    5,
    'completed',
    900,
    NOW() - INTERVAL '5 days'
),
(
    5,
    5,
    'rated',
    0,
    NOW() - INTERVAL '5 days'
),
(
    5,
    6,
    'viewed',
    25,
    NOW() - INTERVAL '5 days'
),
(
    5,
    6,
    'started',
    1200,
    NOW() - INTERVAL '5 days'
),
(
    5,
    6,
    'completed',
    1200,
    NOW() - INTERVAL '5 days'
),

-- День 4: Тесты и системы
(
    5,
    7,
    'viewed',
    15,
    NOW() - INTERVAL '4 days'
),
(
    5,
    7,
    'started',
    600,
    NOW() - INTERVAL '4 days'
),
(
    5,
    7,
    'completed',
    600,
    NOW() - INTERVAL '4 days'
),
(
    5,
    8,
    'viewed',
    35,
    NOW() - INTERVAL '4 days'
),
(
    5,
    8,
    'bookmarked',
    0,
    NOW() - INTERVAL '4 days'
),
(
    5,
    8,
    'started',
    900,
    NOW() - INTERVAL '4 days'
),
(
    5,
    8,
    'completed',
    900,
    NOW() - INTERVAL '4 days'
),

-- День 5: Физика
(
    5,
    20,
    'viewed',
    40,
    NOW() - INTERVAL '3 days'
),
(
    5,
    20,
    'started',
    900,
    NOW() - INTERVAL '3 days'
),
(
    5,
    20,
    'completed',
    900,
    NOW() - INTERVAL '3 days'
),
(
    5,
    21,
    'viewed',
    30,
    NOW() - INTERVAL '3 days'
),
(
    5,
    21,
    'started',
    720,
    NOW() - INTERVAL '3 days'
),
(
    5,
    21,
    'completed',
    720,
    NOW() - INTERVAL '3 days'
),

-- День 6-7: Химия
(
    5,
    42,
    'viewed',
    35,
    NOW() - INTERVAL '2 days'
),
(
    5,
    42,
    'started',
    900,
    NOW() - INTERVAL '2 days'
),
(
    5,
    42,
    'completed',
    900,
    NOW() - INTERVAL '2 days'
),
(
    5,
    43,
    'viewed',
    25,
    NOW() - INTERVAL '1 day'
),
(
    5,
    43,
    'started',
    720,
    NOW() - INTERVAL '1 day'
),
(
    5,
    43,
    'completed',
    720,
    NOW() - INTERVAL '1 day'
);

-- Ученик 6 (Алия) - любит литературу и историю
INSERT INTO
    student_interactions (
        student_id,
        resource_id,
        interaction_type,
        duration,
        timestamp
    )
VALUES
    -- Литература (последние 5 дней)
    (
        6,
        23,
        'viewed',
        50,
        NOW() - INTERVAL '5 days'
    ),
    (
        6,
        23,
        'started',
        1080,
        NOW() - INTERVAL '5 days'
    ),
    (
        6,
        23,
        'completed',
        1080,
        NOW() - INTERVAL '5 days'
    ),
    (
        6,
        23,
        'rated',
        0,
        NOW() - INTERVAL '5 days'
    ),
    (
        6,
        24,
        'viewed',
        45,
        NOW() - INTERVAL '4 days'
    ),
    (
        6,
        24,
        'bookmarked',
        0,
        NOW() - INTERVAL '4 days'
    ),
    (
        6,
        24,
        'started',
        900,
        NOW() - INTERVAL '4 days'
    ),
    (
        6,
        24,
        'completed',
        900,
        NOW() - INTERVAL '4 days'
    ),
    (
        6,
        26,
        'viewed',
        30,
        NOW() - INTERVAL '3 days'
    ),
    (
        6,
        26,
        'started',
        600,
        NOW() - INTERVAL '3 days'
    ),
    (
        6,
        26,
        'completed',
        600,
        NOW() - INTERVAL '3 days'
    ),
    (
        6,
        28,
        'viewed',
        60,
        NOW() - INTERVAL '2 days'
    ),
    (
        6,
        28,
        'started',
        1200,
        NOW() - INTERVAL '2 days'
    ),
    (
        6,
        28,
        'completed',
        1200,
        NOW() - INTERVAL '2 days'
    ),
    (
        6,
        29,
        'viewed',
        70,
        NOW() - INTERVAL '1 day'
    ),
    (
        6,
        29,
        'started',
        1800,
        NOW() - INTERVAL '1 day'
    ),
    (
        6,
        29,
        'completed',
        1800,
        NOW() - INTERVAL '1 day'
    ),

-- История
(
    6,
    32,
    'viewed',
    40,
    NOW() - INTERVAL '4 days'
),
(
    6,
    32,
    'started',
    600,
    NOW() - INTERVAL '4 days'
),
(
    6,
    32,
    'completed',
    600,
    NOW() - INTERVAL '4 days'
),
(
    6,
    37,
    'viewed',
    50,
    NOW() - INTERVAL '2 days'
),
(
    6,
    37,
    'bookmarked',
    0,
    NOW() - INTERVAL '2 days'
),
(
    6,
    37,
    'started',
    900,
    NOW() - INTERVAL '2 days'
),
(
    6,
    37,
    'completed',
    900,
    NOW() - INTERVAL '2 days'
);

-- Ученик 7 (Дмитрий) - начинающий, пробует разное
INSERT INTO
    student_interactions (
        student_id,
        resource_id,
        interaction_type,
        duration,
        timestamp
    )
VALUES
    -- Геометрия (3 дня назад)
    (
        7,
        11,
        'viewed',
        35,
        NOW() - INTERVAL '3 days'
    ),
    (
        7,
        11,
        'started',
        540,
        NOW() - INTERVAL '3 days'
    ),
    (
        7,
        11,
        'completed',
        540,
        NOW() - INTERVAL '3 days'
    ),
    (
        7,
        12,
        'viewed',
        40,
        NOW() - INTERVAL '3 days'
    ),
    (
        7,
        12,
        'started',
        600,
        NOW() - INTERVAL '3 days'
    ),
    (
        7,
        12,
        'completed',
        600,
        NOW() - INTERVAL '3 days'
    ),
    (
        7,
        13,
        'viewed',
        20,
        NOW() - INTERVAL '2 days'
    ),
    (
        7,
        13,
        'started',
        480,
        NOW() - INTERVAL '2 days'
    ),
    (
        7,
        13,
        'skipped',
        240,
        NOW() - INTERVAL '2 days'
    ),

-- Биология (вчера)
(
    7,
    47,
    'viewed',
    30,
    NOW() - INTERVAL '1 day'
),
(
    7,
    47,
    'started',
    600,
    NOW() - INTERVAL '1 day'
),
(
    7,
    47,
    'completed',
    600,
    NOW() - INTERVAL '1 day'
),
(
    7,
    48,
    'viewed',
    45,
    NOW() - INTERVAL '1 day'
),
(
    7,
    48,
    'started',
    900,
    NOW() - INTERVAL '1 day'
),
(
    7,
    48,
    'completed',
    900,
    NOW() - INTERVAL '1 day'
);

-- Ученик 8 (Айгерим) - очень активная
INSERT INTO
    student_interactions (
        student_id,
        resource_id,
        interaction_type,
        duration,
        timestamp
    )
VALUES
    -- Математика (7 дней назад)
    (
        8,
        1,
        'viewed',
        20,
        NOW() - INTERVAL '7 days'
    ),
    (
        8,
        1,
        'started',
        360,
        NOW() - INTERVAL '7 days'
    ),
    (
        8,
        1,
        'completed',
        360,
        NOW() - INTERVAL '7 days'
    ),
    (
        8,
        2,
        'viewed',
        15,
        NOW() - INTERVAL '7 days'
    ),
    (
        8,
        2,
        'started',
        300,
        NOW() - INTERVAL '7 days'
    ),
    (
        8,
        2,
        'completed',
        300,
        NOW() - INTERVAL '7 days'
    ),
    (
        8,
        3,
        'viewed',
        10,
        NOW() - INTERVAL '7 days'
    ),
    (
        8,
        3,
        'started',
        240,
        NOW() - INTERVAL '7 days'
    ),
    (
        8,
        3,
        'completed',
        240,
        NOW() - INTERVAL '7 days'
    ),

-- Быстро проходит все (6 дней назад)
(
    8,
    4,
    'viewed',
    15,
    NOW() - INTERVAL '6 days'
),
(
    8,
    4,
    'completed',
    360,
    NOW() - INTERVAL '6 days'
),
(
    8,
    5,
    'viewed',
    20,
    NOW() - INTERVAL '6 days'
),
(
    8,
    5,
    'completed',
    450,
    NOW() - INTERVAL '6 days'
),
(
    8,
    6,
    'viewed',
    25,
    NOW() - INTERVAL '6 days'
),
(
    8,
    6,
    'completed',
    720,
    NOW() - INTERVAL '6 days'
),

-- Химия (3 дня назад)
(
    8,
    42,
    'viewed',
    30,
    NOW() - INTERVAL '3 days'
),
(
    8,
    42,
    'completed',
    600,
    NOW() - INTERVAL '3 days'
),
(
    8,
    43,
    'viewed',
    25,
    NOW() - INTERVAL '3 days'
),
(
    8,
    43,
    'completed',
    480,
    NOW() - INTERVAL '3 days'
);

-- Ученик 9 (Максим) - младший, медленно учится
INSERT INTO
    student_interactions (
        student_id,
        resource_id,
        interaction_type,
        duration,
        timestamp
    )
VALUES
    -- Простая математика (5 дней назад)
    (
        9,
        1,
        'viewed',
        60,
        NOW() - INTERVAL '5 days'
    ),
    (
        9,
        1,
        'started',
        900,
        NOW() - INTERVAL '5 days'
    ),
    (
        9,
        1,
        'completed',
        900,
        NOW() - INTERVAL '5 days'
    ),
    (
        9,
        2,
        'viewed',
        50,
        NOW() - INTERVAL '4 days'
    ),
    (
        9,
        2,
        'started',
        720,
        NOW() - INTERVAL '4 days'
    ),
    (
        9,
        2,
        'completed',
        720,
        NOW() - INTERVAL '4 days'
    ),

-- Геометрия (3 дня назад)
(
    9,
    11,
    'viewed',
    45,
    NOW() - INTERVAL '3 days'
),
(
    9,
    11,
    'started',
    600,
    NOW() - INTERVAL '3 days'
),
(
    9,
    11,
    'completed',
    600,
    NOW() - INTERVAL '3 days'
),

-- Биология (вчера)
(
    9,
    47,
    'viewed',
    40,
    NOW() - INTERVAL '1 day'
),
(
    9,
    47,
    'started',
    600,
    NOW() - INTERVAL '1 day'
),
(
    9,
    47,
    'completed',
    600,
    NOW() - INTERVAL '1 day'
),

-- Попытка сложного (пропустил)
(
    9,
    8,
    'viewed',
    30,
    NOW() - INTERVAL '2 days'
),
(
    9,
    8,
    'started',
    180,
    NOW() - INTERVAL '2 days'
),
(
    9,
    8,
    'skipped',
    180,
    NOW() - INTERVAL '2 days'
);

-- Ученик 10 (Камила) - гуманитарий
INSERT INTO
    student_interactions (
        student_id,
        resource_id,
        interaction_type,
        duration,
        timestamp
    )
VALUES
    -- Литература (последние 4 дня)
    (
        10,
        23,
        'viewed',
        55,
        NOW() - INTERVAL '4 days'
    ),
    (
        10,
        23,
        'started',
        1080,
        NOW() - INTERVAL '4 days'
    ),
    (
        10,
        23,
        'completed',
        1080,
        NOW() - INTERVAL '4 days'
    ),
    (
        10,
        23,
        'rated',
        0,
        NOW() - INTERVAL '4 days'
    ),
    (
        10,
        26,
        'viewed',
        40,
        NOW() - INTERVAL '3 days'
    ),
    (
        10,
        26,
        'bookmarked',
        0,
        NOW() - INTERVAL '3 days'
    ),
    (
        10,
        26,
        'started',
        600,
        NOW() - INTERVAL '3 days'
    ),
    (
        10,
        26,
        'completed',
        600,
        NOW() - INTERVAL '3 days'
    ),
    (
        10,
        28,
        'viewed',
        65,
        NOW() - INTERVAL '2 days'
    ),
    (
        10,
        28,
        'started',
        1200,
        NOW() - INTERVAL '2 days'
    ),
    (
        10,
        28,
        'completed',
        1200,
        NOW() - INTERVAL '2 days'
    ),

-- История (вчера)
(
    10,
    32,
    'viewed',
    50,
    NOW() - INTERVAL '1 day'
),
(
    10,
    32,
    'started',
    900,
    NOW() - INTERVAL '1 day'
),
(
    10,
    32,
    'completed',
    900,
    NOW() - INTERVAL '1 day'
),
(
    10,
    37,
    'viewed',
    60,
    NOW() - INTERVAL '1 day'
),
(
    10,
    37,
    'started',
    1080,
    NOW() - INTERVAL '1 day'
),
(
    10,
    37,
    'completed',
    1080,
    NOW() - INTERVAL '1 day'
),

-- Критическое мышление
(
    10,
    14,
    'viewed',
    45,
    NOW() - INTERVAL '3 days'
),
(
    10,
    14,
    'started',
    720,
    NOW() - INTERVAL '3 days'
),
(
    10,
    14,
    'completed',
    720,
    NOW() - INTERVAL '3 days'
);

COMMIT;

-- ========================================
-- ИТОГО: ~150 взаимодействий
-- ========================================
-- Теперь есть полная история обучения для аналитики