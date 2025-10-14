from typing import List, Dict, Tuple


class KnowledgeBasedFiltering:
    """
    Knowledge-Based Filtering - рекомендует на основе уровня навыков
    Определяет слабые места и подбирает материалы для улучшения
    """
    
    def __init__(self, db):
        self.db = db
    
    def get_weak_skills(self, student_id: int, threshold: float = 0.5) -> List[Tuple[str, float]]:
        """
        Найти слабые навыки ученика
        Returns: [(skill_name, proficiency_level), ...]
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
        Рекомендовать ресурсы для улучшения слабых навыков
        """
        weak_skills = self.get_weak_skills(student_id)
        
        if not weak_skills:
            return []
        
        # Берем топ-3 самых слабых навыка
        target_skills = [skill for skill, _ in weak_skills[:3]]
        grade = self.get_student_grade(student_id)
        
        # ИСПРАВЛЕНО: Используем IN вместо ANY
        query = """
            WITH skill_resources AS (
                SELECT 
                    r.id as resource_id,
                    r.title,
                    r.difficulty,
                    r.resource_type,
                    t.name as skill_name,
                    rt.weight
                FROM resources r
                JOIN resource_tags rt ON r.id = rt.resource_id
                JOIN tags t ON rt.tag_id = t.id
                WHERE t.name IN %s
                  AND r.difficulty <= 3
                  AND r.id NOT IN (
                      SELECT resource_id 
                      FROM student_progress 
                      WHERE student_id = %s 
                        AND status = 'completed'
                  )
            )
            SELECT 
                resource_id,
                title,
                difficulty,
                resource_type,
                COUNT(DISTINCT skill_name) as skills_covered,
                ARRAY_AGG(DISTINCT skill_name) as skill_names
            FROM skill_resources
            GROUP BY resource_id, title, difficulty, resource_type
            ORDER BY 
                skills_covered DESC,
                difficulty ASC
            LIMIT %s
        """
        
        # ИСПРАВЛЕНО: Передаем tuple вместо списка
        recommendations = self.db.execute(
            query,
            (tuple(target_skills), student_id, top_n)
        )
        
        result = []
        for rec in recommendations:
            # Score зависит от того, сколько слабых навыков покрывает ресурс
            skills_score = float(rec['skills_covered']) / len(target_skills)
            
            # Бонус за подходящую сложность
            difficulty_bonus = (4 - rec['difficulty']) / 10.0
            
            score = min(skills_score + difficulty_bonus, 1.0)
            
            # ИСПРАВЛЕНО: Безопасная обработка массива навыков
            skill_names = rec['skill_names'] if rec['skill_names'] else []
            
            result.append({
                'resource_id': rec['resource_id'],
                'title': rec['title'],
                'score': score,
                'algorithm': 'knowledge_based',
                'reason': f"Поможет улучшить: {', '.join(skill_names)}"
            })
        
        return result
    
    def update_skill(self, student_id: int, skill_name: str, test_score: float):
        """
        Обновить уровень навыка на основе результата теста
        test_score: 0.0-1.0 (процент правильных ответов)
        """
        # Получить текущий уровень
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
        
        # Ограничить 0-1
        new_level = max(0.0, min(1.0, new_level))
        
        # ИСПРАВЛЕНО: Правильный синтаксис для upsert
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