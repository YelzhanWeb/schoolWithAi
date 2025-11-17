from fastapi import FastAPI, HTTPException, Query
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import List, Optional
from enum import Enum
import uvicorn

from .database import db
from .models.hybrid import HybridRecommender
from .models.collaborative import CollaborativeFiltering
from .models.content_based import ContentBasedFiltering
from .models.knowledge_based import KnowledgeBasedFiltering

# FastAPI app
app = FastAPI(
    title="OqysAI ML Service",
    description="ML-powered recommendation system with multiple algorithms",
    version="2.0.0"
)

# CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Global ML models
collaborative_filter = None
content_filter = None
knowledge_filter = None
hybrid_recommender = None


# Enums
class AlgorithmType(str, Enum):
    COLLABORATIVE = "collaborative"
    CONTENT_BASED = "content_based"
    KNOWLEDGE_BASED = "knowledge_based"
    HYBRID = "hybrid"


# Pydantic schemas
class RecommendationResponse(BaseModel):
    course_id: int
    title: str
    score: float
    algorithm: str
    reason: str
    details: Optional[dict] = None


class RecommendationRequest(BaseModel):
    student_id: int
    top_n: Optional[int] = 10


class SkillUpdateRequest(BaseModel):
    student_id: int
    skill_name: str
    test_score: float  # 0.0 - 1.0


class SkillUpdateResponse(BaseModel):
    student_id: int
    skill_name: str
    new_level: float
    message: str


class AlgorithmInfo(BaseModel):
    name: str
    type: str
    description: str
    use_cases: List[str]
    strengths: List[str]
    limitations: List[str]


# Events
@app.on_event("startup")
async def startup_event():
    """Initialize ML models on startup"""
    global collaborative_filter, content_filter, knowledge_filter, hybrid_recommender
    
    print("🚀 Starting ML Service...")
    
    # Connect to database
    db.connect()
    
    # Initialize all ML models
    print("📦 Initializing ML models...")
    collaborative_filter = CollaborativeFiltering(db)
    content_filter = ContentBasedFiltering(db)
    knowledge_filter = KnowledgeBasedFiltering(db)
    hybrid_recommender = HybridRecommender(db)
    
    print("✅ ML Service ready!")


@app.on_event("shutdown")
async def shutdown_event():
    """Cleanup on shutdown"""
    print("Shutting down ML Service...")
    db.disconnect()


# API Endpoints

@app.get("/")
async def root():
    """Health check"""
    return {
        "status": "ok",
        "service": "OqysAI ML Service",
        "version": "2.0.0",
        "algorithms": ["collaborative", "content_based", "knowledge_based", "hybrid"]
    }


@app.get("/health")
async def health_check():
    """Detailed health check"""
    try:
        # Test database connection
        result = db.execute_one("SELECT 1")
        return {
            "status": "healthy",
            "database": "connected",
            "ml_models": {
                "collaborative": collaborative_filter is not None,
                "content_based": content_filter is not None,
                "knowledge_based": knowledge_filter is not None,
                "hybrid": hybrid_recommender is not None
            }
        }
    except Exception as e:
        raise HTTPException(status_code=503, detail=f"Service unhealthy: {str(e)}")


