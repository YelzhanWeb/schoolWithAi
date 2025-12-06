import os
import math
from typing import List

import numpy as np
import pandas as pd
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from sqlalchemy import create_engine

# -----------------------------
# Подключение к БД
# -----------------------------
DATABASE_URL = os.getenv(
    "DATABASE_URL",
    "postgresql://admin:admin123@localhost:5432/education_platform",
)
engine = create_engine(DATABASE_URL)

app = FastAPI(title="AI Recommendation Service")


# -----------------------------
# DTO
# -----------------------------
class RecommendationResponse(BaseModel):
    user_id: str
    recommended_course_ids: List[str]


# -----------------------------
# Вспомогательные функции
# -----------------------------
def _get_user_profile(user_id: str):
    df = pd.read_sql(
        """
        SELECT grade, level
        FROM student_profiles
        WHERE user_id = %(user_id)s
        """,
        engine,
        params={"user_id": user_id},
    )
    if df.empty:
        return {"grade": 1, "level": 1}
    row = df.iloc[0]
    return {
        "grade": int(row.get("grade", 1) or 1),
        "level": int(row.get("level", 1) or 1),
    }


def _get_user_interests(user_id: str) -> List[str]:
    df = pd.read_sql(
        """
        SELECT subject_id
        FROM student_interests
        WHERE user_id = %(user_id)s
        """,
        engine,
        params={"user_id": user_id},
    )
    return df["subject_id"].tolist() if not df.empty else []


def _get_viewed_courses(user_id: str) -> List[str]:
    df = pd.read_sql(
        """
        SELECT DISTINCT course_id
        FROM user_activity_logs
        WHERE user_id = %(user_id)s
          AND course_id IS NOT NULL
        """,
        engine,
        params={"user_id": user_id},
    )
    return df["course_id"].tolist() if not df.empty else []


def _get_all_courses_with_stats(user_grade: int) -> pd.DataFrame:
    """
    Берём все опубликованные курсы и считаем:
    - popularity: сколько раз курс смотрели вообще
    - similar_popularity: сколько раз курс смотрели ученики того же класса (grade)
    Это даёт нам лёгкий "коллаборативный фильтр".
    """
    query = """
    SELECT
        c.id,
        c.subject_id,
        c.difficulty_level,
        c.title,
        c.created_at,
        COALESCE(COUNT(ual.id), 0) AS popularity,
        COALESCE(
            SUM(CASE WHEN sp.grade = %(user_grade)s THEN 1 ELSE 0 END),
            0
        ) AS similar_popularity
    FROM courses c
    LEFT JOIN user_activity_logs ual ON ual.course_id = c.id
    LEFT JOIN student_profiles sp ON sp.user_id = ual.user_id
    WHERE c.is_published = TRUE
    GROUP BY c.id, c.subject_id, c.difficulty_level, c.title, c.created_at
    """
    return pd.read_sql(query, engine, params={"user_grade": user_grade})


# -----------------------------
# Нормализация фич
# -----------------------------
def _normalize_popularity(series: pd.Series) -> pd.Series:
    """
    Логарифмическая нормализация популярности в диапазон (0..1].
    """
    series = series.astype(float)
    max_val = series.max()
    if max_val <= 0:
        return pd.Series(0.0, index=series.index)
    return np.log1p(series) / math.log1p(max_val)


def _difficulty_score(difficulty_level: float, user_grade: float) -> float:
    """
    Нормированная оценка сложности (0..1):
    1.0 — идеальная сложность, дальше плавно падает.
    """
    if pd.isna(difficulty_level):
        return 0.3  # почти базовый шанс, если сложность не задана
    diff = abs(float(difficulty_level) - float(user_grade))
    # считаем, что оптимальная зона +/- 1 класс, дальше резко хуже
    return max(0.0, 1.0 - (diff / 3.0))


def _recency_score(created_at: pd.Timestamp, now: pd.Timestamp) -> float:
    """
    Экспоненциальное затухание по возрасту курса.
    Новые курсы ближе к 1, старые — ближе к 0.
    """
    if pd.isna(created_at):
        return 0.5
    age_days = max(0, (now - created_at).days)
    # 30 дней — характерный масштаб
    return math.exp(-age_days / 30.0)


