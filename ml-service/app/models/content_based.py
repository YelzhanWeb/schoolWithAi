from typing import List, Dict
import numpy as np


class ContentBasedFiltering:
    """
    Content-Based Filtering для КУРСОВ (исправленная версия)
    - Cold start для новых пользователей
    - Правильная нормализация scores
    - Учет истории просмотров
    """
    
    def __init__(self, db):
        self.db = db
    
    def get_student_interests(self, student_id: int) -> List[str]:
        """Получить интересы ученика"""
        query = """
            SELECT interests
            FROM student_profiles
            WHERE user_id = %s
        """
        
        result = self.db.execute_one(query, (student_id,))
        
        if result and result['interests']:
            return result['interests']
        return []
    
    def infer_interests_from_history(self, student_id: int) -> List[str]:
        """
        Вывести интересы из истории: курсы с высокими оценками
        """
        query = """
            SELECT DISTINCT t.name
            FROM resource_ratings rr
            JOIN resources r ON rr.resource_id = r.id
            JOIN resource_tags rt ON r.id = rt.resource_id
            JOIN tags t ON rt.tag_id = t.id
            WHERE rr.student_id = %s
              AND rr.rating >= 4
            GROUP BY t.name
            ORDER BY COUNT(*) DESC
            LIMIT 5
        """
        
        tags = self.db.execute(query, (student_id,))
        return [tag['name'] for tag in tags]
    
    def get_student_profile(self, student_id: int) -> Dict:
        """Получить полный профиль ученика"""
        query = """
            SELECT grade, age_group, learning_style
            FROM student_profiles
            WHERE user_id = %s
        """
        
        result = self.db.execute_one(query, (student_id,))
        
        if result:
            return {
                'grade': result['grade'],
                'age_group': result['age_group'],
                'learning_style': result.get('learning_style', 'mixed')
            }
        
        return {'grade': 5, 'age_group': 'middle', 'learning_style': 'mixed'}
    
    def recommend(self, student_id: int, top_n: int = 10) -> List[Dict]:
        """
        Рекомендовать КУРСЫ с обработкой cold start
        """
        # Получить интересы
        interests = self.get_student_interests(student_id)
        
        # Cold start: вывести интересы из истории
        if not interests:
            interests = self.infer_interests_from_history(student_id)
        
        # Если все еще нет интересов - вернуть популярные для возраста
        if not interests:
            return self._get_popular_for_age_group(student_id, top_n)
        
        profile = self.get_student_profile(student_id)
        
        # Найти курсы по интересам
        query = """
            WITH course_matches AS (
                SELECT 
                    c.id as course_id,
                    c.title,
                    c.description,
                    c.difficulty_level,
                    c.subject,
                    c.age_group,
                    COUNT(DISTINCT t.id) as matching_tags,
                    ARRAY_AGG(DISTINCT t.name) as matched_tag_names,
                    AVG(rr.rating) as avg_rating,
                    COUNT(DISTINCT rr.student_id) as rating_count
                FROM courses c
                JOIN modules m ON c.id = m.course_id
                JOIN resources r ON m.id = r.module_id
                JOIN resource_tags rt ON r.id = rt.resource_id
                JOIN tags t ON rt.tag_id = t.id
                LEFT JOIN resource_ratings rr ON r.id = rr.resource_id
                WHERE t.name IN %s
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
                GROUP BY c.id, c.title, c.description, c.difficulty_level, c.subject, c.age_group
            )
            SELECT *
            FROM course_matches
            ORDER BY matching_tags DESC, avg_rating DESC NULLS LAST
            LIMIT %s
        """
        
        recommendations = self.db.execute(
            query,
            (tuple(interests), student_id, top_n * 2)
        )
        
        result = []
        for rec in recommendations:
            # 1. Score за совпадение интересов (0.0 - 1.0)
            interest_score = min(1.0, float(rec['matching_tags']) / len(interests))
            
            # 2. Бонус за возрастную группу
            age_bonus = 0.2 if rec['age_group'] == profile['age_group'] else 0
            
            # 3. Бонус за подходящую сложность
            expected_difficulty = min(max(1, (profile['grade'] // 2)), 5)
            difficulty_diff = abs(rec['difficulty_level'] - expected_difficulty)
            difficulty_score = max(0, (5 - difficulty_diff) / 5.0) * 0.15
            
            # 4. Бонус за рейтинг (если есть оценки)
            if rec['avg_rating'] and rec['rating_count'] >= 3:
                rating_score = max(0, (rec['avg_rating'] - 3.5) / 1.5) * 0.15
            else:
                rating_score = 0
            
            # Итоговый score
            final_score = min(1.0, 
                interest_score * 0.5 + 
                age_bonus + 
                difficulty_score + 
                rating_score
            )
            
            # Топ-3 совпавших тега
            tag_names = rec['matched_tag_names'][:3] if rec['matched_tag_names'] else []
            
            result.append({
                'course_id': rec['course_id'],
                'title': rec['title'],
                'description': rec['description'],
                'difficulty_level': rec['difficulty_level'],
                'subject': rec['subject'],
                'score': round(final_score, 3),
                'algorithm': 'content_based',
                'reason': f"Соответствует вашим интересам: {', '.join(tag_names)}",
                'details': {
                    'matching_tags': rec['matching_tags'],
                    'matched_interests': tag_names,
                    'age_match': rec['age_group'] == profile['age_group'],
                    'difficulty_match': difficulty_diff <= 1
                }
            })
        
        # Сортировка по score
        result.sort(key=lambda x: x['score'], reverse=True)
        return result[:top_n]
    
    def _get_popular_for_age_group(self, student_id: int, top_n: int = 10) -> List[Dict]:
        """
        Cold start: популярные курсы для возрастной группы студента
        """
        profile = self.get_student_profile(student_id)
        
        query = """
            SELECT 
                c.id as course_id,
                c.title,
                c.description,
                c.difficulty_level,
                c.subject,
                c.age_group,
                COALESCE(AVG(rr.rating), 3.5) as avg_rating,
                COUNT(DISTINCT rr.student_id) as student_count
            FROM courses c
            JOIN modules m ON c.id = m.course_id
            JOIN resources r ON m.id = r.module_id
            LEFT JOIN resource_ratings rr ON r.id = rr.resource_id
            WHERE c.is_published = true
              AND c.age_group = %s
            GROUP BY c.id, c.title, c.description, c.difficulty_level, c.subject, c.age_group
            ORDER BY avg_rating DESC, student_count DESC
            LIMIT %s
        """
        
        courses = self.db.execute(query, (profile['age_group'], top_n))
        
        result = []
        for course in courses:
            # Нормализация рейтинга
            rating_score = max(0.0, min(1.0, (course['avg_rating'] - 3.5) / 1.5))
            
            # Популярность
            popularity_score = min(0.3, np.log1p(course['student_count']) / 10.0)
            
            final_score = min(1.0, 0.7 * rating_score + popularity_score)
            
            result.append({
                'course_id': course['course_id'],
                'title': course['title'],
                'description': course['description'],
                'difficulty_level': course['difficulty_level'],
                'subject': course['subject'],
                'score': round(final_score, 3),
                'algorithm': 'content_based_cold_start',
                'reason': f"Популярный курс для вашей возрастной группы (★{course['avg_rating']:.1f})",
                'details': {
                    'is_cold_start': True,
                    'avg_rating': float(course['avg_rating']),
                    'student_count': course['student_count']
                }
            })
        
        return result