@app.get("/algorithms", response_model=List[AlgorithmInfo])
async def get_algorithms_info():
    """
    Get information about available recommendation algorithms
    """
    return [
        AlgorithmInfo(
            name="Collaborative Filtering",
            type="collaborative",
            description="Рекомендует курсы на основе предпочтений похожих пользователей",
            use_cases=[
                "Когда у студента есть история оценок",
                "Для популярных курсов с множеством отзывов",
                "Для обнаружения неожиданных интересов"
            ],
            strengths=[
                "Не требует знания содержимого курсов",
                "Обнаруживает скрытые паттерны",
                "Хорошо работает для опытных пользователей"
            ],
            limitations=[
                "Cold start для новых пользователей",
                "Требует достаточно данных о рейтингах",
                "Может создавать 'filter bubble'"
            ]
        ),
        AlgorithmInfo(
            name="Content-Based Filtering",
            type="content_based",
            description="Рекомендует курсы на основе интересов и предпочтений студента",
            use_cases=[
                "Когда известны интересы студента",
                "Для персонализированного обучения",
                "Для соответствия возрастной группе"
            ],
            strengths=[
                "Работает для новых пользователей с указанными интересами",
                "Прозрачные рекомендации",
                "Не требует данных о других пользователях"
            ],
            limitations=[
                "Ограничен известными интересами",
                "Может быть слишком узким",
                "Не обнаруживает новые интересы"
            ]
        ),
        AlgorithmInfo(
            name="Knowledge-Based Filtering",
            type="knowledge_based",
            description="Рекомендует курсы для развития конкретных навыков",
            use_cases=[
                "Для целенаправленного обучения",
                "Когда нужно улучшить слабые навыки",
                "Для адаптивного образовательного пути"
            ],
            strengths=[
                "Фокус на развитии навыков",
                "Учитывает текущий уровень студента",
                "Адаптивная сложность"
            ],
            limitations=[
                "Требует профиль навыков студента",
                "Может игнорировать интересы",
                "Фокус только на gap-filling"
            ]
        ),
        AlgorithmInfo(
            name="Hybrid Recommender",
            type="hybrid",
            description="Комбинирует все алгоритмы с адаптивными весами для лучших результатов",
            use_cases=[
                "Для сбалансированных рекомендаций",
                "Когда нужна максимальная точность",
                "Универсальное решение для всех случаев"
            ],
            strengths=[
                "Объединяет преимущества всех алгоритмов",
                "Адаптивные веса",
                "Консенсус-бонус для согласованных рекомендаций",
                "Diversity boost"
            ],
            limitations=[
                "Более сложный для объяснения",
                "Требует больше вычислений"
            ]
        )
    ]


# === RECOMMENDATION ENDPOINTS ===

@app.post("/recommendations/collaborative", response_model=List[RecommendationResponse])
async def get_collaborative_recommendations(request: RecommendationRequest):
    """
    Get recommendations using COLLABORATIVE FILTERING only
    
    Лучше всего работает для:
    - Пользователей с историей оценок
    - Обнаружения популярных курсов среди похожих студентов
    """
    try:
        recs = collaborative_filter.recommend(request.student_id, request.top_n)
        
        if not recs:
            raise HTTPException(
                status_code=404, 
                detail="No collaborative recommendations found. User might be new or have insufficient rating history."
            )
        
        return recs
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/recommendations/content-based", response_model=List[RecommendationResponse])
async def get_content_based_recommendations(request: RecommendationRequest):
    """
    Get recommendations using CONTENT-BASED FILTERING only
    
    Лучше всего работает для:
    - Персонализации по интересам
    - Новых пользователей с заполненным профилем
    - Соответствия возрастной группе
    """
    try:
        recs = content_filter.recommend(request.student_id, request.top_n)
        
        if not recs:
            raise HTTPException(
                status_code=404,
                detail="No content-based recommendations found. User might not have interests specified."
            )
        
        return recs
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/recommendations/knowledge-based", response_model=List[RecommendationResponse])
async def get_knowledge_based_recommendations(request: RecommendationRequest):
    """
    Get recommendations using KNOWLEDGE-BASED FILTERING only
    
    Лучше всего работает для:
    - Целенаправленного развития навыков
    - Улучшения слабых областей
    - Адаптивного обучения
    """
    try:
        recs = knowledge_filter.recommend(request.student_id, request.top_n)
        
        if not recs:
            raise HTTPException(
                status_code=404,
                detail="No knowledge-based recommendations found. User might not have skill profile."
            )
        
        return recs
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/recommendations/hybrid", response_model=List[RecommendationResponse])
async def get_hybrid_recommendations(request: RecommendationRequest):
    """
    Get recommendations using HYBRID SYSTEM (все алгоритмы вместе)
    
    Рекомендуется по умолчанию:
    - Комбинирует все подходы
    - Адаптивные веса
    - Консенсус-бонус
    - Diversity boost
    """
    try:        
        recs = hybrid_recommender.recommend(request.student_id, request.top_n)
        
        if not recs:
            raise HTTPException(
                status_code=404,
                detail="No recommendations found for this student."
            )
        
        return recs
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/recommendations/{student_id}", response_model=List[RecommendationResponse])
async def get_recommendations(
    student_id: int,
    algorithm: AlgorithmType = Query(AlgorithmType.HYBRID, description="Recommendation algorithm to use"),
    top_n: int = Query(10, ge=1, le=50, description="Number of recommendations")
):
    """
    Universal endpoint: Get recommendations using specified algorithm
    
    Parameters:
    - student_id: ID студента
    - algorithm: Тип алгоритма (collaborative, content_based, knowledge_based, hybrid)
    - top_n: Количество рекомендаций (1-50)
    """
    try:
        if algorithm == AlgorithmType.COLLABORATIVE:
            recs = collaborative_filter.recommend(student_id, top_n)
        elif algorithm == AlgorithmType.CONTENT_BASED:
            recs = content_filter.recommend(student_id, top_n)
        elif algorithm == AlgorithmType.KNOWLEDGE_BASED:
            recs = knowledge_filter.recommend(student_id, top_n)
        else:  # HYBRID
            recs = hybrid_recommender.recommend(student_id, top_n)
        
        if not recs:
            raise HTTPException(
                status_code=404,
                detail=f"No {algorithm.value} recommendations found for student {student_id}"
            )
        
        return recs
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


