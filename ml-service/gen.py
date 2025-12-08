"""
Скрипт генерации мок-данных для системы рекомендаций
Совместим с твоей схемой БД (users + first_name, last_name)
"""

import random
import uuid
from datetime import datetime, timedelta

import psycopg2
from psycopg2.extras import execute_values

DB_CONFIG = {
    "host": "localhost",
    "port": 5432,
    "database": "education_platform",
    "user": "admin",
    "password": "admin123"
}

FIRST_NAMES = [
    "Алексей", "Мария", "Иван", "Анна", "Дмитрий",
    "Елена", "Сергей", "Ольга", "Павел", "Наталья"
]

LAST_NAMES = [
    "Иванов", "Петрова", "Сидоров", "Смирнова", "Козлов",
    "Новикова", "Морозов", "Волкова", "Лебедев", "Соколова"
]

COURSE_TITLES = {
    "math-12345": [
        "Основы алгебры",
        "Геометрия для начинающих",
        "Тригонометрия",
        "Математический анализ",
        "Теория вероятностей"
    ],
    "physics-12345": [
        "Механика",
        "Термодинамика",
        "Электричество и магнетизм",
        "Оптика",
        "Квантовая физика"
    ],
    "kaz_lang-12345": [
        "Қазақ тілінің негіздері",
        "Грамматика казахского языка",
        "Казахская литература",
        "Разговорный казахский",
        "Деловой казахский язык"
    ]
}

ACTION_TYPES = ["view", "view", "view", "complete"]  # view чаще


# --------------------------------------------------------------------
# Подключение к БД
# --------------------------------------------------------------------
def create_connection():
    return psycopg2.connect(**DB_CONFIG)


# --------------------------------------------------------------------
# Создание студентов
# --------------------------------------------------------------------
def generate_users(conn, count=30):
    print(f"📝 Генерация {count} студентов...")

    users = []
    profiles = []
    interests = []

    for i in range(count):
        user_id = str(uuid.uuid4())
        email = f"student{i + 1}@test.com"
        first_name = random.choice(FIRST_NAMES)
        last_name = random.choice(LAST_NAMES)

        password_hash = "$2a$10$dummy.hash.for.testing.only"

        grade = random.randint(1, 11)
        level = random.randint(1, 5)
        xp = random.randint(0, 5000)
        weekly_xp = xp % 1000

        users.append((
            user_id, email, password_hash, 'student',
            first_name, last_name,
            'default_avatar.png',
            datetime.now(), datetime.now()
        ))

        profiles.append((
            str(uuid.uuid4()), user_id, grade, xp,
            level, 1, weekly_xp,
            random.randint(0, 10), random.randint(5, 20),
            datetime.now(), datetime.now(), datetime.now()
        ))

        # Интересы — от 1 до 3 предметов
        all_subjects = list(COURSE_TITLES.keys())
        for subj in random.sample(all_subjects, random.randint(1, 3)):
            interests.append((user_id, subj))

    cur = conn.cursor()

    execute_values(cur, """
        INSERT INTO users
        (id, email, password_hash, role, first_name, last_name, avatar_url, created_at, updated_at)
        VALUES %s
        ON CONFLICT (email) DO NOTHING
    """, users)

    execute_values(cur, """
        INSERT INTO student_profiles
        (id, user_id, grade, xp, level, current_league_id, weekly_xp,
         current_streak, max_streak, last_activity_date, created_at, updated_at)
        VALUES %s
        ON CONFLICT (user_id) DO NOTHING
    """, profiles)

    execute_values(cur, """
        INSERT INTO student_interests (user_id, subject_id)
        VALUES %s
        ON CONFLICT DO NOTHING
    """, interests)

    conn.commit()
    print(f"✅ Создано {len(users)} студентов")

    return [u[0] for u in users]


# --------------------------------------------------------------------
# Создание курсов
# --------------------------------------------------------------------
def generate_courses(conn):
    print("📚 Генерация курсов...")

    cur = conn.cursor()

    # Создаём учителя, если нет
    cur.execute("SELECT id FROM users WHERE role = 'teacher' LIMIT 1")
    teacher = cur.fetchone()

    if not teacher:
        teacher_id = str(uuid.uuid4())
        cur.execute("""
            INSERT INTO users
            (id, email, password_hash, role, first_name, last_name, avatar_url, created_at, updated_at)
            VALUES (%s, %s, %s, 'teacher', %s, %s, %s, NOW(), NOW())
        """, (teacher_id, "teacher@test.com", "$2a$10$dummy", "Учитель", "Тестовый", "teacher.png"))
    else:
        teacher_id = teacher[0]

    courses = []
    tags = []

    for subject_id, titles in COURSE_TITLES.items():
        for title in titles:
            cid = str(uuid.uuid4())
            difficulty = random.randint(1, 5)

            courses.append((
                cid, teacher_id, subject_id, title,
                f"Описание курса '{title}'",
                difficulty, "", True, datetime.now()
            ))

            # каждому курсу — 1–2 случайных тега
            for tag_id in random.sample([1, 2, 3], random.randint(1, 2)):
                tags.append((cid, tag_id))

    execute_values(cur, """
        INSERT INTO courses
        (id, author_id, subject_id, title, description,
         difficulty_level, cover_image_url, is_published, created_at)
        VALUES %s
        ON CONFLICT DO NOTHING
    """, courses)

    execute_values(cur, """
        INSERT INTO course_tags (course_id, tag_id)
        VALUES %s
        ON CONFLICT DO NOTHING
    """, tags)

    conn.commit()
    print(f"✅ Создано {len(courses)} курсов")

    return [c[0] for c in courses]


# --------------------------------------------------------------------
# Создание логов активности
# --------------------------------------------------------------------
def generate_interactions(conn, user_ids, course_ids, count=500):
    print(f"🔄 Генерация {count} взаимодействий...")

    logs = []

    for _ in range(count):
        log_id = str(uuid.uuid4())
        user_id = random.choice(user_ids)
        course_id = random.choice(course_ids)
        action = random.choice(ACTION_TYPES)

        created = datetime.now() - timedelta(days=random.randint(0, 30))

        logs.append((
            log_id, user_id, course_id, action,
            '{"duration": 120}', created
        ))

    cur = conn.cursor()

    execute_values(cur, """
        INSERT INTO user_activity_logs
        (id, user_id, course_id, action_type, meta_data, created_at)
        VALUES %s
    """, logs)

    conn.commit()
    print("✅ Взаимодействия созданы")


# --------------------------------------------------------------------
# MAIN
# --------------------------------------------------------------------
def main():
    print("=" * 60)
    print("🚀 Генерация мок-данных")
    print("=" * 60)

    conn = create_connection()
    print("✅ Подключение к БД установлено")

    users = generate_users(conn, 30)
    courses = generate_courses(conn)
    generate_interactions(conn, users, courses, 500)

    print("\n🎉 Готово!")
    print(f"👥 Пользователей: {len(users)}")
    print(f"📚 Курсов: {len(courses)}")
    print("🔄 Логов: 500")
    print("\n💡 Пример user_id:", users[0])


if __name__ == "__main__":
    main()
