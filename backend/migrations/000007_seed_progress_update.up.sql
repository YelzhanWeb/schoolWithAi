-- ========================================
-- SEED: Update Student Progress (More Data)
-- ========================================

-- Удаляем старые данные прогресса (они уже были в 000002)
DELETE FROM student_progress;

-- Ученик 5 (Иван) - математика и физика
INSERT INTO
    student_progress (
        student_id,
        resource_id,
        status,
        score,
        time_spent,
        attempts,
        completed_at
    )
VALUES (
        5,
        1,
        'completed',
        100,
        10,
        1,
        NOW() - INTERVAL '7 days'
    ),
    (
        5,
        2,
        'completed',
        95,
        8,
        1,
        NOW() - INTERVAL '7 days'
    ),
    (
        5,
        3,
        'completed',
        85,
        15,
        2,
        NOW() - INTERVAL '6 days'
    ),
    (
        5,
        4,
        'completed',
        90,
        12,
        1,
        NOW() - INTERVAL '6 days'
    ),
    (
        5,
        5,
        'completed',
        95,
        15,
        1,
        NOW() - INTERVAL '5 days'
    ),
    (
        5,
        6,
        'completed',
        88,
        20,
        1,
        NOW() - INTERVAL '5 days'
    ),
    (
        5,
        7,
        'completed',
        92,
        10,
        1,
        NOW() - INTERVAL '4 days'
    ),
    (
        5,
        8,
        'completed',
        85,
        15,
        1,
        NOW() - INTERVAL '4 days'
    ),
    (
        5,
        9,
        'completed',
        80,
        18,
        2,
        NOW() - INTERVAL '3 days'
    ),
    (
        5,
        10,
        'completed',
        90,
        25,
        1,
        NOW() - INTERVAL '3 days'
    ),
    (
        5,
        20,
        'completed',
        95,
        15,
        1,
        NOW() - INTERVAL '3 days'
    ),
    (
        5,
        21,
        'completed',
        90,
        12,
        1,
        NOW() - INTERVAL '3 days'
    ),
    (
        5,
        22,
        'completed',
        85,
        25,
        1,
        NOW() - INTERVAL '2 days'
    ),
    (
        5,
        42,
        'completed',
        78,
        15,
        1,
        NOW() - INTERVAL '2 days'
    ),
    (
        5,
        43,
        'completed',
        82,
        12,
        1,
        NOW() - INTERVAL '1 day'
    ),
    (
        5,
        44,
        'in_progress',
        60,
        10,
        1,
        NULL
    );

-- Ученик 6 (Алия) - литература и история
INSERT INTO
    student_progress (
        student_id,
        resource_id,
        status,
        score,
        time_spent,
        attempts,
        completed_at
    )
VALUES (
        6,
        23,
        'completed',
        95,
        18,
        1,
        NOW() - INTERVAL '5 days'
    ),
    (
        6,
        24,
        'completed',
        92,
        15,
        1,
        NOW() - INTERVAL '4 days'
    ),
    (
        6,
        25,
        'completed',
        88,
        20,
        1,
        NOW() - INTERVAL '4 days'
    ),
    (
        6,
        26,
        'completed',
        95,
        10,
        1,
        NOW() - INTERVAL '3 days'
    ),
    (
        6,
        27,
        'completed',
        90,
        8,
        1,
        NOW() - INTERVAL '3 days'
    ),
    (
        6,
        28,
        'completed',
        98,
        20,
        1,
        NOW() - INTERVAL '2 days'
    ),
    (
        6,
        29,
        'completed',
        95,
        30,
        1,
        NOW() - INTERVAL '1 day'
    ),
    (
        6,
        30,
        'in_progress',
        70,
        15,
        1,
        NULL
    ),
    (
        6,
        32,
        'completed',
        92,
        10,
        1,
        NOW() - INTERVAL '4 days'
    ),
    (
        6,
        33,
        'completed',
        90,
        12,
        1,
        NOW() - INTERVAL '4 days'
    ),
    (
        6,
        37,
        'completed',
        95,
        15,
        1,
        NOW() - INTERVAL '2 days'
    ),
    (
        6,
        38,
        'completed',
        92,
        20,
        1,
        NOW() - INTERVAL '2 days'
    );

-- Ученик 7 (Дмитрий) - начинающий
INSERT INTO
    student_progress (
        student_id,
        resource_id,
        status,
        score,
        time_spent,
        attempts,
        completed_at
    )
VALUES (
        7,
        11,
        'completed',
        70,
        9,
        1,
        NOW() - INTERVAL '3 days'
    ),
    (
        7,
        12,
        'completed',
        75,
        10,
        1,
        NOW() - INTERVAL '3 days'
    ),
    (
        7,
        13,
        'completed',
        65,
        8,
        2,
        NOW() - INTERVAL '2 days'
    ),
    (
        7,
        47,
        'completed',
        85,
        10,
        1,
        NOW() - INTERVAL '1 day'
    ),
    (
        7,
        48,
        'completed',
        80,
        15,
        1,
        NOW() - INTERVAL '1 day'
    ),
    (
        7,
        49,
        'completed',
        75,
        8,
        1,
        NOW() - INTERVAL '1 day'
    ),
    (
        7,
        50,
        'in_progress',
        60,
        6,
        1,
        NULL
    );

