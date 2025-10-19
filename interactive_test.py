#!/usr/bin/env python3
"""
Интерактивное тестирование ML-сервиса
Позволяет выбирать студента и получать рекомендации
"""

import requests
import json

BASE_URL = "http://localhost:5000"

# Данные студентов из seed data
STUDENTS = {
    5: {"name": "Иван Петров", "interests": ["math", "physics", "programming"]},
    6: {"name": "Алия Нурланова", "interests": ["literature", "history", "art"]},
    7: {"name": "Дмитрий Сидоров", "interests": ["science", "biology", "chemistry"]},
    8: {"name": "Айгерим Асан", "interests": ["math", "physics", "chemistry"]},
    9: {"name": "Максим Волков", "interests": ["math", "art", "music"]},
    10: {"name": "Камила Темиргали", "interests": ["history", "geography", "literature"]},
}

def print_menu():
    print("\n" + "="*60)
    print("  🎓 ML Service - Интерактивное тестирование")
    print("="*60)
    print("\n📚 Выберите действие:\n")
    print("  1. Получить гибридные рекомендации")
    print("  2. Получить collaborative рекомендации")
    print("  3. Получить content-based рекомендации")
    print("  4. Получить knowledge-based рекомендации")
    print("  5. Обновить навык студента")
    print("  6. Сравнить все алгоритмы")
    print("  7. Проверить health")
    print("  0. Выход")
    print()

def print_students():
    print("\n👥 Доступные студенты:")
    print("-" * 60)
    for student_id, info in STUDENTS.items():
        interests = ", ".join(info["interests"])
        print(f"  {student_id}. {info['name']}")
        print(f"     Интересы: {interests}")
    print()

def get_student_choice():
    """Выбор студента"""
    print_students()
    
    while True:
        try:
            choice = int(input("Введите ID студента: "))
            if choice in STUDENTS:
                return choice
            else:
                print("❌ Неверный ID студента!")
        except ValueError:
            print("❌ Введите число!")

def get_recommendations(endpoint: str, student_id: int, top_n: int = 5):
    """Получить рекомендации"""
    url = f"{BASE_URL}/recommendations/{endpoint}"
    
    payload = {
        "student_id": student_id,
        "top_n": top_n
    }
    
    try:
        response = requests.post(url, json=payload, timeout=10)
        
        if response.status_code == 200:
            return response.json()
        else:
            print(f"❌ Ошибка: {response.status_code}")
            print(response.json())
            return None
    except requests.exceptions.ConnectionError:
        print("❌ Не удается подключиться к ML-сервису!")
        print("Убедитесь, что сервис запущен на http://localhost:5000")
        return None
    except Exception as e:
        print(f"❌ Ошибка: {e}")
        return None

def display_recommendations(recommendations, algorithm_name: str):
    """Отобразить рекомендации"""
    if not recommendations:
        print("\n⚠️  Рекомендаций не найдено")
        return
    
    print(f"\n📚 {algorithm_name} - Найдено {len(recommendations)} рекомендаций:")
    print("-" * 60)
    
    for i, rec in enumerate(recommendations, 1):
        print(f"\n{i}. {rec['title']}")
        print(f"   ⭐ Score: {rec['score']:.3f}")
        print(f"   📝 {rec['reason']}")
        
        if 'details' in rec and rec['details']:
            details = rec['details']
            print(f"   📊 Детали:")
            print(f"      Collaborative: {details.get('collaborative_score', 0):.3f}")
            print(f"      Content-Based: {details.get('content_based_score', 0):.3f}")
            print(f"      Knowledge-Based: {details.get('knowledge_based_score', 0):.3f}")

def test_hybrid():
    """Тест гибридных рекомендаций"""
    student_id = get_student_choice()
    student_name = STUDENTS[student_id]["name"]
    
    print(f"\n🔍 Получаем гибридные рекомендации для {student_name}...")
    
    recommendations = get_recommendations("hybrid", student_id)
    display_recommendations(recommendations, "Гибридные рекомендации")

def test_collaborative():
    """Тест collaborative filtering"""
    student_id = get_student_choice()
    student_name = STUDENTS[student_id]["name"]
    
    print(f"\n🔍 Получаем collaborative рекомендации для {student_name}...")
    
    recommendations = get_recommendations("collaborative", student_id)
    display_recommendations(recommendations, "Collaborative Filtering")

