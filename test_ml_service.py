#!/usr/bin/env python3
"""
ML Service Testing Script
Тестирует все endpoints ML-сервиса
"""

import requests
import json
from typing import Dict, Any

BASE_URL = "http://localhost:5000"

def print_header(text: str):
    print("\n" + "="*50)
    print(f"  {text}")
    print("="*50)

def print_test(name: str):
    print(f"\n🧪 {name}")
    print("-" * 50)

def print_result(data: Any):
    print(json.dumps(data, indent=2, ensure_ascii=False))

def test_health():
    """Тест базовой проверки здоровья"""
    print_test("Health Check")
    response = requests.get(f"{BASE_URL}/")
    print(f"Status: {response.status_code}")
    print_result(response.json())
    return response.status_code == 200

def test_detailed_health():
    """Тест детальной проверки здоровья"""
    print_test("Detailed Health Check")
    response = requests.get(f"{BASE_URL}/health")
    print(f"Status: {response.status_code}")
    print_result(response.json())
    return response.status_code == 200

def test_hybrid_recommendations(student_id: int = 5, top_n: int = 5):
    """Тест гибридных рекомендаций"""
    print_test(f"Hybrid Recommendations (Student {student_id})")
    
    payload = {
        "student_id": student_id,
        "top_n": top_n
    }
    
    response = requests.post(
        f"{BASE_URL}/recommendations/hybrid",
        json=payload
    )
    
    print(f"Status: {response.status_code}")
    
    if response.status_code == 200:
        recommendations = response.json()
        print(f"\n📚 Найдено {len(recommendations)} рекомендаций:")
        
        for i, rec in enumerate(recommendations, 1):
            print(f"\n{i}. {rec['title']}")
            print(f"   Score: {rec['score']:.3f}")
            print(f"   Algorithm: {rec['algorithm']}")
            print(f"   Reason: {rec['reason']}")
            
            if 'details' in rec and rec['details']:
                details = rec['details']
                print(f"   Details:")
                print(f"     - Collaborative: {details.get('collaborative_score', 0):.3f}")
                print(f"     - Content-Based: {details.get('content_based_score', 0):.3f}")
                print(f"     - Knowledge-Based: {details.get('knowledge_based_score', 0):.3f}")
        
        return True
    else:
        print("❌ Error:")
        print_result(response.json())
        return False

def test_collaborative(student_id: int = 5, top_n: int = 5):
    """Тест коллаборативной фильтрации"""
    print_test(f"Collaborative Filtering (Student {student_id})")
    
    payload = {
        "student_id": student_id,
        "top_n": top_n
    }
    
    response = requests.post(
        f"{BASE_URL}/recommendations/collaborative",
        json=payload
    )
    
    print(f"Status: {response.status_code}")
    
    if response.status_code == 200:
        recommendations = response.json()
        print(f"\n📚 Найдено {len(recommendations)} рекомендаций:")
        
        for i, rec in enumerate(recommendations, 1):
            print(f"{i}. {rec['title']} (Score: {rec['score']:.3f})")
            print(f"   {rec['reason']}")
        
        return True
    else:
        print("❌ Error:")
        print_result(response.json())
        return False

def test_content_based(student_id: int = 6, top_n: int = 5):
    """Тест контентной фильтрации"""
    print_test(f"Content-Based Filtering (Student {student_id})")
    
    payload = {
        "student_id": student_id,
        "top_n": top_n
    }
    
    response = requests.post(
        f"{BASE_URL}/recommendations/content-based",
        json=payload
    )
    
    print(f"Status: {response.status_code}")
    
    if response.status_code == 200:
        recommendations = response.json()
        print(f"\n📚 Найдено {len(recommendations)} рекомендаций:")
        
        for i, rec in enumerate(recommendations, 1):
            print(f"{i}. {rec['title']} (Score: {rec['score']:.3f})")
            print(f"   {rec['reason']}")
        
        return True
    else:
        print("❌ Error:")
        print_result(response.json())
        return False

def test_knowledge_based(student_id: int = 7, top_n: int = 5):
    """Тест знаниевой фильтрации"""
    print_test(f"Knowledge-Based Filtering (Student {student_id})")
    
    payload = {
        "student_id": student_id,
        "top_n": top_n
    }
    
    response = requests.post(
        f"{BASE_URL}/recommendations/knowledge-based",
        json=payload
    )
    
    print(f"Status: {response.status_code}")
    
    if response.status_code == 200:
        recommendations = response.json()
        print(f"\n📚 Найдено {len(recommendations)} рекомендаций:")
        
        for i, rec in enumerate(recommendations, 1):
            print(f"{i}. {rec['title']} (Score: {rec['score']:.3f})")
            print(f"   {rec['reason']}")
        
        return True
    else:
        print("❌ Error:")
        print_result(response.json())
        return False

def test_skill_update(student_id: int = 5, skill_name: str = "algebra", test_score: float = 0.9):
    """Тест обновления навыка"""
    print_test(f"Update Skill (Student {student_id}, Skill: {skill_name})")
    
    payload = {
        "student_id": student_id,
        "skill_name": skill_name,
        "test_score": test_score
    }
    
    response = requests.post(
        f"{BASE_URL}/skills/update",
        json=payload
    )
    
    print(f"Status: {response.status_code}")
    
    if response.status_code == 200:
        result = response.json()
        print(f"\n✅ Skill updated!")
        print(f"   Student: {result['student_id']}")
        print(f"   Skill: {result['skill_name']}")
        print(f"   New Level: {result['new_level']:.3f}")
        print(f"   Message: {result['message']}")
        return True
    else:
        print("❌ Error:")
        print_result(response.json())
        return False

def run_all_tests():
    """Запустить все тесты"""
    print_header("🚀 ML Service Testing")
    
    results = {
        "Health Check": test_health(),
        "Detailed Health": test_detailed_health(),
        "Hybrid Recommendations": test_hybrid_recommendations(5, 5),
        "Collaborative Filtering": test_collaborative(5, 5),
        "Content-Based Filtering": test_content_based(6, 5),
        "Knowledge-Based Filtering": test_knowledge_based(7, 5),
        "Skill Update": test_skill_update(5, "algebra", 0.9),
    }
    
    # Итоговый отчет
    print_header("📊 Test Results")
    
    passed = sum(results.values())
    total = len(results)
    
    for test_name, result in results.items():
        status = "✅ PASSED" if result else "❌ FAILED"
        print(f"{status} - {test_name}")
    
    print(f"\n{'='*50}")
    print(f"Total: {passed}/{total} tests passed")
    
    if passed == total:
        print("🎉 All tests passed!")
    else:
        print("⚠️  Some tests failed")
    
    print(f"{'='*50}\n")
    
    return passed == total

if __name__ == "__main__":
    try:
        success = run_all_tests()
        exit(0 if success else 1)
    except requests.exceptions.ConnectionError:
        print("\n❌ Cannot connect to ML Service!")
        print("Make sure the service is running on http://localhost:5000")
        exit(1)
    except Exception as e:
        print(f"\n❌ Unexpected error: {e}")
        exit(1)