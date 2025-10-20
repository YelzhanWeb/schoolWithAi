from ..models.hybrid import HybridRecommender
from ..models.collaborative import CollaborativeFiltering
from ..models.content_based import ContentBasedFiltering
from ..models.knowledge_based import KnowledgeBasedFiltering


class RecommendationService:
    def __init__(self, db):
        self.db = db
        self.hybrid = HybridRecommender(db)
        self.collaborative = CollaborativeFiltering(db)
        self.content_based = ContentBasedFiltering(db)
        self.knowledge_based = KnowledgeBasedFiltering(db)

    # --- HYBRID ---
    def get_hybrid_recommendations(self, student_id: int, top_n: int = 10):
        recs = self.hybrid.recommend(student_id, top_n)
        self.hybrid.save_recommendations(student_id, recs)
        return recs