def test_content_based():
    """Тест content-based filtering"""
    student_id = get_student_choice()
    student_name = STUDENTS[student_id]["name"]
    
    print(f"\n🔍 Получаем content-based рекомендации для {student_name}...")
    
    recommendations = get_recommendations("content-based", student_id)
    display_recommendations(recommendations, "Content-Based Filtering")

def test_knowledge_based():
    """Тест knowledge-based filtering"""
    student_id = get_student_choice()
    student_name = STUDENTS[student_id]["name"]
    
    print(f"\n🔍 Получаем knowledge-based рекомендации для {student_name}...")
    
    recommendations = get_recommendations("knowledge-based", student_id)
    display_recommendations(recommendations, "Knowledge-Based Filtering")

def update_skill():
    """Обновление навыка"""
    student_id = get_student_choice()
    student_name = STUDENTS[student_id]["name"]
    
    print(f"\n📝 Обновление навыка для {student_name}")
    
    skill_name = input("Введите название навыка (например, algebra): ")
    
    while True:
        try:
            score = float(input("Введите результат теста (0.0 - 1.0): "))
            if 0.0 <= score <= 1.0:
                break
            else:
                print("❌ Введите число от 0.0 до 1.0")
        except ValueError:
            print("❌ Введите число!")
    
    payload = {
        "student_id": student_id,
        "skill_name": skill_name,
        "test_score": score
    }
    
    try:
        response = requests.post(
            f"{BASE_URL}/skills/update",
            json=payload,
            timeout=10
        )
        
        if response.status_code == 200:
            result = response.json()
            print(f"\n✅ {result['message']}")
            print(f"   Новый уровень: {result['new_level']:.3f}")
        else:
            print(f"❌ Ошибка: {response.status_code}")
            print(response.json())
    except Exception as e:
        print(f"❌ Ошибка: {e}")

def compare_algorithms():
    """Сравнить все алгоритмы"""
    student_id = get_student_choice()
    student_name = STUDENTS[student_id]["name"]
    
    print(f"\n🔬 Сравнение алгоритмов для {student_name}...")
    print("="*60)
    
    algorithms = [
        ("collaborative", "Collaborative Filtering"),
        ("content-based", "Content-Based Filtering"),
        ("knowledge-based", "Knowledge-Based Filtering"),
        ("hybrid", "Hybrid (Combined)")
    ]
    
    for endpoint, name in algorithms:
        print(f"\n{'='*60}")
        print(f"📊 {name}")
        print('='*60)
        
        recommendations = get_recommendations(endpoint, student_id, top_n=3)
        
        if recommendations:
            for i, rec in enumerate(recommendations, 1):
                print(f"{i}. {rec['title']} (Score: {rec['score']:.3f})")
        else:
            print("⚠️  Нет рекомендаций")
        
        input("\nНажмите Enter для продолжения...")

def test_health():
    """Проверка здоровья сервиса"""
    print("\n🏥 Проверка состояния ML-сервиса...")
    
    try:
        # Basic health
        response = requests.get(f"{BASE_URL}/", timeout=5)
        print(f"\n✅ Базовая проверка: {response.status_code}")
        print(json.dumps(response.json(), indent=2, ensure_ascii=False))
        
        # Detailed health
        response = requests.get(f"{BASE_URL}/health", timeout=5)
        print(f"\n✅ Детальная проверка: {response.status_code}")
        print(json.dumps(response.json(), indent=2, ensure_ascii=False))
        
    except requests.exceptions.ConnectionError:
        print("❌ Не удается подключиться к сервису!")
    except Exception as e:
        print(f"❌ Ошибка: {e}")

def main():
    """Главный цикл"""
    
    while True:
        print_menu()
        
        try:
            choice = input("Ваш выбор: ").strip()
            
            if choice == "1":
                test_hybrid()
            elif choice == "2":
                test_collaborative()
            elif choice == "3":
                test_content_based()
            elif choice == "4":
                test_knowledge_based()
            elif choice == "5":
                update_skill()
            elif choice == "6":
                compare_algorithms()
            elif choice == "7":
                test_health()
            elif choice == "0":
                print("\n👋 До свидания!")
                break
            else:
                print("❌ Неверный выбор!")
            
            input("\n\nНажмите Enter для продолжения...")
            
        except KeyboardInterrupt:
            print("\n\n👋 До свидания!")
            break
        except Exception as e:
            print(f"\n❌ Ошибка: {e}")
            input("\nНажмите Enter для продолжения...")

if __name__ == "__main__":
    main()