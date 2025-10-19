#!/bin/bash

API_URL="http://localhost:8080/api"

echo "🔹 1. Авторизация пользователя..."
RESPONSE=$(curl -s -X POST "$API_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"test2@edu.kz","password":"password123"}')

echo "Ответ: $RESPONSE"
TOKEN=$(echo $RESPONSE | jq -r '.token')
echo "🪪 Получен токен: $TOKEN"
echo

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
  echo "❌ Ошибка: токен не получен. Проверь данные авторизации."
  exit 1
fi

echo "🔹 2. Получаем рекомендации..."
curl -s "$API_URL/recommendations" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo

echo "🔹 3. Обновляем рекомендации..."
curl -s -X POST "$API_URL/recommendations/refresh" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo

echo "🔹 4. Получаем список курсов..."
curl -s "$API_URL/courses" | jq .
echo

echo "🔹 5. Получаем конкретный курс (ID=1)..."
curl -s "$API_URL/courses/1" | jq .
echo

echo "✅ Все запросы выполнены!"
