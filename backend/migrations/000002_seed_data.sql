-- ========================================
-- SEED DATA: Sample content for testing
-- ========================================

-- ========================================
-- 1. USERS (тестовые пользователи)
-- ========================================
-- Password для всех: "password123" (bcrypt hash)
INSERT INTO
    users (
        email,
        password_hash,
        role,
        full_name,
        is_active
    )
VALUES
    -- Админ
    (
        'admin@edu.kz',
        '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
        'admin',
        'Админ Главный',
        true
    ),

-- Учителя
(
    'teacher1@edu.kz',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'teacher',
    'Анна Смирнова',
    true
),
(
    'teacher2@edu.kz',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'teacher',
    'Петр Иванов',
    true
),
(
    'teacher3@edu.kz',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'teacher',
    'Мария Козлова',
    true
),

-- Ученики
(
    'student1@edu.kz',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'student',
    'Иван Петров',
    true
),
(
    'student2@edu.kz',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'student',
    'Алия Нурланова',
    true
),
(
    'student3@edu.kz',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'student',
    'Дмитрий Сидоров',
    true
),
(
    'student4@edu.kz',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'student',
    'Айгерим Асан',
    true
),
(
    'student5@edu.kz',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'student',
    'Максим Волков',
    true
),
(
    'student6@edu.kz',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'student',
    'Камила Темиргали',
    true
);

-- ========================================
-- 2. STUDENT PROFILES (профили учеников)
-- ========================================
INSERT INTO
    student_profiles (
        user_id,
        grade,
        age_group,
        interests,
        learning_style,
        preferences
    )
VALUES (
        5,
        7,
        'middle',
        ARRAY[
            'math',
            'physics',
            'programming'
        ],
        'visual',
        '{"difficulty_preference": "medium", "preferred_resource_types": ["video", "interactive"]}'::jsonb
    ),
    (
        6,
        9,
        'senior',
        ARRAY[
            'literature',
            'history',
            'art'
        ],
        'auditory',
        '{"difficulty_preference": "hard", "preferred_resource_types": ["reading", "video"]}'::jsonb
    ),
    (
        7,
        6,
        'middle',
        ARRAY[
            'science',
            'biology',
            'chemistry'
        ],
        'kinesthetic',
        '{"difficulty_preference": "easy", "preferred_resource_types": ["interactive", "exercise"]}'::jsonb
    ),
    (
        8,
        10,
        'senior',
        ARRAY[
            'math',
            'physics',
            'chemistry'
        ],
        'visual',
        '{"difficulty_preference": "hard", "preferred_resource_types": ["reading", "video"]}'::jsonb
    ),
    (
        9,
        5,
        'junior',
        ARRAY['math', 'art', 'music'],
        'visual',
        '{"difficulty_preference": "easy", "preferred_resource_types": ["video", "interactive"]}'::jsonb
    ),
    (
        10,
        8,
        'middle',
        ARRAY[
            'history',
            'geography',
            'literature'
        ],
        'auditory',
        '{"difficulty_preference": "medium", "preferred_resource_types": ["reading", "quiz"]}'::jsonb
    );

-- ========================================
-- 3. COURSES (курсы)
-- ========================================
INSERT INTO
    courses (
        title,
        description,
        created_by,
        difficulty_level,
        age_group,
        subject,
        is_published,
        thumbnail_url
    )
VALUES
    -- Математика
    (
        'Основы алгебры',
        'Введение в алгебру: уравнения, неравенства, функции',
        2,
        3,
        'middle',
        'mathematics',
        true,
        '/images/algebra_basics.jpg'
    ),
    (
        'Геометрия для начинающих',
        'Основные геометрические фигуры и их свойства',
        2,
        2,
        'junior',
        'mathematics',
        true,
        '/images/geometry.jpg'
    ),
    (
        'Квадратные уравнения',
        'Решение квадратных уравнений различными способами',
        2,
        4,
        'senior',
        'mathematics',
        true,
        '/images/quadratic.jpg'
    ),

