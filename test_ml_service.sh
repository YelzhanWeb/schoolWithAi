#!/bin/bash

# Цвета для вывода
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:5000"

echo -e "${GREEN}=================================${NC}"
echo -e "${GREEN}🧪 Testing ML Service${NC}"
echo -e "${GREEN}=================================${NC}\n"

# 1. Health Check
echo -e "${YELLOW}1️⃣  Testing Health Check...${NC}"
curl -s "${BASE_URL}/" | jq .
echo -e "\n"

# 2. Detailed Health
echo -e "${YELLOW}2️⃣  Testing Detailed Health...${NC}"
curl -s "${BASE_URL}/health" | jq .
echo -e "\n"

# 3. Hybrid Recommendations
echo -e "${YELLOW}3️⃣  Testing Hybrid Recommendations (Student ID: 5)...${NC}"
curl -s -X POST "${BASE_URL}/recommendations/hybrid" \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 5,
    "top_n": 5
  }' | jq .
echo -e "\n"

# 4. Collaborative Filtering
echo -e "${YELLOW}4️⃣  Testing Collaborative Filtering (Student ID: 5)...${NC}"
curl -s -X POST "${BASE_URL}/recommendations/collaborative" \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 5,
    "top_n": 5
  }' | jq .
echo -e "\n"

# 5. Content-Based Filtering
echo -e "${YELLOW}5️⃣  Testing Content-Based Filtering (Student ID: 6)...${NC}"
curl -s -X POST "${BASE_URL}/recommendations/content-based" \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 6,
    "top_n": 5
  }' | jq .
echo -e "\n"

# 6. Knowledge-Based Filtering
echo -e "${YELLOW}6️⃣  Testing Knowledge-Based Filtering (Student ID: 7)...${NC}"
curl -s -X POST "${BASE_URL}/recommendations/knowledge-based" \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 7,
    "top_n": 5
  }' | jq .
echo -e "\n"

# 7. Update Skill
echo -e "${YELLOW}7️⃣  Testing Skill Update (Student ID: 5, algebra)...${NC}"
curl -s -X POST "${BASE_URL}/skills/update" \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": 5,
    "skill_name": "algebra",
    "test_score": 0.9
  }' | jq .
echo -e "\n"

echo -e "${GREEN}=================================${NC}"
echo -e "${GREEN}✅ Testing Complete!${NC}"
echo -e "${GREEN}=================================${NC}"