# -----------------------------
# Основной эндпоинт
# -----------------------------
@app.get("/recommend/{user_id}", response_model=RecommendationResponse)
def get_recommendations(user_id: str):
    """
    Алгоритм рекомендаций:
    1. Собираем профиль пользователя (grade, level, интересы, просмотренные курсы).
    2. Для всех опубликованных курсов считаем признаки:
        - subject_match: совпадение по предмету
        - difficulty_score: близость сложности к grade
        - popularity_score: общая популярность
        - similar_popularity_score: популярность среди учеников того же класса
        - recency_score: свежесть курса
    3. Складываем признаки с весами (линейная модель).
    4. Исключаем уже просмотренные курсы.
    5. Делаем лёгкую диверсификацию: большинство курсов по интересам,
       + иногда 1–2 «новых» направления.
    """
    try:
        # 1. Данные о пользователе
        profile = _get_user_profile(user_id)
        user_grade = profile["grade"]
        user_interests = _get_user_interests(user_id)
        viewed_courses = _get_viewed_courses(user_id)

        # 2. Все курсы со статистикой
        df_courses = _get_all_courses_with_stats(user_grade=user_grade)
        if df_courses.empty:
            return {"user_id": user_id, "recommended_course_ids": []}

        # 3. Убираем уже просмотренные
        if viewed_courses:
            df_courses = df_courses[~df_courses["id"].isin(viewed_courses)].copy()
        else:
            df_courses = df_courses.copy()

        if df_courses.empty:
            return {"user_id": user_id, "recommended_course_ids": []}

        now = pd.Timestamp.utcnow()
        df_courses["created_at"] = pd.to_datetime(df_courses["created_at"])

        # ---------- Фичи ----------
        # Совпадение по интересам (0 или 1)
        if user_interests:
            df_courses["subject_match"] = df_courses["subject_id"].isin(user_interests).astype(float)
        else:
            df_courses["subject_match"] = 0.0

        # Сложность (0..1)
        df_courses["difficulty_score"] = df_courses["difficulty_level"].apply(
            lambda d: _difficulty_score(d, user_grade)
        )

        # Популярность (0..1, лог-нормализация)
        df_courses["popularity_score"] = _normalize_popularity(df_courses["popularity"])

        # Популярность среди «похожих» (одноклассников) (0..1)
        df_courses["similar_popularity_score"] = _normalize_popularity(
            df_courses["similar_popularity"]
        )

        # Свежесть (0..1, экспонента)
        df_courses["recency_score"] = df_courses["created_at"].apply(
            lambda dt: _recency_score(dt, now)
        )

        # ---------- Итоговый скор ----------
        # Веса можно вынести в конфиг / .env, сейчас подобраны эмпирически.
        W_SUBJECT = 0.30
        W_DIFFICULTY = 0.25
        W_POP = 0.15
        W_SIMILAR_POP = 0.20
        W_RECENCY = 0.10

        df_courses["score"] = (
            W_SUBJECT * df_courses["subject_match"]
            + W_DIFFICULTY * df_courses["difficulty_score"]
            + W_POP * df_courses["popularity_score"]
            + W_SIMILAR_POP * df_courses["similar_popularity_score"]
            + W_RECENCY * df_courses["recency_score"]
        )

        # ---------- Ранжирование + диверсификация ----------
        top_n = int(os.getenv("RECOMMENDATIONS_TOP_N", 5))
        df_sorted = df_courses.sort_values(by="score", ascending=False)

        # Сначала набираем курсы по интересам
        df_main = df_sorted[df_sorted["subject_match"] == 1.0].head(top_n)

        # Если мало или хотим добавить «новые» предметы — добираем
        if len(df_main) < top_n:
            remaining = top_n - len(df_main)
            df_explore = df_sorted[df_sorted["subject_match"] == 0.0].head(remaining)
            df_final = pd.concat([df_main, df_explore], ignore_index=True)
        else:
            df_final = df_main

        # На всякий случай ещё раз упорядочим по score
        df_final = df_final.sort_values(by="score", ascending=False).head(top_n)

        return {
            "user_id": user_id,
            "recommended_course_ids": df_final["id"].tolist(),
        }

    except Exception as e:
        print(f"Error in get_recommendations: {e}", flush=True)
        raise HTTPException(status_code=500, detail="AI Service Error")


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=5000,
        reload=True,
    )
