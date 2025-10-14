from fastapi import FastAPI, HTTPException, Depends
from contextlib import asynccontextmanager
from .database import get_db, db as db_instance
from .services.recommendation_service import RecommendationService
from .schemas.recommendations import (
    RecommendationsRequest,
    RecommendationsListResponse,
    RecommendationResponse,
)


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Управление жизненным циклом приложения"""
    # Startup
    db_instance.connect()
    yield
    # Shutdown
    db_instance.disconnect()


app = FastAPI(
    title="ML Recommendation Service",
    version="1.0.0",
    lifespan=lifespan
)


@app.get("/")
def root():
    """Health check endpoint"""
    return {
        "service": "ML Recommendation Service",
        "status": "running",
        "version": "1.0.0"
    }


@app.post("/recommendations/hybrid", response_model=RecommendationsListResponse)
def get_hybrid_recommendations(
    payload: RecommendationsRequest,
    db = Depends(get_db)
):
    """
    Гибридные рекомендации (объединяет все три метода)
    """
    try:
        service = RecommendationService(db)
        recs = service.get_hybrid_recommendations(payload.student_id, payload.top_n)
        return {"recommendations": [RecommendationResponse(**r) for r in recs]}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/recommendations/collaborative", response_model=RecommendationsListResponse)
def get_collaborative_recommendations(
    payload: RecommendationsRequest,
    db = Depends(get_db)
):
    """
    Коллаборативная фильтрация (по похожим студентам)
    """
    try:
        service = RecommendationService(db)
        recs = service.get_collaborative_recommendations(payload.student_id, payload.top_n)
        return {"recommendations": [RecommendationResponse(**r) for r in recs]}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/recommendations/content-based", response_model=RecommendationsListResponse)
def get_content_based_recommendations(
    payload: RecommendationsRequest,
    db = Depends(get_db)
):
    """
    Контентная фильтрация (по похожим материалам)
    """
    try:
        service = RecommendationService(db)
        recs = service.get_content_based_recommendations(payload.student_id, payload.top_n)
        return {"recommendations": [RecommendationResponse(**r) for r in recs]}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/recommendations/knowledge-based", response_model=RecommendationsListResponse)
def get_knowledge_based_recommendations(
    payload: RecommendationsRequest,
    db = Depends(get_db)
):
    """
    Знаниевая модель (на основе интересов и уровня знаний)
    """
    try:
        service = RecommendationService(db)
        recs = service.get_knowledge_based_recommendations(payload.student_id, payload.top_n)
        return {"recommendations": [RecommendationResponse(**r) for r in recs]}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/health")
def health_check():
    """Проверка состояния сервиса"""
    try:
        # Проверяем подключение к БД
        result = db_instance.execute_one("SELECT 1")
        return {
            "status": "healthy",
            "database": "connected"
        }
    except Exception as e:
        raise HTTPException(
            status_code=503,
            detail=f"Service unhealthy: {str(e)}"
        )