import numpy as np
from typing import List, Dict, Tuple
from datetime import datetime, timedelta


class CollaborativeFiltering:
    """
    Collaborative Filtering для КУРСОВ (исправленная версия)
    - Centered cosine similarity
    - Обработка cold start
    - Decay factor для старых оценок
    - Правильная нормализация scores
    """
    
    def __init__(self, db):
        self.db = db
        self.user_course_matrix = None
        self.user_means = None  # Средние оценки каждого пользователя
        self.user_ids = []
        self.course_ids = []
        self.decay_days = 180  # Оценки старше 180 дней весят меньше
    
    def build_matrix(self):
        """
        Построить матрицу User-Course с учетом времени оценок
        """
        query = """
            SELECT 
                rr.student_id,
                c.id as course_id,
                AVG(rr.rating) as avg_rating,
                MAX(rr.created_at) as last_rating_date,
                COUNT(*) as rating_count
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
        
        self.user_ids = sorted(set(r['student_id'] for r in ratings))
        self.course_ids = sorted(set(r['course_id'] for r in ratings))
        
        n_users = len(self.user_ids)
        n_courses = len(self.course_ids)
        self.user_course_matrix = np.zeros((n_users, n_courses))
        
        user_idx_map = {uid: idx for idx, uid in enumerate(self.user_ids)}
        course_idx_map = {cid: idx for idx, cid in enumerate(self.course_ids)}
        
        # Заполнить матрицу с учетом decay factor
        now = datetime.now()
        
        for rating in ratings:
            user_idx = user_idx_map[rating['student_id']]
            course_idx = course_idx_map[rating['course_id']]
            
            # Decay factor: старые оценки весят меньше
            days_old = (now - rating['last_rating_date']).days
            decay = np.exp(-days_old / self.decay_days)
            
            weighted_rating = rating['avg_rating'] * decay
            self.user_course_matrix[user_idx, course_idx] = weighted_rating
        
        # Вычислить средние оценки для каждого пользователя (для centered cosine)
        self.user_means = np.zeros(n_users)
        for i in range(n_users):
            rated_items = self.user_course_matrix[i] > 0
            if rated_items.any():
                self.user_means[i] = self.user_course_matrix[i][rated_items].mean()
        
        return self.user_course_matrix
    
    def find_similar_users(self, student_id: int, top_n: int = 10) -> List[Tuple[int, float]]:
        """
        Найти похожих пользователей через centered cosine similarity
        """
        if self.user_course_matrix is None:
            self.build_matrix()
        
        if student_id not in self.user_ids:
            return []
        
        user_idx = self.user_ids.index(student_id)
        user_vector = self.user_course_matrix[user_idx].copy()
        
        # Центрировать вектор (вычесть среднее)
        rated_mask = user_vector > 0
        if not rated_mask.any():
            return []
        
        user_vector[rated_mask] -= self.user_means[user_idx]
        
        similarities = []
        for idx, other_user_id in enumerate(self.user_ids):
            if other_user_id == student_id:
                continue
            
            other_vector = self.user_course_matrix[idx].copy()
            other_rated_mask = other_vector > 0
            
            if not other_rated_mask.any():
                continue
            
            # Центрировать
            other_vector[other_rated_mask] -= self.user_means[idx]
            
            # Найти общие оценки
            common_mask = rated_mask & other_rated_mask
            
            if not common_mask.any():
                continue
            
            # Centered cosine similarity только по общим элементам
            dot_product = np.dot(user_vector[common_mask], other_vector[common_mask])
            norm_user = np.linalg.norm(user_vector[common_mask])
            norm_other = np.linalg.norm(other_vector[common_mask])
            
            if norm_user > 0 and norm_other > 0:
                similarity = dot_product / (norm_user * norm_other)
                # Учитываем количество общих оценок
                confidence = min(1.0, common_mask.sum() / 5.0)  # Полная уверенность при 5+ общих
                adjusted_similarity = similarity * confidence
                similarities.append((other_user_id, adjusted_similarity))
        
        similarities.sort(key=lambda x: x[1], reverse=True)
        return similarities[:top_n]
    
    def recommend(self, student_id: int, top_n: int = 10) -> List[Dict]:
        """
        Рекомендовать КУРСЫ с правильной нормализацией scores
        """
        similar_users = self.find_similar_users(student_id, top_n=10)
        
        # Cold start: если нет похожих пользователей
        if not similar_users:
            return self._get_popular_courses_smart(student_id, top_n)
        
        similar_user_ids = [uid for uid, _ in similar_users]
        similarity_scores = {uid: sim for uid, sim in similar_users}
        
        # Найти курсы с взвешенными оценками
        query = """
            WITH similar_user_courses AS (
                SELECT 
                    c.id as course_id,
                    c.title,
                    c.description,
                    c.difficulty_level,
                    c.subject,
                    rr.student_id,
                    AVG(rr.rating) as avg_rating,
                    COUNT(DISTINCT rr.student_id) as rating_count
                FROM resource_ratings rr
                JOIN resources r ON rr.resource_id = r.id
                JOIN modules m ON r.module_id = m.id
                JOIN courses c ON m.course_id = c.id
                WHERE rr.student_id IN %s
                  AND rr.rating >= 3.5
                  AND c.is_published = true
                  AND c.id NOT IN (
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
                      ) * 0.9
                  )
                GROUP BY c.id, c.title, c.description, c.difficulty_level, c.subject, rr.student_id
            )
            SELECT 
                course_id,
                title,
                description,
                difficulty_level,
                subject,
                AVG(avg_rating) as final_rating,
                COUNT(DISTINCT student_id) as num_similar_users
            FROM similar_user_courses
            GROUP BY course_id, title, description, difficulty_level, subject
            ORDER BY final_rating DESC, num_similar_users DESC
            LIMIT %s
        """
        
        recommendations = self.db.execute(
            query, 
            (tuple(similar_user_ids), student_id, top_n * 2)
        )
        
        result = []
        for rec in recommendations:
            # Нормализованный score: (rating - 3.5) / 1.5 для диапазона 3.5-5
            # Результат: 0.0 (rating=3.5) до 1.0 (rating=5.0)
            base_score = max(0.0, min(1.0, (rec['final_rating'] - 3.5) / 1.5))
            
            # Бонус за количество похожих пользователей, которым понравилось
            popularity_bonus = min(0.2, rec['num_similar_users'] / len(similar_users) * 0.2)
            
            final_score = min(1.0, base_score + popularity_bonus)
            
            result.append({
                'course_id': rec['course_id'],
                'title': rec['title'],
                'description': rec['description'],
                'difficulty_level': rec['difficulty_level'],
                'subject': rec['subject'],
                'score': round(final_score, 3),
                'algorithm': 'collaborative',
                'reason': f"Похожим ученикам понравилось (★{rec['final_rating']:.1f}, {rec['num_similar_users']} чел.)",
                'details': {
                    'avg_rating': float(rec['final_rating']),
                    'num_similar_users': rec['num_similar_users']
                }
            })
        
        return result[:top_n]
    
    def _get_popular_courses_smart(self, student_id: int, top_n: int = 10) -> List[Dict]:
        """
        Cold start: популярные курсы с правильными scores
        Учитываем возрастную группу студента
        """
        # Получить возрастную группу студента
        age_query = """
            SELECT age_group, grade
            FROM student_profiles
            WHERE user_id = %s
        """
        student_info = self.db.execute_one(age_query, (student_id,))
        age_group = student_info['age_group'] if student_info else None
        
        query = """
            SELECT 
                c.id as course_id,
                c.title,
                c.description,
                c.difficulty_level,
                c.subject,
                c.age_group,
                COALESCE(AVG(rr.rating), 3.5) as avg_rating,
                COUNT(DISTINCT rr.student_id) as student_count,
                COUNT(rr.id) as total_ratings
            FROM courses c
            JOIN modules m ON c.id = m.course_id
            JOIN resources r ON m.id = r.module_id
            LEFT JOIN resource_ratings rr ON r.id = rr.resource_id
            WHERE c.is_published = true
            GROUP BY c.id, c.title, c.description, c.difficulty_level, c.subject, c.age_group
            HAVING COUNT(rr.id) >= 3
            ORDER BY avg_rating DESC, student_count DESC
            LIMIT %s
        """
        
        courses = self.db.execute(query, (top_n * 2,))
        
        result = []
        for course in courses:
            # Нормализация: (rating - 3.5) / 1.5
            rating_score = max(0.0, min(1.0, (course['avg_rating'] - 3.5) / 1.5))
            
            # Популярность: нормализовать по логарифму
            popularity_score = min(1.0, np.log1p(course['total_ratings']) / 10.0)
            
            # Бонус за совпадение возрастной группы
            age_bonus = 0.15 if course['age_group'] == age_group else 0
            
            # Итоговый score: 60% рейтинг + 30% популярность + 10% возраст
            final_score = min(1.0, 
                0.6 * rating_score + 
                0.3 * popularity_score + 
                age_bonus
            )
            
            result.append({
                'course_id': course['course_id'],
                'title': course['title'],
                'description': course['description'],
                'difficulty_level': course['difficulty_level'],
                'subject': course['subject'],
                'score': round(final_score, 3),
                'algorithm': 'collaborative_cold_start',
                'reason': f"Популярный курс (★{course['avg_rating']:.1f}, {course['student_count']} учеников)",
                'details': {
                    'avg_rating': float(course['avg_rating']),
                    'student_count': course['student_count'],
                    'is_cold_start': True
                }
            })
        
        # Сортировать по финальному score
        result.sort(key=lambda x: x['score'], reverse=True)
        return result[:top_n]