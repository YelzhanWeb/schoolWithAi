from typing import List, Dict, Tuple


class KnowledgeBasedFiltering:
    """
    Knowledge-Based Filtering для КУРСОВ (исправленная версия)
    - Адаптивный подбор курсов для развития навыков
    - Работает даже если нет явно слабых навыков
    - Правильная нормализация scores
    """
    
    def __init__(self, db):
        self.db = db
    
    def get_skill_gaps(self, student_id: int) -> List[Tuple[str, float, str]]:
        """
        Найти gaps в навыках:
        1. Слабые навыки (< 0.5)
        2. Средние навыки для улучшения (0.5 - 0.75)
        3. Навыки для углубления (> 0.75)
        
        Возвращает: [(skill_name, proficiency_level, priority)]
        """
        query = """
            SELECT 
                skill_name, 
                proficiency_level,
                CASE 
                    WHEN proficiency_level < 0.5 THEN 'weak'
                    WHEN proficiency_level < 0.75 THEN 'medium'
                    ELSE 'advanced'
                END as priority
            FROM student_skills
            WHERE student_id = %s
            ORDER BY proficiency_level ASC, skill_name
        """
        
        skills = self.db.execute(query, (student_id,))
        
        if not skills:
            # Если нет навыков - вернуть базовые для класса
            return self._get_grade_level_skills(student_id)
        
        return [(s['skill_name'], s['proficiency_level'], s['priority']) for s in skills]
    
    def _get_grade_level_skills(self, student_id: int) -> List[Tuple[str, float, str]]:
        """
        Получить базовые навыки для класса студента (cold start)
        """
        query = """
            SELECT grade, subject
            FROM student_profiles
            WHERE user_id = %s
        """
        
        profile = self.db.execute_one(query, (student_id,))
        
        # Базовые навыки по предметам для разных классов
        grade_skills = {
            'math': ['algebra', 'geometry', 'arithmetic', 'problem_solving'],
            'science': ['biology', 'chemistry', 'physics', 'experiments'],
            'language': ['reading', 'writing', 'grammar', 'vocabulary'],
            'history': ['world_history', 'critical_thinking', 'analysis'],
            'programming': ['algorithms', 'data_structures', 'python', 'problem_solving']
        }
        
        # Вернуть базовые навыки с низким уровнем для развития
        default_skills = []
        for subject, skills in grade_skills.items():
            for skill in skills[:3]:  # Первые 3 навыка
                default_skills.append((skill, 0.3, 'weak'))
        
        return default_skills[:5]
    
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
        Рекомендовать КУРСЫ для развития навыков с адаптивной стратегией
        """
        skill_gaps = self.get_skill_gaps(student_id)
        
        if not skill_gaps:
            return []
        
        grade = self.get_student_grade(student_id)
        
        # Приоритизация: сначала слабые, потом средние, потом продвинутые
        weak_skills = [s for s, _, p in skill_gaps if p == 'weak']
        medium_skills = [s for s, _, p in skill_gaps if p == 'medium']
        advanced_skills = [s for s, _, p in skill_gaps if p == 'advanced']
        
        # Стратегия: 70% слабые, 20% средние, 10% продвинутые
        target_skills = []
        target_skills.extend(weak_skills[:5])
        target_skills.extend(medium_skills[:3])
        target_skills.extend(advanced_skills[:2])
        
        if not target_skills:
            target_skills = [s for s, _, _ in skill_gaps[:5]]
        
        # Рекомендуемая сложность
        min_difficulty = max(1, (grade // 3))
        max_difficulty = min(5, (grade // 2) + 2)
        
        query = """
            WITH skill_courses AS (
                SELECT 
                    c.id as course_id,
                    c.title,
                    c.description,
                    c.difficulty_level,
                    c.subject,
                    COUNT(DISTINCT t.name) as skills_covered,
                    ARRAY_AGG(DISTINCT t.name) as skill_names,
                    COALESCE(AVG(rr.rating), 3.5) as avg_rating,
                    COUNT(DISTINCT rr.student_id) as rating_count
                FROM courses c
                JOIN modules m ON c.id = m.course_id
                JOIN resources r ON m.id = r.module_id
                JOIN resource_tags rt ON r.id = rt.resource_id
                JOIN tags t ON rt.tag_id = t.id
                LEFT JOIN resource_ratings rr ON r.id = rr.resource_id
                WHERE t.name IN %s
                  AND c.difficulty_level BETWEEN %s AND %s
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
                GROUP BY c.id, c.title, c.description, c.difficulty_level, c.subject
            )
            SELECT *
            FROM skill_courses
            ORDER BY 
                skills_covered DESC,
                avg_rating DESC,
                difficulty_level ASC
            LIMIT %s
        """
        
        recommendations = self.db.execute(
            query,
            (tuple(target_skills), min_difficulty, max_difficulty, student_id, top_n * 2)
        )
        
        # Создать карту приоритетов навыков
        skill_priority_map = {}
        for skill, level, priority in skill_gaps:
            if priority == 'weak':
                skill_priority_map[skill] = 1.0
            elif priority == 'medium':
                skill_priority_map[skill] = 0.7
            else:
                skill_priority_map[skill] = 0.4
        
        result = []
        for rec in recommendations:
            # 1. Score за покрытие навыков
            covered_skills = rec['skill_names'] if rec['skill_names'] else []
            
            # Взвешенная сумма приоритетов покрытых навыков
            priority_sum = sum(skill_priority_map.get(s, 0.5) for s in covered_skills if s in target_skills)
            max_priority = len([s for s in covered_skills if s in target_skills])
            
            if max_priority > 0:
                skill_coverage_score = min(1.0, priority_sum / max_priority)
            else:
                skill_coverage_score = 0
            
            # 2. Бонус за количество покрытых навыков
            coverage_bonus = min(0.2, (rec['skills_covered'] / len(target_skills)) * 0.2)
            
            # 3. Бонус за подходящую сложность (легче для слабых навыков)
            if weak_skills and any(s in covered_skills for s in weak_skills):
                # Для слабых навыков лучше простые курсы
                ideal_difficulty = min_difficulty + 1
            else:
                ideal_difficulty = (min_difficulty + max_difficulty) // 2
            
            difficulty_diff = abs(rec['difficulty_level'] - ideal_difficulty)
            difficulty_score = max(0, (5 - difficulty_diff) / 5.0) * 0.15
            
            # 4. Бонус за рейтинг
            if rec['rating_count'] >= 3:
                rating_score = max(0, (rec['avg_rating'] - 3.5) / 1.5) * 0.1
            else:
                rating_score = 0
            
            # Итоговый score
            final_score = min(1.0,
                skill_coverage_score * 0.55 +
                coverage_bonus +
                difficulty_score +
                rating_score
            )
            
            # Определить какие навыки курс развивает
            target_covered = [s for s in covered_skills if s in target_skills][:3]
            
            # Сформировать причину
            if any(s in weak_skills for s in target_covered):
                reason_prefix = "Поможет укрепить слабые навыки"
            elif any(s in medium_skills for s in target_covered):
                reason_prefix = "Развитие навыков"
            else:
                reason_prefix = "Углубленное изучение"
            
            result.append({
                'course_id': rec['course_id'],
                'title': rec['title'],
                'description': rec['description'],
                'difficulty_level': rec['difficulty_level'],
                'subject': rec['subject'],
                'score': round(final_score, 3),
                'algorithm': 'knowledge_based',
                'reason': f"{reason_prefix}: {', '.join(target_covered)}",
                'details': {
                    'skills_covered': rec['skills_covered'],
                    'target_skills': target_covered,
                    'difficulty_match': difficulty_diff <= 1,
                    'avg_rating': float(rec['avg_rating'])
                }
            })
        
        # Сортировка по score
        result.sort(key=lambda x: x['score'], reverse=True)
        return result[:top_n]
    
    def update_skill(self, student_id: int, skill_name: str, test_score: float):
        """
        Обновить уровень навыка с улучшенным алгоритмом обновления
        """
        query_get = """
            SELECT proficiency_level, updated_at
            FROM student_skills
            WHERE student_id = %s AND skill_name = %s
        """
        
        result = self.db.execute_one(query_get, (student_id, skill_name))
        
        if result:
            current_level = result['proficiency_level']
            
            # Adaptive learning rate: быстрее обучение для низких уровней
            if current_level < 0.3:
                alpha = 0.4  # Быстрое обучение
            elif current_level < 0.6:
                alpha = 0.3
            else:
                alpha = 0.2  # Медленное обучение на высоких уровнях
            
            # Exponential moving average
            new_level = (1 - alpha) * current_level + alpha * test_score
        else:
            # Первый тест - используем результат напрямую, но с консервативностью
            new_level = test_score * 0.8
        
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