from fastapi import FastAPI, Query
from pydantic import BaseModel
import uvicorn
import json
from pathlib import Path

app = FastAPI(title="AI Рекомендательная система")

# ==========================
# ⚙️ НАСТРОЙКИ
# ==========================
DATA_DIR = Path(__file__).parent / "data"

SIMILARITY_THRESHOLD = 0.3
WEIGHTS = {
    "interest": 0.5,
    "level": 0.3,
    "age": 0.2,
}


# ==========================
# 🧩 МОДЕЛИ
# ==========================
class User(BaseModel):
    id: int
    name: str
    age: int
    interest: str
    level: str


# ==========================
# 📚 ЗАГРУЗКА ДАННЫХ
# ==========================
def load_json(filename):
    with open(DATA_DIR / filename, "r", encoding="utf-8") as f:
        return json.load(f)


users = load_json("/home/luka/schoolWithAi/go_recommender/data/users.json")
courses = load_json("/home/luka/schoolWithAi/go_recommender/data/courses.json")
user_courses = load_json("/home/luka/schoolWithAi/go_recommender/data/user_courses.json")


# ==========================
# 🧮 ЛОГИКА РЕКОМЕНДАЦИЙ
# ==========================
def similarity(user_a: dict, user_b: dict) -> float:
    """Вычисляем похожесть между двумя пользователями."""
    score = 0.0

    if user_a["interest"] == user_b["interest"]:
        score += WEIGHTS["interest"]

    if user_a["level"] == user_b["level"]:
        score += WEIGHTS["level"]

    age_diff = abs(user_a["age"] - user_b["age"])
    score += max(0, WEIGHTS["age"] - age_diff * 0.02)

    return round(score, 3)


def recommend(user: User, limit: int = 5):
    """Рекомендует курсы пользователю на основе похожих пользователей."""

    # === 1. Ищем похожих пользователей ===
    similar_users = []
    for other in users:
        if other["id"] == user.id:
            continue

        sim_score = similarity(user.dict(), other)
        if sim_score > SIMILARITY_THRESHOLD:
            similar_users.append((other["id"], sim_score))

    # === 2. Собираем курсы похожих пользователей ===
    course_scores = {}
    for uid, sim in similar_users:
        for uc in user_courses:
            if uc["user_id"] == uid:
                weighted_score = sim * uc["rating"]
                course_scores[uc["course_id"]] = course_scores.get(uc["course_id"], 0) + weighted_score

    # === 3. Исключаем уже пройденные ===
    completed = {uc["course_id"] for uc in user_courses if uc["user_id"] == user.id}

    # === 4. Формируем рекомендации ===
    recommendations = []
    for course_id, score in sorted(course_scores.items(), key=lambda x: x[1], reverse=True):
        if course_id in completed:
            continue

        course = next((c for c in courses if c["id"] == course_id), None)
        if not course:
            continue

        # Финальный бонус
        interest_bonus = 1.2 if course["category"] == user.interest else 1.0
        level_bonus = 1.1 if course["level"] == user.level else 1.0
        final_score = score * interest_bonus * level_bonus

        recommendations.append({**course, "score": round(final_score, 2)})

    return recommendations[:limit]


# ==========================
# 🌐 API
# ==========================
@app.get("/users")
def list_users():
    """Посмотреть всех пользователей"""
    return users


@app.post("/recommend")
def recommend_api(user: User, limit: int = Query(5, ge=1, le=20)):
    """
    Эндпоинт: получить рекомендации для пользователя.
    limit — ограничение количества возвращаемых курсов (по умолчанию 5).
    """
    recs = recommend(user, limit)
    return {
        "user": user.name,
        "limit": limit,
        "recommended": recs
    }


# ==========================
# 🚀 ЗАПУСК
# ==========================
if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
