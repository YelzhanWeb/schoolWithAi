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

    # --- COLLABORATIVE ---
    def get_collaborative_recommendations(self, student_id: int, top_n: int = 10):
        recs = self.collaborative.recommend(student_id, top_n)
        return recs

    # --- CONTENT-BASED ---
    def get_content_based_recommendations(self, student_id: int, top_n: int = 10):
        recs = self.content_based.recommend(student_id, top_n)
        return recs

    # --- KNOWLEDGE-BASED ---
    def get_knowledge_based_recommendations(self, student_id: int, top_n: int = 10):
        recs = self.knowledge_based.recommend(student_id, top_n)
        return recs
