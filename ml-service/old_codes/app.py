from fastapi import FastAPI
from pydantic import BaseModel
import uvicorn
import json
import os

app = FastAPI(title="AI Рекомендательная система")

# ==========================
# 🔧 ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ
# ==========================
def load_json(filename: str):
    """Загружает JSON-файл из папки data"""
    path = os.path.join("data", filename)
    with open(path, "r", encoding="utf-8") as f:
        return json.load(f)


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
# 📚 ЗАГРУЗКА ДАННЫХ
# ==========================
users = load_json("/home/luka/schoolWithAi/go_recommender/data/users.json")
courses = load_json("/home/luka/schoolWithAi/go_recommender/data/courses.json")
user_courses = load_json("/home/luka/schoolWithAi/go_recommender/data/user_courses.json")


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
    # 1️⃣ Поиск похожих пользователей
    similar_users = []
    for other in users:
        if other["id"] == user.id:
            continue

        sim_score = similarity(user.dict(), other)
        if sim_score > 0.3:
            similar_users.append((other["id"], sim_score))

    # 2️⃣ Сбор курсов от похожих пользователей
    course_scores = {}
    for uid, sim in similar_users:
        for uc in user_courses:
            if uc["user_id"] == uid:
                weighted_score = sim * uc["rating"]
                course_scores[uc["course_id"]] = course_scores.get(uc["course_id"], 0) + weighted_score

    # 3️⃣ Убираем уже пройденные курсы
    completed = {uc["course_id"] for uc in user_courses if uc["user_id"] == user.id}

    # 4️⃣ Формируем рекомендации
    recommendations = []
    sorted_courses = sorted(course_scores.items(), key=lambda x: x[1], reverse=True)

    for course_id, score in sorted_courses:
        if course_id in completed:
            continue

        course = next((c for c in courses if c["id"] == course_id), None)
        if not course:
            continue

        interest_bonus = 1.2 if course["category"] == user.interest else 1.0
        level_bonus = 1.1 if course["level"] == user.level else 1.0
        final_score = score * interest_bonus * level_bonus

        recommendations.append({
            **course,
            "score": round(final_score, 2)
        })

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
