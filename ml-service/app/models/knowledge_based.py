from typing import List, Dict, Tuple


class KnowledgeBasedFiltering:
    """
    Knowledge-Based Filtering для КУРСОВ
    Рекомендует курсы для улучшения слабых навыков
    """
    
    def __init__(self, db):
        self.db = db
    
    def get_weak_skills(self, student_id: int, threshold: float = 0.5) -> List[Tuple[str, float]]:
        """
        Найти слабые навыки ученика
        """
        query = """
            SELECT skill_name, proficiency_level
            FROM student_skills
            WHERE student_id = %s
              AND proficiency_level < %s
            ORDER BY proficiency_level ASC
        """
        
        skills = self.db.execute(query, (student_id, threshold))
        return [(s['skill_name'], s['proficiency_level']) for s in skills]
    
    def get_student_grade(self, student_id: int) -> int:
        """Получить класс ученика"""
        query = """
            SELECT grade
            FROM student_profiles
            WHERE user_id = %s
        """
        
        result = self.db.execute_one(query, (student_id,))
        return result['grade'] if result else 5
    
    def recommend(self, student_id: int, top_n: int = 10) -> List[Dict]:
        """
        Рекомендовать КУРСЫ для улучшения слабых навыков
        """
        weak_skills = self.get_weak_skills(student_id)
        
        if not weak_skills:
            return []
        
        # Берем топ-3 самых слабых навыка
        target_skills = [skill for skill, _ in weak_skills[:3]]
        grade = self.get_student_grade(student_id)
        
        # Рекомендуемая сложность на основе класса
        max_difficulty = min(max(1, (grade // 2) + 1), 5)
        
        query = """
            WITH skill_courses AS (
                SELECT 
                    c.id as course_id,
                    c.title,
                    c.description,
                    c.difficulty_level,
                    c.subject,
                    COUNT(DISTINCT t.name) as skills_covered,
                    ARRAY_AGG(DISTINCT t.name) as skill_names
                FROM courses c
                JOIN modules m ON c.id = m.course_id
                JOIN resources r ON m.id = r.module_id
                JOIN resource_tags rt ON r.id = rt.resource_id
                JOIN tags t ON rt.tag_id = t.id
                WHERE t.name IN %s
                  AND c.difficulty_level <= %s
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
                GROUP BY c.id, c.title, c.description, c.difficulty_level, c.subject
            )
            SELECT *
            FROM skill_courses
            ORDER BY 
                skills_covered DESC,
                difficulty_level ASC
            LIMIT %s
        """
        
        recommendations = self.db.execute(
            query,
            (tuple(target_skills), max_difficulty, student_id, top_n)
        )
        
        result = []
        for rec in recommendations:
            # Score зависит от того, сколько слабых навыков покрывает курс
            skills_score = float(rec['skills_covered']) / len(target_skills)
            
            # Бонус за подходящую сложность
            difficulty_bonus = (max_difficulty - rec['difficulty_level'] + 1) / 10.0
            
            score = min(skills_score + difficulty_bonus, 1.0)
            
            skill_names = rec['skill_names'] if rec['skill_names'] else []
            
            result.append({
                'course_id': rec['course_id'],
                'title': rec['title'],
                'description': rec['description'],
                'difficulty_level': rec['difficulty_level'],
                'subject': rec['subject'],
                'score': score,
                'algorithm': 'knowledge_based',
                'reason': f"Поможет улучшить: {', '.join(skill_names[:3])}"
            })
        
        return result
    
    def update_skill(self, student_id: int, skill_name: str, test_score: float):
        """
        Обновить уровень навыка на основе результата теста
        """
        query_get = """
            SELECT proficiency_level
            FROM student_skills
            WHERE student_id = %s AND skill_name = %s
        """
        
        result = self.db.execute_one(query_get, (student_id, skill_name))
        
        if result:
            current_level = result['proficiency_level']
            # Exponential moving average
            new_level = 0.7 * current_level + 0.3 * test_score
        else:
            new_level = test_score
        
        new_level = max(0.0, min(1.0, new_level))
        
        query_upsert = """
            INSERT INTO student_skills (student_id, skill_name, proficiency_level, updated_at)
            VALUES (%s, %s, %s, NOW())
            ON CONFLICT (student_id, skill_name) 
            DO UPDATE SET 
                proficiency_level = EXCLUDED.proficiency_level,
                updated_at = NOW()
        """
        
        self.db.cursor.execute(query_upsert, (student_id, skill_name, new_level))
        self.db.conn.commit()
        
        return new_level