-- Ученик 8 (Айгерим) - отличница
INSERT INTO
    student_progress (
        student_id,
        resource_id,
        status,
        score,
        time_spent,
        attempts,
        completed_at
    )
VALUES (
        8,
        1,
        'completed',
        100,
        6,
        1,
        NOW() - INTERVAL '7 days'
    ),
    (
        8,
        2,
        'completed',
        100,
        5,
        1,
        NOW() - INTERVAL '7 days'
    ),
    (
        8,
        3,
        'completed',
        100,
        4,
        1,
        NOW() - INTERVAL '7 days'
    ),
    (
        8,
        4,
        'completed',
        100,
        6,
        1,
        NOW() - INTERVAL '6 days'
    ),
    (
        8,
        5,
        'completed',
        98,
        8,
        1,
        NOW() - INTERVAL '6 days'
    ),
    (
        8,
        6,
        'completed',
        100,
        12,
        1,
        NOW() - INTERVAL '6 days'
    ),
    (
        8,
        7,
        'completed',
        100,
        8,
        1,
        NOW() - INTERVAL '6 days'
    ),
    (
        8,
        8,
        'completed',
        95,
        10,
        1,
        NOW() - INTERVAL '5 days'
    ),
    (
        8,
        9,
        'completed',
        98,
        12,
        1,
        NOW() - INTERVAL '5 days'
    ),
    (
        8,
        10,
        'completed',
        100,
        15,
        1,
        NOW() - INTERVAL '5 days'
    ),
    (
        8,
        20,
        'completed',
        95,
        12,
        1,
        NOW() - INTERVAL '4 days'
    ),
    (
        8,
        21,
        'completed',
        98,
        10,
        1,
        NOW() - INTERVAL '4 days'
    ),
    (
        8,
        22,
        'completed',
        95,
        18,
        1,
        NOW() - INTERVAL '4 days'
    ),
    (
        8,
        42,
        'completed',
        100,
        10,
        1,
        NOW() - INTERVAL '3 days'
    ),
    (
        8,
        43,
        'completed',
        98,
        8,
        1,
        NOW() - INTERVAL '3 days'
    ),
    (
        8,
        44,
        'completed',
        95,
        15,
        1,
        NOW() - INTERVAL '3 days'
    ),
    (
        8,
        45,
        'completed',
        92,
        18,
        1,
        NOW() - INTERVAL '2 days'
    ),
    (
        8,
        46,
        'completed',
        95,
        10,
        1,
        NOW() - INTERVAL '2 days'
    );

-- Ученик 9 (Максим) - младший
INSERT INTO
    student_progress (
        student_id,
        resource_id,
        status,
        score,
        time_spent,
        attempts,
        completed_at
    )
VALUES (
        9,
        1,
        'completed',
        70,
        15,
        2,
        NOW() - INTERVAL '5 days'
    ),
    (
        9,
        2,
        'completed',
        75,
        12,
        1,
        NOW() - INTERVAL '4 days'
    ),
    (
        9,
        11,
        'completed',
        68,
        10,
        1,
        NOW() - INTERVAL '3 days'
    ),
    (
        9,
        12,
        'completed',
        72,
        12,
        1,
        NOW() - INTERVAL '3 days'
    ),
    (
        9,
        47,
        'completed',
        80,
        10,
        1,
        NOW() - INTERVAL '1 day'
    ),
    (
        9,
        48,
        'completed',
        78,
        15,
        1,
        NOW() - INTERVAL '1 day'
    ),
    (
        9,
        50,
        'in_progress',
        50,
        6,
        1,
        NULL
    );

-- Ученик 10 (Камила) - гуманитарий
INSERT INTO
    student_progress (
        student_id,
        resource_id,
        status,
        score,
        time_spent,
        attempts,
        completed_at
    )
VALUES (
        10,
        23,
        'completed',
        95,
        18,
        1,
        NOW() - INTERVAL '4 days'
    ),
    (
        10,
        24,
        'completed',
        90,
        15,
        1,
        NOW() - INTERVAL '4 days'
    ),
    (
        10,
        26,
        'completed',
        92,
        10,
        1,
        NOW() - INTERVAL '3 days'
    ),
    (
        10,
        28,
        'completed',
        98,
        20,
        1,
        NOW() - INTERVAL '2 days'
    ),
    (
        10,
        29,
        'in_progress',
        75,
        15,
        1,
        NULL
    ),
    (
        10,
        32,
        'completed',
        90,
        15,
        1,
        NOW() - INTERVAL '1 day'
    ),
    (
        10,
        33,
        'completed',
        88,
        12,
        1,
        NOW() - INTERVAL '1 day'
    ),
    (
        10,
        37,
        'completed',
        95,
        18,
        1,
        NOW() - INTERVAL '1 day'
    ),
    (
        10,
        14,
        'completed',
        92,
        12,
        1,
        NOW() - INTERVAL '3 days'
    ),
    (
        10,
        15,
        'completed',
        88,
        15,
        1,
        NOW() - INTERVAL '3 days'
    );

COMMIT;

-- ========================================
-- ИТОГО: Обновлен прогресс всех учеников
-- ========================================
-- Теперь у каждого ученика реалистичный прогресс