-- Критическое мышление
(
    'Критическое мышление: Основы',
    'Развитие навыков анализа и логики',
    3,
    3,
    'middle',
    'logic',
    true,
    '/images/critical_thinking.jpg'
),
(
    'Логические задачи',
    'Сборник логических задач разной сложности',
    3,
    4,
    'senior',
    'logic',
    true,
    '/images/logic_puzzles.jpg'
),

-- Физика
(
    'Основы физики',
    'Механика, силы, движение',
    2,
    3,
    'middle',
    'physics',
    true,
    '/images/physics_basics.jpg'
),
(
    'Электричество и магнетизм',
    'Изучение электрических и магнитных явлений',
    2,
    5,
    'senior',
    'physics',
    true,
    '/images/electricity.jpg'
),

-- Литература
(
    'Анализ литературных произведений',
    'Учимся анализировать тексты',
    4,
    3,
    'middle',
    'literature',
    true,
    '/images/literature.jpg'
),
(
    'Русская классика',
    'Изучение произведений русских классиков',
    4,
    4,
    'senior',
    'literature',
    true,
    '/images/russian_classics.jpg'
),

-- История
(
    'История Казахстана',
    'Основные этапы развития Казахстана',
    4,
    2,
    'middle',
    'history',
    true,
    '/images/kazakhstan_history.jpg'
),
(
    'Всемирная история',
    'Ключевые события мировой истории',
    4,
    3,
    'senior',
    'history',
    true,
    '/images/world_history.jpg'
),

-- Химия
(
    'Основы химии',
    'Атомы, молекулы, химические реакции',
    3,
    3,
    'middle',
    'chemistry',
    true,
    '/images/chemistry.jpg'
),

-- Биология
(
    'Биология: Живые организмы',
    'Изучение живой природы',
    3,
    2,
    'junior',
    'biology',
    true,
    '/images/biology.jpg'
);

-- ========================================
-- 4. MODULES (модули курсов)
-- ========================================

-- Курс 1: Основы алгебры
INSERT INTO
    modules (
        course_id,
        title,
        description,
        order_index
    )
VALUES (
        1,
        'Введение в алгебру',
        'Основные понятия и термины',
        1
    ),
    (
        1,
        'Линейные уравнения',
        'Решение уравнений первой степени',
        2
    ),
    (
        1,
        'Системы уравнений',
        'Решение систем линейных уравнений',
        3
    ),
    (
        1,
        'Неравенства',
        'Линейные и квадратные неравенства',
        4
    );

-- Курс 2: Геометрия для начинающих
INSERT INTO
    modules (
        course_id,
        title,
        description,
        order_index
    )
VALUES (
        2,
        'Точки и прямые',
        'Основные геометрические объекты',
        1
    ),
    (
        2,
        'Углы',
        'Виды углов и их измерение',
        2
    ),
    (
        2,
        'Треугольники',
        'Свойства треугольников',
        3
    ),
    (
        2,
        'Четырехугольники',
        'Квадраты, прямоугольники, параллелограммы',
        4
    );

-- Курс 3: Квадратные уравнения
INSERT INTO
    modules (
        course_id,
        title,
        description,
        order_index
    )
VALUES (
        3,
        'Определение квадратного уравнения',
        'Что такое квадратное уравнение',
        1
    ),
    (
        3,
        'Формула корней',
        'Дискриминант и формула корней',
        2
    ),
    (
        3,
        'Теорема Виета',
        'Связь корней и коэффициентов',
        3
    );

-- Курс 4: Критическое мышление
INSERT INTO
    modules (
        course_id,
        title,
        description,
        order_index
    )
VALUES (
        4,
        'Что такое критическое мышление',
        'Основы и важность',
        1
    ),
    (
        4,
        'Анализ информации',
        'Как оценивать источники',
        2
    ),
    (
        4,
        'Логические ошибки',
        'Распространенные заблуждения',
        3
    ),
    (
        4,
        'Аргументация',
        'Как строить убедительные аргументы',
        4
    );

-- Курс 5: Логические задачи
INSERT INTO
    modules (
        course_id,
        title,
        description,
        order_index
    )
