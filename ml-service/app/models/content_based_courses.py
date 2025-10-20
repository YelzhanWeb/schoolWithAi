from typing import List, Dict


class ContentBasedFilteringCourses:
    """
    Content-Based Filtering для КУРСОВ
    Рекомендует курсы на основе интересов и предпочтений ученика
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
    
    def get_student_grade(self, student_id: int) -> int:
        """Получить класс ученика"""
        query = """
            SELECT grade, age_group
            FROM student_profiles
            WHERE user_id = %s
        """
        
        result = self.db.execute_one(query, (student_id,))
        if result:
            return result['grade'], result['age_group']
        return 5, 'middle'
    
    def recommend(self, student_id: int, top_n: int = 10) -> List[Dict]:
        """
        Рекомендовать КУРСЫ на основе интересов
        """
        interests = self.get_student_interests(student_id)
        grade, age_group = self.get_student_grade(student_id)
        
        if not interests:
            return []
        
        # Найти курсы, которые соответствуют интересам студента
        query = """
            WITH course_tags AS (
                SELECT 
                    c.id as course_id,
                    c.title,
                    c.description,
                    c.difficulty_level,
                    c.subject,
                    c.age_group,
                    COUNT(DISTINCT t.id) as matching_tags,
                    ARRAY_AGG(DISTINCT t.name) as matched_tag_names
                FROM courses c
                JOIN modules m ON c.id = m.course_id
                JOIN resources r ON m.id = r.module_id
                JOIN resource_tags rt ON r.id = rt.resource_id
                JOIN tags t ON rt.tag_id = t.id
                WHERE t.name IN %s
                  AND c.is_published = true
                  AND c.id NOT IN (
                      -- Исключить завершенные курсы
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
                      ) * 0.8
                  )
                GROUP BY c.id, c.title, c.description, c.difficulty_level, c.subject, c.age_group
            )
            SELECT *
            FROM course_tags
            ORDER BY matching_tags DESC
            LIMIT %s
        """
        
        recommendations = self.db.execute(
            query,
            (tuple(interests), student_id, top_n)
        )
        
        result = []
        for rec in recommendations:
            # Рассчитать score
            base_score = float(rec['matching_tags']) / len(interests)
            
            # Бонус за подходящую возрастную группу
            age_bonus = 0.2 if rec['age_group'] == age_group else 0
            
            # Бонус за подходящую сложность (относительно класса)
            expected_difficulty = min(max(1, (grade // 2)), 5)
            difficulty_diff = abs(rec['difficulty_level'] - expected_difficulty)
            difficulty_bonus = max(0, (5 - difficulty_diff) / 10)
            
            score = min(base_score + age_bonus + difficulty_bonus, 1.0)
            
            tag_names = rec['matched_tag_names'][:3] if rec['matched_tag_names'] else []
            
            result.append({
                'course_id': rec['course_id'],
                'title': rec['title'],
                'description': rec['description'],
                'difficulty_level': rec['difficulty_level'],
                'subject': rec['subject'],
                'score': score,
                'algorithm': 'content_based',
                'reason': f"Соответствует вашим интересам: {', '.join(tag_names)}"
            })
        
        return result