import numpy as np
from typing import List, Dict, Tuple


class CollaborativeFilteringCourses:
    """
    Collaborative Filtering для КУРСОВ (не ресурсов!)
    Рекомендует целые курсы на основе рейтингов и прогресса похожих пользователей
    """
    
    def __init__(self, db):
        self.db = db
        self.user_course_matrix = None
        self.user_ids = []
        self.course_ids = []
    
    def build_matrix(self):
        """
        Построить матрицу User-Course на основе:
        1. Средних оценок ресурсов курса
        2. Прогресса по курсу
        """
        # Получить средние оценки по курсам для каждого студента
        query = """
            SELECT 
                rr.student_id,
                c.id as course_id,
                AVG(rr.rating) as avg_rating
            FROM resource_ratings rr
            JOIN resources r ON rr.resource_id = r.id
            JOIN modules m ON r.module_id = m.id
            JOIN courses c ON m.course_id = c.id
            GROUP BY rr.student_id, c.id
            ORDER BY rr.student_id, c.id
        """
        
        ratings = self.db.execute(query)
        
        if not ratings:
            return None
        
        # Получить уникальные ID
        self.user_ids = sorted(set(r['student_id'] for r in ratings))
        self.course_ids = sorted(set(r['course_id'] for r in ratings))
        
        # Создать матрицу (users x courses)
        n_users = len(self.user_ids)
        n_courses = len(self.course_ids)
        self.user_course_matrix = np.zeros((n_users, n_courses))
        
        # Заполнить матрицу
        user_idx_map = {uid: idx for idx, uid in enumerate(self.user_ids)}
        course_idx_map = {cid: idx for idx, cid in enumerate(self.course_ids)}
        
        for rating in ratings:
            user_idx = user_idx_map[rating['student_id']]
            course_idx = course_idx_map[rating['course_id']]
            self.user_course_matrix[user_idx, course_idx] = rating['avg_rating']
        
        return self.user_course_matrix
    
    def find_similar_users(self, student_id: int, top_n: int = 5) -> List[Tuple[int, float]]:
        """
        Найти похожих пользователей через cosine similarity
        """
        if self.user_course_matrix is None:
            self.build_matrix()
        
        if student_id not in self.user_ids:
            return []
        
        user_idx = self.user_ids.index(student_id)
        user_vector = self.user_course_matrix[user_idx]
        
        # Вычислить cosine similarity
        similarities = []
        for idx, other_user_id in enumerate(self.user_ids):
            if other_user_id == student_id:
                continue
            
            other_vector = self.user_course_matrix[idx]
            
            dot_product = np.dot(user_vector, other_vector)
            norm_user = np.linalg.norm(user_vector)
            norm_other = np.linalg.norm(other_vector)
            
            if norm_user > 0 and norm_other > 0:
                similarity = dot_product / (norm_user * norm_other)
                similarities.append((other_user_id, similarity))
        
        similarities.sort(key=lambda x: x[1], reverse=True)
        return similarities[:top_n]
    
    def recommend(self, student_id: int, top_n: int = 10) -> List[Dict]:
        """
        Рекомендовать КУРСЫ на основе похожих пользователей
        """
        similar_users = self.find_similar_users(student_id, top_n=10)
        
        if not similar_users:
           return self._get_popular_items(top_n)
        
        similar_user_ids = [uid for uid, _ in similar_users]
        
        if not similar_user_ids:
            return []
        
        # Найти курсы, которые нравятся похожим пользователям
        query = """
            WITH similar_user_courses AS (
                SELECT 
                    c.id as course_id,
                    c.title,
                    c.description,
                    c.difficulty_level,
                    c.subject,
                    AVG(rr.rating) as avg_rating,
                    COUNT(DISTINCT rr.student_id) as rating_count
                FROM resource_ratings rr
                JOIN resources r ON rr.resource_id = r.id
                JOIN modules m ON r.module_id = m.id
                JOIN courses c ON m.course_id = c.id
                WHERE rr.student_id IN %s
                  AND rr.rating >= 4
                  AND c.id NOT IN (
                      -- Исключить курсы, которые студент уже полностью завершил
                      SELECT DISTINCT c2.id
                      FROM student_progress sp
                      JOIN resources r2 ON sp.resource_id = r2.id
                      JOIN modules m2 ON r2.module_id = m2.id
                      JOIN courses c2 ON m2.course_id = c2.id
                      WHERE sp.student_id = %s
                        AND sp.status = 'completed'
                      GROUP BY c2.id
                      HAVING COUNT(DISTINCT sp.resource_id) >= (
                          SELECT COUNT(*) 
                          FROM resources r3
                          JOIN modules m3 ON r3.module_id = m3.id
                          WHERE m3.course_id = c2.id
                      ) * 0.8  -- Завершено 80%+ курса
                  )
                GROUP BY c.id, c.title, c.description, c.difficulty_level, c.subject
            )
            SELECT *
            FROM similar_user_courses
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
                'course_id': rec['course_id'],
                'title': rec['title'],
                'description': rec['description'],
                'difficulty_level': rec['difficulty_level'],
                'subject': rec['subject'],
                'score': float(rec['avg_rating']) / 5.0,
                'algorithm': 'collaborative',
                'reason': f"Похожим ученикам понравилось (оценка {rec['avg_rating']:.1f})"
            })
        
        return result