VALUES (
        5,
        'Простые логические задачи',
        'Разминка для ума',
        1
    ),
    (
        5,
        'Математическая логика',
        'Задачи на числа и последовательности',
        2
    ),
    (
        5,
        'Задачи на внимательность',
        'Развитие наблюдательности',
        3
    );

-- Курс 6: Основы физики
INSERT INTO
    modules (
        course_id,
        title,
        description,
        order_index
    )
VALUES (
        6,
        'Кинематика',
        'Описание движения',
        1
    ),
    (
        6,
        'Динамика',
        'Законы Ньютона',
        2
    ),
    (
        6,
        'Работа и энергия',
        'Энергия и ее виды',
        3
    );

-- ========================================
-- 5. RESOURCES (учебные материалы)
-- ========================================

-- Модуль 1: Введение в алгебру
INSERT INTO
    resources (
        module_id,
        title,
        content,
        resource_type,
        difficulty,
        estimated_time,
        file_url
    )
VALUES (
        1,
        'Что такое алгебра?',
        'Алгебра - это раздел математики, изучающий общие свойства действий над различными величинами...',
        'reading',
        2,
        10,
        NULL
    ),
    (
        1,
        'Основные термины',
        'Переменная, константа, выражение, уравнение...',
        'reading',
        2,
        8,
        NULL
    ),
    (
        1,
        'Тест: Основы алгебры',
        '{"questions": ["Что такое переменная?", "Чем отличается уравнение от выражения?"]}',
        'quiz',
        2,
        5,
        NULL
    );

-- Модуль 2: Линейные уравнения
INSERT INTO
    resources (
        module_id,
        title,
        content,
        resource_type,
        difficulty,
        estimated_time
    )
VALUES (
        2,
        'Понятие линейного уравнения',
        'Линейное уравнение - это уравнение вида ax + b = 0...',
        'reading',
        3,
        12
    ),
    (
        2,
        'Видео: Как решать линейные уравнения',
        'Пошаговый разбор примеров',
        'video',
        3,
        15
    ),
    (
        2,
        'Упражнение: Решите уравнения',
        '{"equations": ["2x + 3 = 7", "5x - 1 = 14", "-3x + 6 = 0"]}',
        'exercise',
        3,
        20
    ),
    (
        2,
        'Тест: Линейные уравнения',
        '{"questions_count": 10}',
        'quiz',
        3,
        15
    );

-- Модуль 3: Системы уравнений
INSERT INTO
    resources (
        module_id,
        title,
        content,
        resource_type,
        difficulty,
        estimated_time
    )
VALUES (
        3,
        'Системы линейных уравнений',
        'Система уравнений - это несколько уравнений с общими переменными...',
        'reading',
        4,
        15
    ),
    (
        3,
        'Метод подстановки',
        'Как решать системы методом подстановки',
        'reading',
        4,
        12
    ),
    (
        3,
        'Интерактивное упражнение',
        'Решайте системы уравнений онлайн с проверкой',
        'interactive',
        4,
        25
    );

-- Модуль 5: Точки и прямые (Геометрия)
INSERT INTO
    resources (
        module_id,
        title,
        content,
        resource_type,
        difficulty,
        estimated_time
    )
VALUES (
        5,
        'Основные понятия геометрии',
        'Точка, прямая, отрезок, луч...',
        'reading',
        1,
        8
    ),
    (
        5,
        'Видео: Рисуем геометрические фигуры',
        'Практическое занятие',
        'video',
        1,
        10
    ),
    (
        5,
        'Тест: Геометрические объекты',
        '{"questions_count": 5}',
        'quiz',
        1,
        8
    );

-- Модуль 9: Что такое критическое мышление
INSERT INTO
    resources (
        module_id,
        title,
        content,
        resource_type,
        difficulty,
        estimated_time
    )
VALUES (
        9,
        'Введение в критическое мышление',
        'Что это и почему это важно',
        'reading',
        3,
        10
    ),
    (
        9,
        'Видео: Примеры критического мышления',
        'Разбор реальных ситуаций',
        'video',
        3,
        15
    ),
    (
        9,
        'Упражнение: Оцените аргументы',
        'Анализ различных утверждений',
        'exercise',
        3,
        20
    );

