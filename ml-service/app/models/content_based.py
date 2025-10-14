from typing import List, Dict


class ContentBasedFiltering:
    """
    Content-Based Filtering - рекомендует на основе интересов и предпочтений ученика
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
    
    def get_student_preferred_types(self, student_id: int) -> List[str]:
        """Получить предпочитаемые типы ресурсов"""
        query = """
            SELECT preferences
            FROM student_profiles
            WHERE user_id = %s
        """
        
        result = self.db.execute_one(query, (student_id,))
        
        if result and result['preferences']:
            prefs = result['preferences']
            # ИСПРАВЛЕНО: Безопасная проверка типа
            if isinstance(prefs, dict):
                return prefs.get('preferred_resource_types', [])
        return []
    
    def recommend(self, student_id: int, top_n: int = 10) -> List[Dict]:
        """
        Рекомендовать ресурсы на основе интересов и предпочтений
        """
        interests = self.get_student_interests(student_id)
        preferred_types = self.get_student_preferred_types(student_id)
        
        if not interests:
            return []
        
        # ИСПРАВЛЕНО: Используем IN вместо ANY
        query = """
            WITH matched_resources AS (
                SELECT 
                    r.id as resource_id,
                    r.title,
                    r.resource_type,
                    r.difficulty,
                    COUNT(DISTINCT t.id) as matching_tags,
                    SUM(rt.weight) as tag_weight,
                    ARRAY_AGG(DISTINCT t.name) as matched_tag_names
                FROM resources r
                JOIN resource_tags rt ON r.id = rt.resource_id
                JOIN tags t ON rt.tag_id = t.id
                WHERE t.name IN %s
                  AND r.id NOT IN (
                      SELECT resource_id 
                      FROM student_progress 
                      WHERE student_id = %s 
                        AND status = 'completed'
                  )
                GROUP BY r.id, r.title, r.resource_type, r.difficulty
            )
            SELECT *
            FROM matched_resources
            ORDER BY matching_tags DESC, tag_weight DESC
            LIMIT %s
        """
        
        # ИСПРАВЛЕНО: Передаем tuple вместо списка
        recommendations = self.db.execute(
            query,
            (tuple(interests), student_id, top_n)
        )
        
        result = []
        for rec in recommendations:
            # Рассчитать score
            base_score = float(rec['matching_tags']) / len(interests)
            
            # Бонус если тип ресурса совпадает с предпочтениями
            type_bonus = 0.2 if rec['resource_type'] in preferred_types else 0
            
            score = min(base_score + type_bonus, 1.0)
            
            # ИСПРАВЛЕНО: Безопасная обработка массива тегов
            tag_names = rec['matched_tag_names'][:3] if rec['matched_tag_names'] else []
            
            result.append({
                'resource_id': rec['resource_id'],
                'title': rec['title'],
                'score': score,
                'algorithm': 'content_based',
                'reason': f"Соответствует вашим интересам: {', '.join(tag_names)}"
            })
        
        return result