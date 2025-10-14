from .collaborative import CollaborativeFiltering
from .content_based import ContentBasedFiltering
from .knowledge_based import KnowledgeBasedFiltering
from .hybrid import HybridRecommender

__all__ = [
    'CollaborativeFiltering',
    'ContentBasedFiltering',
    'KnowledgeBasedFiltering',
    'HybridRecommender'
]