-- Модуль 10: Анализ информации
INSERT INTO
    resources (
        module_id,
        title,
        content,
        resource_type,
        difficulty,
        estimated_time
    )
VALUES (
        10,
        'Как проверять источники',
        'Критерии надежности информации',
        'reading',
        3,
        12
    ),
    (
        10,
        'Практика: Фейковые новости',
        'Определите правду и ложь',
        'interactive',
        4,
        18
    ),
    (
        10,
        'Тест: Анализ информации',
        '{"questions_count": 8}',
        'quiz',
        3,
        10
    );

-- Модуль 14: Кинематика (Физика)
INSERT INTO
    resources (
        module_id,
        title,
        content,
        resource_type,
        difficulty,
        estimated_time
    )
VALUES (
        14,
        'Движение и скорость',
        'Основы кинематики',
        'reading',
        3,
        15
    ),
    (
        14,
        'Видео: Расчет скорости',
        'Примеры задач',
        'video',
        3,
        12
    ),
    (
        14,
        'Задачи на движение',
        'Решите 10 задач на скорость и время',
        'exercise',
        4,
        25
    );

-- ========================================
-- 6. QUESTIONS (вопросы для тестов)
-- ========================================

-- Тест: Основы алгебры (resource_id = 3)
INSERT INTO
    questions (
        resource_id,
        question_text,
        question_type,
        correct_answer,
        options,
        explanation,
        points,
        order_index
    )
VALUES (
        3,
        'Что такое переменная в алгебре?',
        'multiple_choice',
        'B',
        '{"A": "Постоянное число", "B": "Величина, которая может принимать разные значения", "C": "Знак операции", "D": "Результат вычисления"}'::jsonb,
        'Переменная - это символ, обозначающий величину, которая может принимать различные значения.',
        5,
        1
    ),
    (
        3,
        'Чем уравнение отличается от выражения?',
        'multiple_choice',
        'C',
        '{"A": "Ничем не отличается", "B": "Уравнение короче", "C": "Уравнение содержит знак равенства", "D": "Выражение сложнее"}'::jsonb,
        'Уравнение всегда содержит знак равенства (=), а выражение - нет.',
        5,
        2
    ),
    (
        3,
        'Верно ли, что 2x + 3 - это уравнение?',
        'true_false',
        'false',
        '{"true": "Да", "false": "Нет"}'::jsonb,
        'Это выражение, а не уравнение, так как нет знака равенства.',
        3,
        3
    );

-- Тест: Линейные уравнения (resource_id = 6)
INSERT INTO
    questions (
        resource_id,
        question_text,
        question_type,
        correct_answer,
        options,
        explanation,
        points,
        order_index
    )
VALUES (
        6,
        'Решите уравнение: 2x + 4 = 10',
        'multiple_choice',
        'B',
        '{"A": "x = 2", "B": "x = 3", "C": "x = 4", "D": "x = 5"}'::jsonb,
        '2x = 10 - 4 = 6, поэтому x = 3',
        10,
        1
    ),
    (
        6,
        'Чему равен x в уравнении: 5x - 15 = 0?',
        'multiple_choice',
        'A',
        '{"A": "x = 3", "B": "x = -3", "C": "x = 5", "D": "x = 0"}'::jsonb,
        '5x = 15, значит x = 3',
        10,
        2
    ),
    (
        6,
        'Имеет ли уравнение 0x + 5 = 5 решение?',
        'true_false',
        'false',
        '{"true": "Да, x может быть любым числом", "false": "Нет, это тождество"}'::jsonb,
        'Это тождество, верное при любом x, а не уравнение с конкретным решением.',
        5,
        3
    );

-- Тест: Геометрические объекты (resource_id = 17)
INSERT INTO
    questions (
        resource_id,
        question_text,
        question_type,
        correct_answer,
        options,
        explanation,
        points,
        order_index
    )
