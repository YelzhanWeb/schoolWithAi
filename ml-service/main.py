import os
from typing import List
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import pandas as pd
from sqlalchemy import create_engine

DATABASE_URL = "postgresql://admin:admin123@postgres:5432/education_platform"
engine = create_engine(DATABASE_URL)

app = FastAPI(title="AI Recommendation Service")

# DTO для ответа (только ID курсов)
class RecommendationResponse(BaseModel):
    user_id: str
    recommended_course_ids: List[str]

@app.get("/recommend/{user_id}", response_model=RecommendationResponse)
def get_recommendations(user_id: str):
    """
    Генерирует рекомендации на основе интересов и уровня.
    """
    try:
        # 2. Загружаем данные из БД в Pandas DataFrame
        # Читаем интересы студента
        query_interests = f"SELECT subject_id FROM student_interests WHERE user_id = '{user_id}'"
        user_interests = pd.read_sql(query_interests, engine)['subject_id'].tolist()

        # Читаем профиль (уровень)
        query_profile = f"SELECT grade, level FROM student_profiles WHERE user_id = '{user_id}'"
        df_profile = pd.read_sql(query_profile, engine)
        
        if df_profile.empty:
             # Если профиля нет, считаем новичком
            user_grade = 1
        else:
            user_grade = df_profile.iloc[0]['grade']

        # Читаем историю (чтобы не советовать просмотренное)
        query_logs = f"SELECT course_id FROM user_activity_logs WHERE user_id = '{user_id}'"
        viewed_courses = pd.read_sql(query_logs, engine)['course_id'].tolist()

        # Читаем ВСЕ опубликованные курсы
        query_courses = "SELECT id, subject_id, difficulty_level, title FROM courses WHERE is_published = true"
        df_courses = pd.read_sql(query_courses, engine)

        # ---------------------------------------------------------
        # 3. ЛОГИКА ИИ (Здесь магия)
        # ---------------------------------------------------------

        # Шаг 0: Сразу создаем колонку Score
        df_courses['score'] = 0

        # Шаг 1: Фильтрация просмотренных
        # Убираем курсы, ID которых есть в viewed_courses
        df_courses = df_courses[~df_courses['id'].isin(viewed_courses)].copy()

        # Шаг 2: Начисляем баллы за Интересы
        # Если предмет курса совпадает с интересом юзера -> +10 баллов
        if user_interests:
            df_courses.loc[df_courses['subject_id'].isin(user_interests), 'score'] += 10

        # Шаг 3: Начисляем баллы за Сложность (Personalization)
        # Если сложность курса близка к классу ученика (примерная логика)
        # Допустим: курс 1-2 сложность для 1-4 класса, 3-5 для старших.
        # Это упрощенная формула: чем ближе сложность к уровню, тем выше балл
        # (в реальности формула будет сложнее)
        df_courses['score'] += 5  # Базовый бонус всем оставшимся

        # Шаг 4: Ранжирование
        # Сортируем по очкам (score) сверху вниз
        recommendations = df_courses.sort_values(by='score', ascending=False).head(5)

        # 4. Возвращаем результат
        return {
            "user_id": user_id,
            "recommended_course_ids": recommendations['id'].tolist()
        }

    except Exception as e:
        print(f"Error: {e}")
        raise HTTPException(status_code=500, detail="AI Service Error")

# Запуск: uvicorn main:app --reload --port 8000