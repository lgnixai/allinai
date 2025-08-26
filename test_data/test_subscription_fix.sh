#!/bin/bash

# 测试订阅创建功能修复的脚本
# 使用方法: ./test_data/test_subscription_fix.sh

BASE_URL="http://localhost:3000/api"
TOKEN="your_auth_token_here"  # 请替换为实际的认证token

echo "=== 测试订阅创建功能修复 ==="

# 1. 创建第一个订阅
echo "1. 创建第一个订阅..."
CREATE_SUBSCRIPTION_1_RESPONSE=$(curl -s -X POST "$BASE_URL/subscriptions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "topic_name": "人工智能技术",
    "topic_description": "探索AI技术的最新发展"
  }')

echo "创建第一个订阅响应: $CREATE_SUBSCRIPTION_1_RESPONSE"

# 提取订阅ID
SUBSCRIPTION_ID_1=$(echo $CREATE_SUBSCRIPTION_1_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
echo "第一个订阅ID: $SUBSCRIPTION_ID_1"

# 2. 尝试创建相同主题的订阅（应该返回已存在的订阅）
echo "2. 尝试创建相同主题的订阅..."
CREATE_SUBSCRIPTION_2_RESPONSE=$(curl -s -X POST "$BASE_URL/subscriptions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "topic_name": "人工智能技术",
    "topic_description": "重复的主题名称"
  }')

echo "创建第二个订阅响应: $CREATE_SUBSCRIPTION_2_RESPONSE"

# 3. 创建不同主题的订阅
echo "3. 创建不同主题的订阅..."
CREATE_SUBSCRIPTION_3_RESPONSE=$(curl -s -X POST "$BASE_URL/subscriptions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "topic_name": "区块链技术",
    "topic_description": "区块链技术发展动态"
  }')

echo "创建第三个订阅响应: $CREATE_SUBSCRIPTION_3_RESPONSE"

# 4. 获取用户的所有订阅
echo "4. 获取用户的所有订阅..."
GET_SUBSCRIPTIONS_RESPONSE=$(curl -s -X GET "$BASE_URL/subscriptions" \
  -H "Authorization: Bearer $TOKEN")

echo "获取订阅列表响应: $GET_SUBSCRIPTIONS_RESPONSE"

# 5. 获取订阅下的文章
echo "5. 获取订阅下的文章..."
if [ ! -z "$SUBSCRIPTION_ID_1" ]; then
  GET_ARTICLES_RESPONSE=$(curl -s -X GET "$BASE_URL/subscriptions/$SUBSCRIPTION_ID_1/articles" \
    -H "Authorization: Bearer $TOKEN")
  
  echo "获取文章列表响应: $GET_ARTICLES_RESPONSE"
else
  echo "无法获取订阅ID，跳过文章获取测试"
fi

echo "=== 测试完成 ==="
echo "请检查："
echo "1. 第一个订阅是否创建成功"
echo "2. 第二个相同主题的订阅是否返回已存在的订阅"
echo "3. 第三个不同主题的订阅是否创建成功"
echo "4. 订阅列表是否包含所有订阅"
echo "5. 文章列表是否包含新字段"