VALUES (
        17,
        'Сколько точек нужно, чтобы провести прямую?',
        'multiple_choice',
        'B',
        '{"A": "Одна", "B": "Две", "C": "Три", "D": "Бесконечно много"}'::jsonb,
        'Через две точки можно провести единственную прямую.',
        5,
        1
    ),
    (
        17,
        'Что такое отрезок?',
        'multiple_choice',
        'C',
        '{"A": "Бесконечная линия", "B": "Линия с одним концом", "C": "Часть прямой между двумя точками", "D": "Кривая линия"}'::jsonb,
        'Отрезок - это часть прямой, ограниченная двумя точками.',
        5,
        2
    );

-- ========================================
-- 7. STUDENT PROGRESS (прогресс учеников)
-- ========================================

-- Ученик 1 (Иван) - прошел несколько уроков по алгебре
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
        8,
        1,
        NOW() - INTERVAL '5 days'
    ),
    (
        5,
        2,
        'completed',
        95,
        7,
        1,
        NOW() - INTERVAL '5 days'
    ),
    (
        5,
        3,
        'completed',
        80,
        12,
        2,
        NOW() - INTERVAL '4 days'
    ),
    (
        5,
        4,
        'completed',
        90,
        10,
        1,
        NOW() - INTERVAL '3 days'
    ),
    (
        5,
        5,
        'in_progress',
        60,
        8,
        1,
        NULL
    );

-- Ученик 2 (Алия) - любит литературу
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
        18,
        'completed',
        95,
        15,
        1,
        NOW() - INTERVAL '2 days'
    ),
    (
        6,
        19,
        'completed',
        88,
        20,
        1,
        NOW() - INTERVAL '1 day'
    ),
    (
        6,
        20,
        'in_progress',
        70,
        10,
        1,
        NULL
    );

-- Ученик 3 (Дмитрий) - начал геометрию
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
        15,
        'completed',
        75,
        9,
        1,
        NOW() - INTERVAL '1 day'
    ),
    (
        7,
        16,
        'completed',
        82,
        11,
        1,
        NOW() - INTERVAL '1 day'
    ),
    (
        7,
        17,
        'completed',
        70,
        15,
        2,
        NOW() - INTERVAL '12 hours'
    );

-- ========================================
-- 8. RESOURCE RATINGS (оценки ресурсов)
-- ========================================

-- Ученики оценивают материалы
INSERT INTO
    resource_ratings (
        student_id,
        resource_id,
        rating,
        review
    )
VALUES (
        5,
        1,
        5,
        'Очень понятное объяснение!'
    ),
    (
        5,
        2,
        5,
        'Хорошо структурирован материал'
    ),
    (
        5,
        3,
        4,
        'Тест немного сложный, но полезный'
    ),
    (5, 4, 5, 'Отличное видео'),
    (
        6,
        18,
        5,
        'Интересный подход к анализу'
    ),
    (
        6,
        19,
        4,
        'Много полезной информации'
    ),
    (
        7,
        15,
        4,
        'Понятно для начинающих'
    ),
    (
        7,
        16,
        5,
        'Видео помогло разобраться'
    ),
    (
        7,
        17,
        3,
        'Тест оказался сложнее, чем думал'
    );

-- ========================================
-- 9. STUDENT INTERACTIONS (взаимодействия)
-- ========================================

-- История действий учеников
INSERT INTO
    student_interactions (
        student_id,
        resource_id,
        interaction_type,
        duration,
        timestamp
    )
VALUES
    -- Иван
    (
        5,
        1,
        'viewed',
        30,
        NOW() - INTERVAL '5 days 2 hours'
    ),
    (
        5,
        1,
        'started',
        480,
        NOW() - INTERVAL '5 days 2 hours'
    ),
    (
        5,
        1,
        'completed',
        480,
        NOW() - INTERVAL '5 days 1 hour'
    ),
    (
        5,
        2,
        'viewed',
        20,
        NOW() - INTERVAL '5 days'
    ),
    (
        5,
        2,
        'started',
        420,
        NOW() - INTERVAL '5 days'
    ),
    (
        5,
        2,
        'completed',
        420,
        NOW() - INTERVAL '5 days'
    ),