# === COMPARISON ENDPOINT ===

@app.get("/recommendations/{student_id}/compare")
async def compare_algorithms(
    student_id: int,
    top_n: int = Query(5, ge=1, le=20, description="Number of recommendations per algorithm")
):
    """
    Compare all algorithms side-by-side for analysis
    
    Возвращает рекомендации от всех 4 алгоритмов для сравнения
    """
    try:
        return {
            "student_id": student_id,
            "collaborative": collaborative_filter.recommend(student_id, top_n),
            "content_based": content_filter.recommend(student_id, top_n),
            "knowledge_based": knowledge_filter.recommend(student_id, top_n),
            "hybrid": hybrid_recommender.recommend(student_id, top_n)
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


# === SKILL MANAGEMENT ===

@app.post("/skills/update", response_model=SkillUpdateResponse)
async def update_skill_level(request: SkillUpdateRequest):
    """
    Update student skill level based on test results
    """
    try:
        new_level = knowledge_filter.update_skill(
            request.student_id,
            request.skill_name,
            request.test_score
        )
        
        return SkillUpdateResponse(
            student_id=request.student_id,
            skill_name=request.skill_name,
            new_level=new_level,
            message=f"Skill '{request.skill_name}' updated to {new_level:.2f}"
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/skills/{student_id}")
async def get_student_skills(student_id: int):
    """
    Get all skills for a student with proficiency levels
    """
    try:
        query = """
            SELECT 
                skill_name,
                proficiency_level,
                updated_at,
                CASE 
                    WHEN proficiency_level < 0.5 THEN 'weak'
                    WHEN proficiency_level < 0.75 THEN 'medium'
                    ELSE 'advanced'
                END as category
            FROM student_skills
            WHERE student_id = %s
            ORDER BY proficiency_level ASC
        """
        
        skills = db.execute(query, (student_id,))
        
        if not skills:
            return {
                "student_id": student_id,
                "skills": [],
                "message": "No skills found for this student"
            }
        
        return {
            "student_id": student_id,
            "skills": [
                {
                    "skill_name": s['skill_name'],
                    "proficiency_level": float(s['proficiency_level']),
                    "category": s['category'],
                    "updated_at": s['updated_at'].isoformat() if s['updated_at'] else None
                }
                for s in skills
            ],
            "weak_skills": [s['skill_name'] for s in skills if s['category'] == 'weak'],
            "medium_skills": [s['skill_name'] for s in skills if s['category'] == 'medium'],
            "advanced_skills": [s['skill_name'] for s in skills if s['category'] == 'advanced']
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


# === ANALYTICS ===

@app.get("/analytics/algorithm-performance/{student_id}")
async def get_algorithm_performance(student_id: int):
    """
    Analyze which algorithm performs best for this student
    Based on historical click-through and completion rates
    """
    try:
        query = """
            SELECT 
                algorithm_type,
                COUNT(*) as total_recommendations,
                SUM(CASE WHEN is_viewed THEN 1 ELSE 0 END) as viewed_count,
                AVG(score) as avg_score
            FROM course_recommendations
            WHERE student_id = %s
              AND created_at > NOW() - INTERVAL '30 days'
            GROUP BY algorithm_type
            ORDER BY viewed_count DESC
        """
        
        stats = db.execute(query, (student_id,))
        
        return {
            "student_id": student_id,
            "period": "last_30_days",
            "algorithm_stats": [
                {
                    "algorithm": s['algorithm_type'],
                    "total_recommendations": s['total_recommendations'],
                    "viewed": s['viewed_count'],
                    "click_through_rate": round(s['viewed_count'] / s['total_recommendations'], 3) if s['total_recommendations'] > 0 else 0,
                    "avg_score": round(float(s['avg_score']), 3)
                }
                for s in stats
            ]
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


# Run with: uvicorn app.main:app --reload --host 0.0.0.0 --port 5000
if __name__ == "__main__":
    uvicorn.run("app.main:app", host="0.0.0.0", port=5000, reload=True)