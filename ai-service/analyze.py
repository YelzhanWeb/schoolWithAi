import json
import statistics
import matplotlib.pyplot as plt
from collections import defaultdict


def load_json(path):
    """Безопасная загрузка JSON-файла"""
    try:
        with open(path, "r", encoding="utf-8") as f:
            return json.load(f)
    except FileNotFoundError:
        print(f"❌ Файл '{path}' не найден.")
        return None
    except json.JSONDecodeError:
        print(f"❌ Ошибка формата JSON в файле '{path}'.")
        return None


def analyze_performance(
    user_file="/home/aibolat/schoolWithAi/go_recommender/data/users.json",
    grades_file="/home/aibolat/schoolWithAi/go_recommender/data/grades.json",
    courses_file="/home/aibolat/schoolWithAi/go_recommender/data/courses.json"
):
    # 1️⃣ Загрузка данных
    users = load_json(user_file)
    grades = load_json(grades_file)
    courses = load_json(courses_file)

    if not users or not grades or not courses:
        return

    # 2️⃣ Создаём словарь курсов {id: название}
    course_names = {c["id"]: c["name"] for c in courses}

    # 3️⃣ Вывод списка пользователей
    print("\n👥 Список пользователей:\n")
    for user in users:
        print(f"ID: {user['id']} | Имя: {user['name']}")

    # 4️⃣ Ввод ID ученика
    try:
        user_id = int(input("\nВведите ID ученика для анализа: "))
    except ValueError:
        print("❌ Некорректный ввод ID.")
        return

    # 5️⃣ Проверка существования ученика
    selected_user = next((u for u in users if u["id"] == user_id), None)
    if not selected_user:
        print("❌ Пользователь с таким ID не найден.")
        return

    # 6️⃣ Фильтруем оценки по пользователю
    user_grades = [g for g in grades if g["user_id"] == user_id]
    if not user_grades:
        print(f"ℹ️ У пользователя {selected_user['name']} пока нет оценок.")
        return

    # 7️⃣ Анализ по курсам
    print(f"\n📊 Анализ успеваемости ученика: {selected_user['name']} (ID: {user_id})\n")
    course_scores = defaultdict(list)

    for rec in user_grades:
        score = rec.get("score") or rec.get("rating")  # поддержка rating
        if score is not None:
            course_scores[rec["course_id"]].append(score)
        else:
            print(f"⚠️ Запись без оценки: {rec}")

    if not course_scores:
        print("⚠️ Не найдено ни одной валидной оценки.")
        return

    overall_scores = []
    avg_per_course = {}
    for course_id, scores in course_scores.items():
        avg = sum(scores) / len(scores)
        avg_per_course[course_id] = avg
        overall_scores.extend(scores)
        course_name = course_names.get(course_id, f"Курс {course_id}")
        print(f"  • {course_name}: средний балл — {avg:.2f}")

    # 8️⃣ Общая статистика
    avg_score = sum(overall_scores) / len(overall_scores)
    min_score = min(overall_scores)
    max_score = max(overall_scores)
    std_dev = statistics.stdev(overall_scores) if len(overall_scores) > 1 else 0.0

    print("\n📈 Общая статистика:")
    print(f"  Средний балл: {avg_score:.2f}")
    print(f"  Минимум: {min_score}, Максимум: {max_score}")
    print(f"  Стабильность (разброс оценок): {std_dev:.2f}")

    # 9️⃣ Рекомендация
    print("\n💡 Рекомендация:")
    if avg_score < 60:
        print("🔸 Необходимо повторить базовые темы и пройти дополнительные занятия.")
    elif avg_score < 80:
        if std_dev > 15:
            print("⚪ Результаты нестабильны — рекомендуется больше практики и самоконтроля.")
        else:
            print("⚪ Хорошие результаты, но можно улучшить — стоит решать больше задач.")
    else:
        if std_dev < 10:
            print("✅ Отличная работа! Стабильные высокие результаты.")
        else:
            print("✅ Отличные знания, но важно поддерживать стабильность.")

    # 🔟 Визуализация
    plt.figure(figsize=(8, 5))
    labels = [course_names.get(cid, f"Курс {cid}") for cid in avg_per_course.keys()]
    plt.bar(labels, avg_per_course.values(), color="#00DFB3")
    plt.title(f"Средний балл по курсам — {selected_user['name']}")
    plt.xlabel("Курс")
    plt.ylabel("Средний балл")
    plt.xticks(rotation=20, ha="right")
    plt.grid(axis="y", linestyle="--", alpha=0.6)
    plt.tight_layout()
    plt.show()


if __name__ == "__main__":
    analyze_performance()