-- Алия
(
    6,
    18,
    'viewed',
    45,
    NOW() - INTERVAL '2 days'
),
(
    6,
    18,
    'started',
    900,
    NOW() - INTERVAL '2 days'
),
(
    6,
    18,
    'completed',
    900,
    NOW() - INTERVAL '2 days'
),
(
    6,
    19,
    'bookmarked',
    0,
    NOW() - INTERVAL '1 day'
),

-- Дмитрий
(
    7,
    15,
    'viewed',
    15,
    NOW() - INTERVAL '1 day'
),
(
    7,
    15,
    'started',
    540,
    NOW() - INTERVAL '1 day'
),
(
    7,
    17,
    'skipped',
    60,
    NOW() - INTERVAL '12 hours'
);

-- ========================================
-- 10. STUDENT SKILLS (навыки учеников)
-- ========================================

INSERT INTO
    student_skills (
        student_id,
        skill_name,
        proficiency_level
    )
VALUES
    -- Иван (хорош в алгебре)
    (5, 'algebra', 0.75),
    (5, 'geometry', 0.45),
    (5, 'problem_solving', 0.68),
    (5, 'critical_thinking', 0.55),

-- Алия (сильна в гуманитарных науках)
(
    6,
    'reading_comprehension',
    0.85
),
(6, 'writing', 0.78),
(6, 'analysis', 0.80),
(6, 'critical_thinking', 0.72),

-- Дмитрий (начинающий)
(7, 'geometry', 0.50), (7, 'algebra', 0.40), (7, 'logic', 0.48);

-- ========================================
-- 11. STUDENT POINTS (очки и уровни)
-- ========================================

INSERT INTO
    student_points (
        student_id,
        total_points,
        level,
        experience
    )
VALUES (5, 450, 5, 50), -- Иван - активный ученик
    (6, 320, 4, 20), -- Алия
    (7, 180, 2, 80), -- Дмитрий - новичок
    (8, 520, 6, 20), -- Айгерим
    (9, 95, 1, 95), -- Максим - только начал
    (10, 280, 3, 80);
-- Камила

-- ========================================
-- 12. STUDENT ACHIEVEMENTS (полученные достижения)
-- ========================================

-- Ученики получили достижения
INSERT INTO
    student_achievements (
        student_id,
        achievement_id,
        earned_at
    )
VALUES
    -- Иван
    (
        5,
        1,
        NOW() - INTERVAL '5 days'
    ), -- Первые шаги
    (
        5,
        2,
        NOW() - INTERVAL '3 days'
    ), -- Быстрый ученик
    (
        5,
        3,
        NOW() - INTERVAL '4 days'
    ), -- Идеальный результат

-- Алия
(
    6,
    1,
    NOW() - INTERVAL '10 days'
),
(
    6,
    3,
    NOW() - INTERVAL '2 days'
),

-- Дмитрий
( 7, 1, NOW() - INTERVAL '1 day' );

-- ========================================
-- 13. LEADERBOARD (рейтинги)
-- ========================================

INSERT INTO
    leaderboard (
        student_id,
        period,
        rank,
        points
    )
VALUES
    -- Рейтинг за неделю
    (8, 'weekly', 1, 420),
    (5, 'weekly', 2, 380),
    (6, 'weekly', 3, 290),
    (10, 'weekly', 4, 240),
    (7, 'weekly', 5, 180),
    (9, 'weekly', 6, 95),

-- Общий рейтинг
(8, 'all_time', 1, 520),
(5, 'all_time', 2, 450),
(6, 'all_time', 3, 320),
(10, 'all_time', 4, 280),
(7, 'all_time', 5, 180),
(9, 'all_time', 6, 95);

-- ========================================
-- SEED DATA COMPLETE!
-- ========================================
-- Создано:
-- - 10 пользователей (1 админ, 3 учителя, 6 учеников)
-- - 13 курсов (математика, физика, литература, история, химия, биология)
-- - 15 модулей
-- - 23 ресурса (уроки, видео, тесты, упражнения)
-- - 8 вопросов для тестов
-- - Прогресс для 3 учеников
-- - О