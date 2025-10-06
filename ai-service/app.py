from fastapi import FastAPI
from pydantic import BaseModel
import uvicorn

app = FastAPI(title="AI Рекомендательная система")

# ==========================
# 🧩 МОДЕЛИ ДАННЫХ
# ==========================
class User(BaseModel):
    id: int
    name: str
    age: int
    interest: str
    level: str


# ==========================
# 📚 ПРИМЕРНЫЕ ДАННЫЕ
# ==========================
courses = [
    {"id": 1, "name": "Python Start", "category": "Программирование", "level": "Начальный"},
    {"id": 2, "name": "Go Advanced", "category": "Программирование", "level": "Продвинутый"},
    {"id": 3, "name": "Canva Design", "category": "Дизайн", "level": "Средний"},
    {"id": 4, "name": "Lego Robots", "category": "Робототехника", "level": "Начальный"},
    {"id": 5, "name": "AI for Kids", "category": "Программирование", "level": "Средний"},
]

user_courses = [
    {"user_id": 1, "course_id": 1, "rating": 5},
    {"user_id": 1, "course_id": 3, "rating": 2},
    {"user_id": 2, "course_id": 3, "rating": 5},
    {"user_id": 3, "course_id": 4, "rating": 4},
    {"user_id": 4, "course_id": 1, "rating": 4},
    {"user_id": 4, "course_id": 2, "rating": 5},
]

users = [
    {"id": 1, "name": "Али", "age": 13, "interest": "Программирование", "level": "Начальный"},
    {"id": 2, "name": "Дана", "age": 14, "interest": "Дизайн", "level": "Средний"},
    {"id": 3, "name": "Тимур", "age": 12, "interest": "Робототехника", "level": "Начальный"},
    {"id": 4, "name": "Айбек", "age": 15, "interest": "Программирование", "level": "Средний"},
]


# ==========================
# 🧮 ФУНКЦИИ РЕКОМЕНДАЦИИ
# ==========================

def similarity(user_a: dict, user_b: dict) -> float:
    """
    Вычисляем похожесть между двумя пользователями.
    Основано на совпадении интересов, уровня и разнице в возрасте.
    """
    score = 0.0

    if user_a["interest"] == user_b["interest"]:
        score += 0.5

    if user_a["level"] == user_b["level"]:
        score += 0.3

    age_diff = abs(user_a["age"] - user_b["age"])
    score += max(0, 0.2 - age_diff * 0.02)

    return score


def recommend(user: User):
    """
    Рекомендует курсы пользователю на основе похожих пользователей.
    Использует простую коллаборативную фильтрацию.
    """
    # === 1. Ищем похожих пользователей ===
    similar_users = []
    for other in users:
        if other["id"] == user.id:
            continue

        sim_score = similarity(user.dict(), other)
        if sim_score > 0.3:
            similar_users.append((other["id"], sim_score))

    # === 2. Собираем курсы, которые они оценили ===
    course_scores = {}
    for uid, sim in similar_users:
        for uc in user_courses:
            if uc["user_id"] == uid:
                weighted_score = sim * uc["rating"]
                course_scores[uc["course_id"]] = course_scores.get(uc["course_id"], 0) + weighted_score

    # === 3. Исключаем уже пройденные курсы пользователя ===
    completed_courses = {uc["course_id"] for uc in user_courses if uc["user_id"] == user.id}

    # === 4. Формируем итоговые рекомендации ===
    recommendations = []
    sorted_courses = sorted(course_scores.items(), key=lambda x: x[1], reverse=True)

    for course_id, score in sorted_courses:
        if course_id in completed_courses:
            continue

        # находим курс по id
        course = next((c for c in courses if c["id"] == course_id), None)
        if not course:
            continue

        # финальный скоринг с учетом интереса и уровня
        interest_bonus = 1.2 if course["category"] == user.interest else 1.0
        level_bonus = 1.1 if course["level"] == user.level else 1.0
        final_score = score * interest_bonus * level_bonus

        course_with_score = {
            **course,
            "score": round(final_score, 2)
        }
        recommendations.append(course_with_score)

    return recommendations[:3]


# ==========================
# 🌐 API ЭНДПОИНТЫ
# ==========================
@app.post("/recommend")
def recommend_api(user: User):
    """Эндпоинт: получить рекомендации для пользователя"""
    recs = recommend(user)
    return {
        "user": user.name,
        "recommended": recs
    }


# ==========================
# 🚀 ЗАПУСК
# ==========================
if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
