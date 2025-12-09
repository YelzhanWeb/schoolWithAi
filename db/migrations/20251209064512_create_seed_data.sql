-- +goose Up
-- +goose StatementBegin
INSERT INTO
    subjects (id, slug, name_ru, name_kz)
VALUES (
        'math-12345',
        'math',
        'Математика',
        'Математика'
    ),
    (
        'kaz_lang-12345',
        'kaz_lang',
        'Казахский язык',
        'Қазақ тілі'
    ),
    (
        'history_kz-12345',
        'history_kz',
        'История Казахстана',
        'Қазақстан тарихы'
    ),
    (
        'physics-12345',
        'physics',
        'Физика',
        'Физика'
    ),
    (
        'informatics-12345',
        'informatics',
        'Информатика',
        'Информатика'
    ),
    (
        'english-12345',
        'english',
        'Английский язык',
        'Ағылшын тілі'
    ),
    (
        'biology-12345',
        'biology',
        'Биология',
        'Биология'
    ),
    (
        'chemistry-12345',
        'chemistry',
        'Химия',
        'Химия'
    ),
    (
        'geography-12345',
        'geography',
        'География',
        'География'
    ),
    (
        'literature-12345',
        'literature',
        'Русская литература',
        'Орыс әдебиеті'
    ),
    (
        'world_history-12345',
        'world_history',
        'Всемирная история',
        'Дүниежүзі тарихы'
    ),
    (
        'computer_science-12345',
        'computer_science',
        'Компьютерные науки',
        'Компьютерлік ғылымдар'
    ),
    (
        'logic-12345',
        'logic',
        'Логика',
        'Логика'
    ),
    (
        'economics-12345',
        'economics',
        'Экономика',
        'Экономика'
    ),
    (
        'finance-12345',
        'finance',
        'Финансовая грамотность',
        'Қаржылық сауаттылық'
    ),
    (
        'law-12345',
        'law',
        'Право',
        'Құқық негіздері'
    ),
    (
        'art-12345',
        'art',
        'Искусство',
        'Өнер'
    ),
    (
        'music-12345',
        'music',
        'Музыка',
        'Музыка'
    ),
    (
        'pe-12345',
        'physical_education',
        'Физическая культура',
        'Дене шынықтыру'
    ),
    (
        'philosophy-12345',
        'philosophy',
        'Философия',
        'Философия'
    ),
    (
        'sociology-12345',
        'sociology',
        'Социология',
        'Социология'
    ),
    (
        'psychology-12345',
        'psychology',
        'Психология',
        'Психология'
    ),
    (
        'env_science-12345',
        'environment',
        'Экология',
        'Экология'
    ),
    (
        'astronomy-12345',
        'astronomy',
        'Астрономия',
        'Астрономия'
    ),
    (
        'statistics-12345',
        'statistics',
        'Статистика',
        'Статистика'
    );

INSERT INTO
    tags (name, slug)
VALUES ('Алгебра', 'algebra'),
    ('Геометрия', 'geometry'),
    ('Механика', 'mechanics'),
    (
        'Электричество и магнетизм',
        'electricity_magnetism'
    ),
    ('Грамматика', 'grammar'),
    (
        'Чтение и анализ текста',
        'text_analysis'
    ),
    (
        'Программирование',
        'programming'
    ),
    ('Алгоритмы', 'algorithms'),
    (
        'Статистика и вероятности',
        'statistics'
    ),
    (
        'История Казахстана',
        'history_kz'
    ),
    ('Литература', 'literature'),
    (
        'Лексика и словарный запас',
        'vocabulary'
    ),
    (
        'Подготовка к экзаменам',
        'exam_prep'
    ),
    (
        'Олимпиадный уровень',
        'olympiad'
    ),
    (
        'Картографический анализ',
        'map_analysis'
    ),
    (
        'Химия: основы',
        'chemistry_basics'
    ),
    (
        'Биология: анатомия',
        'biology_anatomy'
    ),
    (
        'География: карты',
        'geography_maps'
    ),
    (
        'Физика: оптика',
        'physics_optics'
    ),
    (
        'История: даты и события',
        'history_dates'
    ),
    (
        'Экономика: базовые понятия',
        'economics_basics'
    ),
    (
        'Финансовая грамотность',
        'finance_literacy'
    ),
    (
        'Психология: основы',
        'psychology_basics'
    ),
    (
        'Социология: исследования',
        'sociology_studies'
    ),
    (
        'Философия: термины',
        'philosophy_terms'
    ),
    (
        'Астрономия: космос',
        'astronomy_space'
    ),
    (
        'Биология: экология',
        'biology_ecology'
    ),
    (
        'Программирование: Python задачи',
        'python_tasks'
    ),
    (
        'Алгоритмы: графы',
        'algorithms_graphs'
    ),
    (
        'Геометрия: стереометрия',
        'geometry_3d'
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd