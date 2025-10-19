import numpy as np
from typing import List, Dict, Tuple


class CollaborativeFiltering:
    """
    Collaborative Filtering - находит похожих пользователей
    и рекомендует то, что им понравилось
    """
    
    def __init__(self, db):
        self.db = db
        self.user_item_matrix = None
        self.user_ids = []
        self.resource_ids = []
    
    def build_matrix(self):
        """Построить матрицу User-Item из ratings"""
        query = """
            SELECT student_id, resource_id, rating
            FROM resource_ratings
            ORDER BY student_id, resource_id
        """
        
        ratings = self.db.execute(query)
        
        if not ratings:
            return None
        
        # Получить уникальные ID
        self.user_ids = sorted(set(r['student_id'] for r in ratings))
        self.resource_ids = sorted(set(r['resource_id'] for r in ratings))
        
        # Создать матрицу (users x resources)
        n_users = len(self.user_ids)
        n_resources = len(self.resource_ids)
        self.user_item_matrix = np.zeros((n_users, n_resources))
        
        # Заполнить матрицу
        user_idx_map = {uid: idx for idx, uid in enumerate(self.user_ids)}
        resource_idx_map = {rid: idx for idx, rid in enumerate(self.resource_ids)}
        
        for rating in ratings:
            user_idx = user_idx_map[rating['student_id']]
            resource_idx = resource_idx_map[rating['resource_id']]
            self.user_item_matrix[user_idx, resource_idx] = rating['rating']
        
        return self.user_item_matrix
    
    def find_similar_users(self, student_id: int, top_n: int = 5) -> List[Tuple[int, float]]:
        """
        Найти похожих пользователей через cosine similarity
        Returns: [(user_id, similarity_score), ...]
        """
        if self.user_item_matrix is None:
            self.build_matrix()
        
        if student_id not in self.user_ids:
            return []
        
        user_idx = self.user_ids.index(student_id)
        user_vector = self.user_item_matrix[user_idx]
        
        # Вычислить cosine similarity со всеми пользователями
        similarities = []
        for idx, other_user_id in enumerate(self.user_ids):
            if other_user_id == student_id:
                continue
            
            other_vector = self.user_item_matrix[idx]
            
            # Cosine similarity
            dot_product = np.dot(user_vector, other_vector)
            norm_user = np.linalg.norm(user_vector)
            norm_other = np.linalg.norm(other_vector)
            
            if norm_user > 0 and norm_other > 0:
                similarity = dot_product / (norm_user * norm_other)
                similarities.append((other_user_id, similarity))
        
        # Сортировать по similarity
        similarities.sort(key=lambda x: x[1], reverse=True)
        
        return similarities[:top_n]
    
    def recommend(self, student_id: int, top_n: int = 10) -> List[Dict]:
        """
        Рекомендовать ресурсы на основе похожих пользователей
        """
        # Найти похожих пользователей
        similar_users = self.find_similar_users(student_id, top_n=10)
        
        # ДОБАВЛЕНО: Проверка на пустой список
        if not similar_users:
            print(f"⚠️  No similar users found for student {student_id}")
            return []
        
        similar_user_ids = [uid for uid, _ in similar_users]
        
        # УЛУЧШЕНО: Добавлена проверка перед запросом
        if not similar_user_ids:
            return []
        
        query = """
            SELECT 
                rr.resource_id,
                r.title,
                AVG(rr.rating) as avg_rating,
                COUNT(*) as rating_count,
                ARRAY_AGG(rr.student_id) as rated_by
            FROM resource_ratings rr
            JOIN resources r ON rr.resource_id = r.id
            WHERE rr.student_id IN %s
              AND rr.rating >= 4
              AND rr.resource_id NOT IN (
                  SELECT resource_id 
                  FROM student_progress 
                  WHERE student_id = %s 
                    AND status = 'completed'
              )
            GROUP BY rr.resource_id, r.title
            ORDER BY avg_rating DESC, rating_count DESC
            LIMIT %s
        """
        
        recommendations = self.db.execute(
            query, 
            (tuple(similar_user_ids), student_id, top_n)
        )
        
        result = []
        for rec in recommendations:
            result.append({
                'resource_id': rec['resource_id'],
                'title': rec['title'],
                'score': float(rec['avg_rating']) / 5.0,
                'algorithm': 'collaborative',
                'reason': f"Похожим ученикам понравилось (оценка {rec['avg_rating']:.1f})"
            })
        
        return result