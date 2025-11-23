import os
from typing import List
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import pandas as pd
from sqlalchemy import create_engine, text

DATABASE_URL = os.getenv(
    "DATABASE_URL",
    "postgresql://admin:admin123@postgres:5432/education_platform"
)
engine = create_engine(DATABASE_URL)

app = FastAPI(title="AI Recommendation Service")


# ========================================
# DTOs
# ========================================

class RecommendationRequest(BaseModel):
    user_id: str
    top_n: int = 5


class RecommendationResponse(BaseModel):
    user_id: str
    recommended_course_ids: List[str]


# ========================================
# API Endpoints
# ========================================

@app.get("/health")
def health_check():
    """Health check endpoint"""
    return {"status": "ok", "service": "ml-recommendation"}


@app.post("/recommendations/hybrid", response_model=RecommendationResponse)
def get_hybrid_recommendations(request: RecommendationRequest):
    """
    Генерирует персонализированные рекомендации курсов.
    
    Алгоритм:
    1. Загружает интересы студента (предметы)
    2. Загружает профиль (класс, уровень)
    3. Фильтрует просмотренные курсы
    4. Начисляет баллы за релевантность:
       - +10 за совпадение предмета с интересами
       - +5 базовый бонус
       - +3 за подходящую сложность
    5. Сортирует по рейтингу и возвращает топ N
    """
    try:
        user_id = request.user_id
        top_n = request.top_n

        # ========================================
        # 1. Загрузка данных из БД (с параметризацией)
        # ========================================
        
        # Интересы студента
        query_interests = text("""
            SELECT subject_id 
            FROM student_interests 
            WHERE user_id = :user_id
        """)
        user_interests = pd.read_sql(
            query_interests, 
            engine, 
            params={"user_id": user_id}
        )['subject_id'].tolist()

        # Профиль студента
        query_profile = text("""
            SELECT grade, level, xp 
            FROM student_profiles 
            WHERE user_id = :user_id
        """)
        df_profile = pd.read_sql(
            query_profile, 
            engine, 
            params={"user_id": user_id}
        )

        if df_profile.empty:
            user_grade = 1
            user_level = 1
        else:
            user_grade = int(df_profile.iloc[0]['grade'])
            user_level = int(df_profile.iloc[0]['level'])

        # История просмотров
        query_logs = text("""
            SELECT DISTINCT course_id 
            FROM user_activity_logs 
            WHERE user_id = :user_id 
            AND course_id IS NOT NULL
        """)
        viewed_courses = pd.read_sql(
            query_logs, 
            engine, 
            params={"user_id": user_id}
        )['course_id'].tolist()

        # Все опубликованные курсы
        query_courses = text("""
            SELECT id, subject_id, difficulty_level, title, tags
            FROM courses 
            WHERE is_published = true
        """)
        df_courses = pd.read_sql(query_courses, engine)

        # ========================================
        # 2. AI SCORING ALGORITHM
        # ========================================

        if df_courses.empty:
            return RecommendationResponse(
                user_id=user_id,
                recommended_course_ids=[]
            )

        # Инициализация score
        df_courses['score'] = 0

        # Фильтрация просмотренных курсов
        if viewed_courses:
            df_courses = df_courses[
                ~df_courses['id'].isin(viewed_courses)
            ].copy()

        if df_courses.empty:
            return RecommendationResponse(
                user_id=user_id,
                recommended_course_ids=[]
            )

        # Скоринг по интересам
        if user_interests:
            df_courses.loc[
                df_courses['subject_id'].isin(user_interests), 
                'score'
            ] += 10

        # Базовый бонус
        df_courses['score'] += 5

        # Скоринг по сложности (чем ближе к уровню, тем лучше)
        # Упрощенная формула: -abs(difficulty - user_level) * 2
        df_courses['difficulty_score'] = -abs(
            df_courses['difficulty_level'] - user_level
        ) * 2
        df_courses['score'] += df_courses['difficulty_score']

        # Бонус за подходящую сложность
        suitable_difficulty = df_courses['difficulty_level'] <= user_level + 1
        df_courses.loc[suitable_difficulty, 'score'] += 3

        # ========================================
        # 3. Ранжирование и возврат
        # ========================================

        recommendations = df_courses.sort_values(
            by='score', 
            ascending=False
        ).head(top_n)

        return RecommendationResponse(
            user_id=user_id,
            recommended_course_ids=recommendations['id'].tolist()
        )

    except Exception as e:
        print(f"❌ Error in ML Service: {e}")
        raise HTTPException(
            status_code=500, 
            detail=f"AI Service Error: {str(e)}"
        )


# ========================================
# Legacy endpoint для обратной совместимости
# ========================================

@app.get("/recommend/{user_id}", response_model=RecommendationResponse)
def get_recommendations_legacy(user_id: str):
    """Legacy endpoint - использует POST метод внутри"""
    return get_hybrid_recommendations(
        RecommendationRequest(user_id=user_id, top_n=5)
    )


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)