from typing import List, Dict
from .collaborative_courses import CollaborativeFilteringCourses
from .content_based_courses import ContentBasedFilteringCourses
from .knowledge_based_courses import KnowledgeBasedFilteringCourses


class HybridRecommenderCourses:
    """
    Гибридная рекомендательная система для КУРСОВ
    Объединяет: Collaborative + Content-Based + Knowledge-Based
    """
    
    def __init__(self, db):
        self.db = db
        self.collaborative = CollaborativeFilteringCourses(db)
        self.content_based = ContentBasedFilteringCourses(db)
        self.knowledge_based = KnowledgeBasedFilteringCourses(db)
        
        # Веса для каждого алгоритма
        self.weights = {
            'collaborative': 0.4,
            'content_based': 0.3,
            'knowledge_based': 0.3
        }
    
    def recommend(self, student_id: int, top_n: int = 10) -> List[Dict]:
        """
        Генерация гибридных рекомендаций КУРСОВ
        """
        # 1. Получить рекомендации от каждого алгоритма
        collab_recs = self.collaborative.recommend(student_id, top_n=top_n)
        content_recs = self.content_based.recommend(student_id, top_n=top_n)
        knowledge_recs = self.knowledge_based.recommend(student_id, top_n=top_n)
        
        # 2. Объединить все рекомендации
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
                'reasons': [rec['reason']]
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
                    'reasons': [rec['reason']]
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
                    'reasons': [rec['reason']]
                }
        
        # 3. Вычислить финальный score для каждого курса
        final_recommendations = []
        
        for course_id, data in all_recommendations.items():
            scores = data['scores']
            
            # Взвешенная сумма
            final_score = (
                scores['collaborative'] * self.weights['collaborative'] +
                scores['content_based'] * self.weights['content_based'] +
                scores['knowledge_based'] * self.weights['knowledge_based']
            )
            
            # Бонус если курс рекомендован несколькими алгоритмами
            num_algorithms = sum(1 for s in scores.values() if s > 0)
            if num_algorithms > 1:
                final_score *= 1.15  # +15% бонус
            
            # Ограничить 0-1
            final_score = min(final_score, 1.0)
            
            # Определить главную причину рекомендации
            main_reason = data['reasons'][0] if data['reasons'] else "Рекомендовано для вас"
            
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
                    'num_algorithms': num_algorithms
                }
            })
        
        # 4. Сортировать по финальному score
        final_recommendations.sort(key=lambda x: x['score'], reverse=True)
        
        # 5. Вернуть топ-N
        return final_recommendations[:top_n]
    
    def save_recommendations(self, student_id: int, recommendations: List[Dict]):
        """
        Сохранить рекомендации КУРСОВ в БД
        """
        # Удалить старые рекомендации
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
                    rec['algorithm']
                )
            )
        
        self.db.conn.commit()