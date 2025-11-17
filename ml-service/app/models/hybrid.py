from typing import List, Dict
from .collaborative import CollaborativeFiltering
from .content_based import ContentBasedFiltering
from .knowledge_based import KnowledgeBasedFiltering


class HybridRecommender:
    """
    Гибридная рекомендательная система для КУРСОВ (исправленная версия)
    
    Улучшения:
    - Адаптивные веса в зависимости от доступности данных
    - Правильная нормализация и объединение scores
    - Diversity boost для разнообразия рекомендаций
    - Умная обработка случаев, когда некоторые алгоритмы не дают результатов
    """
    
    def __init__(self, db):
        self.db = db
        self.collaborative = CollaborativeFiltering(db)
        self.content_based = ContentBasedFiltering(db)
        self.knowledge_based = KnowledgeBasedFiltering(db)
        
        # Базовые веса (будут адаптироваться)
        self.base_weights = {
            'collaborative': 0.4,
            'content_based': 0.3,
            'knowledge_based': 0.3
        }
    
    def _calculate_adaptive_weights(self, collab_recs, content_recs, knowledge_recs) -> Dict[str, float]:
        """
        Адаптивные веса: перераспределяем вес отсутствующих алгоритмов
        """
        available = {
            'collaborative': len(collab_recs) > 0,
            'content_based': len(content_recs) > 0,
            'knowledge_based': len(knowledge_recs) > 0
        }
        
        # Если все алгоритмы работают - используем базовые веса
        if all(available.values()):
            return self.base_weights.copy()
        
        # Перераспределяем веса
        weights = {}
        total_base_weight = 0
        available_base_weight = 0
        
        for algo, base_weight in self.base_weights.items():
            total_base_weight += base_weight
            if available[algo]:
                available_base_weight += base_weight
        
        # Нормализуем веса только по доступным алгоритмам
        for algo, base_weight in self.base_weights.items():
            if available[algo]:
                weights[algo] = base_weight / available_base_weight if available_base_weight > 0 else 0
            else:
                weights[algo] = 0
        
        return weights
    
    def _diversity_boost(self, recommendations: List[Dict]) -> List[Dict]:
        """
        Boost разнообразия: понижаем score дубликатов по subject
        """
        seen_subjects = {}
        
        for rec in recommendations:
            subject = rec.get('subject', 'other')
            
            if subject in seen_subjects:
                # Понижаем score за повторяющийся предмет
                seen_subjects[subject] += 1
                penalty = 0.05 * seen_subjects[subject]
                rec['score'] = max(0, rec['score'] - penalty)
            else:
                seen_subjects[subject] = 0
        
        return recommendations
    
    def recommend(self, student_id: int, top_n: int = 10) -> List[Dict]:
        """
        Генерация гибридных рекомендаций КУРСОВ с адаптивными весами
        """
        # 1. Получить рекомендации от каждого алгоритма
        collab_recs = self.collaborative.recommend(student_id, top_n=top_n * 2)
        content_recs = self.content_based.recommend(student_id, top_n=top_n * 2)
        knowledge_recs = self.knowledge_based.recommend(student_id, top_n=top_n * 2)
        
        # Если НИ ОДИН алгоритм не дал результатов
        if not collab_recs and not content_recs and not knowledge_recs:
            print(f"⚠️ No recommendations from any algorithm for student {student_id}")
            return self._get_fallback_recommendations(student_id, top_n)
        
        # 2. Адаптивные веса
        weights = self._calculate_adaptive_weights(collab_recs, content_recs, knowledge_recs)
        
        print(f"📊 Adaptive weights: {weights}")
        print(f"📈 Recs count - Collab: {len(collab_recs)}, Content: {len(content_recs)}, Knowledge: {len(knowledge_recs)}")
        
        # 3. Объединить все рекомендации
        all_recommendations = {}
        
        # Добавить collaborative
        for rec in collab_recs:
            course_id = rec['course_id']
            all_recommendations[course_id] = {
                'course_id': course_id,
                'title': rec['title'],
                'description': rec.get('description', ''),
                'difficulty_level': rec.get('difficulty_level', 3),
                'subject': rec.get('subject', ''),
                'scores': {
                    'collaborative': rec['score'],
                    'content_based': 0.0,
                    'knowledge_based': 0.0
                },
                'reasons': [rec['reason']],
                'details': rec.get('details', {})
            }
        
        # Добавить content-based
        for rec in content_recs:
            course_id = rec['course_id']
            if course_id in all_recommendations:
                all_recommendations[course_id]['scores']['content_based'] = rec['score']
                all_recommendations[course_id]['reasons'].append(rec['reason'])
            else:
                all_recommendations[course_id] = {
                    'course_id': course_id,
                    'title': rec['title'],
                    'description': rec.get('description', ''),
                    'difficulty_level': rec.get('difficulty_level', 3),
                    'subject': rec.get('subject', ''),
                    'scores': {
                        'collaborative': 0.0,
                        'content_based': rec['score'],
                        'knowledge_based': 0.0
                    },
                    'reasons': [rec['reason']],
                    'details': rec.get('details', {})
                }
        
        # Добавить knowledge-based
        for rec in knowledge_recs:
            course_id = rec['course_id']
            if course_id in all_recommendations:
                all_recommendations[course_id]['scores']['knowledge_based'] = rec['score']
                all_recommendations[course_id]['reasons'].append(rec['reason'])
            else:
                all_recommendations[course_id] = {
                    'course_id': course_id,
                    'title': rec['title'],
                    'description': rec.get('description', ''),
                    'difficulty_level': rec.get('difficulty_level', 3),
                    'subject': rec.get('subject', ''),
                    'scores': {
                        'collaborative': 0.0,
                        'content_based': 0.0,
                        'knowledge_based': rec['score']
                    },
                    'reasons': [rec['reason']],
                    'details': rec.get('details', {})
                }
        
        # 4. Вычислить финальный score для каждого курса
        final_recommendations = []
        
        for course_id, data in all_recommendations.items():
            scores = data['scores']
            
            # Взвешенная сумма с АДАПТИВНЫМИ весами
            final_score = (
                scores['collaborative'] * weights['collaborative'] +
                scores['content_based'] * weights['content_based'] +
                scores['knowledge_based'] * weights['knowledge_based']
            )
            
            # Бонус если курс рекомендован несколькими алгоритмами
            num_algorithms = sum(1 for s in scores.values() if s > 0)
            if num_algorithms >= 2:
                # Бонус пропорционален количеству алгоритмов
                consensus_bonus = 0.1 * (num_algorithms - 1)
                final_score = min(1.0, final_score * (1 + consensus_bonus))
            
            # Определить главный алгоритм
            main_algo = max(
                [('collaborative', scores['collaborative']),
                 ('content_based', scores['content_based']),
                 ('knowledge_based', scores['knowledge_based'])],
                key=lambda x: x[1]
            )[0]
            
            main_reason = next(
                (r for r in data['reasons'] if main_algo in r.lower() or 
                 ('похожим' in r.lower() and main_algo == 'collaborative') or
                 ('интересам' in r.lower() and main_algo == 'content_based') or
                 ('навык' in r.lower() and main_algo == 'knowledge_based')),
                data['reasons'][0] if data['reasons'] else "Рекомендовано для вас"
            )
            
            final_recommendations.append({
                'course_id': course_id,
                'title': data['title'],
                'description': data['description'],
                'difficulty_level': data['difficulty_level'],
                'subject': data['subject'],
                'score': round(final_score, 3),
                'algorithm': 'hybrid',
                'reason': main_reason,
                'details': {
                    'collaborative_score': round(scores['collaborative'], 3),
                    'content_based_score': round(scores['content_based'], 3),
                    'knowledge_based_score': round(scores['knowledge_based'], 3),
                    'all_reasons': data['reasons'],
                    'num_algorithms': num_algorithms,
                    'weights_used': weights,
                    'main_algorithm': main_algo
                }
            })
        
        # 5. Сортировать по финальному score
        final_recommendations.sort(key=lambda x: x['score'], reverse=True)
        
        # 6. Применить diversity boost
        final_recommendations = self._diversity_boost(final_recommendations)
        
        # 7. Пересортировать после diversity boost
        final_recommendations.sort(key=lambda x: x['score'], reverse=True)
        
        # 8. Вернуть топ-N
        return final_recommendations[:top_n]
    
    def _get_fallback_recommendations(self, student_id: int, top_n: int = 10) -> List[Dict]:
        """
        Fallback: популярные курсы с учетом профиля студента
        """
        query = """
            SELECT 
                sp.age_group, sp.grade
            FROM student_profiles sp
            WHERE sp.user_id = %s
        """
        profile = self.db.execute_one(query, (student_id,))
        
        age_group = profile['age_group'] if profile else 'middle'
        grade = profile['grade'] if profile else 5
        
        query = """
            SELECT 
                c.id as course_id,
                c.title,
                c.description,
                c.difficulty_level,
                c.subject,
                COALESCE(AVG(rr.rating), 3.5) as avg_rating,
                COUNT(DISTINCT rr.student_id) as student_count,
                COUNT(rr.id) as total_ratings,
                CASE WHEN c.age_group = %s THEN 1 ELSE 0 END as age_match
            FROM courses c
            JOIN modules m ON c.id = m.course_id
            JOIN resources r ON m.id = r.module_id
            LEFT JOIN resource_ratings rr ON r.id = rr.resource_id
            WHERE c.is_published = true
            GROUP BY c.id, c.title, c.description, c.difficulty_level, c.subject, c.age_group
            HAVING COUNT(rr.id) >= 5
            ORDER BY 
                age_match DESC,
                avg_rating DESC, 
                student_count DESC
            LIMIT %s
        """
        
        courses = self.db.execute(query, (age_group, top_n))
        
        result = []
        for course in courses:
            rating_score = max(0, (course['avg_rating'] - 3.5) / 1.5)
            popularity_score = min(0.3, (course['total_ratings'] / 100.0))
            age_bonus = 0.2 if course['age_match'] else 0
            
            final_score = min(1.0, 0.5 * rating_score + popularity_score + age_bonus)
            
            result.append({
                'course_id': course['course_id'],
                'title': course['title'],
                'description': course['description'],
                'difficulty_level': course['difficulty_level'],
                'subject': course['subject'],
                'score': round(final_score, 3),
                'algorithm': 'fallback',
                'reason': f"Популярный курс (★{course['avg_rating']:.1f}, {course['student_count']} учеников)",
                'details': {
                    'is_fallback': True,
                    'avg_rating': float(course['avg_rating']),
                    'student_count': course['student_count']
                }
            })
        
        return result
    
    def save_recommendations(self, student_id: int, recommendations: List[Dict]):
        """
        Сохранить рекомендации КУРСОВ в БД
        """
        # Удалить старые рекомендации (старше 7 дней)
        delete_query = """
            DELETE FROM course_recommendations
            WHERE student_id = %s
              AND created_at < NOW() - INTERVAL '7 days'
        """
        self.db.cursor.execute(delete_query, (student_id,))
        
        # Вставить новые рекомендации
        insert_query = """
            INSERT INTO course_recommendations 
                (student_id, course_id, score, reason, algorithm_type, is_viewed)
            VALUES (%s, %s, %s, %s, %s, false)
            ON CONFLICT (student_id, course_id) 
            DO UPDATE SET
                score = EXCLUDED.score,
                reason = EXCLUDED.reason,
                algorithm_type = EXCLUDED.algorithm_type,
                created_at = NOW()
        """
        
        for rec in recommendations:
            self.db.cursor.execute(
                insert_query,
                (
                    student_id,
                    rec['course_id'],
                    rec['score'],
                    rec['reason'],
                    rec.get('details', {}).get('main_algorithm', rec['algorithm'])
                )
            )
        
        self.db.conn.commit()