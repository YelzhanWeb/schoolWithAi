from fastapi import FastAPI, HTTPException
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
from typing import List, Optional
import uvicorn

from .database import db
from .models.hybrid import HybridRecommender
from .models.knowledge_based import KnowledgeBasedFiltering

# FastAPI app
app = FastAPI(
    title="OqysAI ML Service",
    description="ML-powered recommendation system",
    version="1.0.0"
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
recommender = None
knowledge_filter = None


# Pydantic schemas
class RecommendationResponse(BaseModel):
    resource_id: int
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


# Events
@app.on_event("startup")
async def startup_event():
    """Initialize ML models on startup"""
    global recommender, knowledge_filter
    
    print("🚀 Starting ML Service...")
    
    # Connect to database
    db.connect()
    
    # Initialize ML models
    recommender = HybridRecommender(db)
    knowledge_filter = KnowledgeBasedFiltering(db)
    
    # Build collaborative filtering matrix
    print("Building collaborative filtering matrix...")
    recommender.collaborative.build_matrix()
    
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
        "version": "1.0.0"
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
            "ml_models": "loaded"
        }
    except Exception as e:
        raise HTTPException(status_code=503, detail=f"Service unhealthy: {str(e)}")


@app.post("/recommendations/hybrid", response_model=List[RecommendationResponse])
async def get_hybrid_recommendations(request: RecommendationRequest):
    """
    Get hybrid recommendations (combines all three methods)
    """
    try:
        recs = recommender.recommend(request.student_id, request.top_n)
        
        # Save to database
        recommender.save_recommendations(request.student_id, recs)
        
        return recs
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/recommendations/collaborative", response_model=List[RecommendationResponse])
async def get_collaborative_recommendations(request: RecommendationRequest):
    """
    Get collaborative filtering recommendations
    """
    try:
        recs = recommender.collaborative.recommend(request.student_id, request.top_n)
        return recs
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/recommendations/content-based", response_model=List[RecommendationResponse])
async def get_content_based_recommendations(request: RecommendationRequest):
    """
    Get content-based recommendations
    """
    try:
        recs = recommender.content_based.recommend(request.student_id, request.top_n)
        return recs
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/recommendations/knowledge-based", response_model=List[RecommendationResponse])
async def get_knowledge_based_recommendations(request: RecommendationRequest):
    """
    Get knowledge-based recommendations
    """
    try:
        recs = recommender.knowledge_based.recommend(request.student_id, request.top_n)
        return recs
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


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

# ml-service/app/main.py - ДОБАВИТЬ:

@app.post("/recommendations/courses/hybrid", response_model=List[RecommendationResponse])
async def get_hybrid_course_recommendations(request: RecommendationRequest):
    """
    Get hybrid course recommendations
    """
    try:
        from .models.hybrid_courses import HybridRecommenderCourses
        
        courses_recommender = HybridRecommenderCourses(db)
        recs = courses_recommender.recommend(request.student_id, request.top_n)
        
        # Конвертировать course_id -> resource_id для совместимости
        # (или создать отдельную схему CourseRecommendationResponse)
        result = []
        for rec in recs:
            result.append({
                "resource_id": rec['course_id'],  # Используем course_id
                "title": rec['title'],
                "score": rec['score'],
                "algorithm": rec['algorithm'],
                "reason": rec['reason'],
                "details": rec.get('details')
            })
        
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
    
# Run with: uvicorn app.main:app --reload --host 0.0.0.0 --port 5000
if __name__ == "__main__":
    uvicorn.run("app.main:app", host="0.0.0.0", port=5000, reload=True)