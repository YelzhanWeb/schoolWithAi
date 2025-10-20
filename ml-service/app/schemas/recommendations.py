from pydantic import BaseModel
from typing import List, Optional


class RecommendationsRequest(BaseModel):
    student_id: int
    top_n: Optional[int] = 10


class RecommendationResponse(BaseModel):
    course_id: int
    title: str
    score: float
    algorithm: str
    reason: str
    details: Optional[dict] = None


class RecommendationsListResponse(BaseModel):
    recommendations: List[RecommendationResponse]
