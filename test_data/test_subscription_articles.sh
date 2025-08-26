#!/bin/bash

# 测试订阅文章新字段功能的脚本
# 使用方法: ./test_subscription_articles.sh

BASE_URL="http://localhost:3000/api"
TOKEN="your_auth_token_here"  # 请替换为实际的认证token

echo "=== 测试订阅文章新字段功能 ==="

# 1. 创建订阅
echo "1. 创建订阅..."
CREATE_SUBSCRIPTION_RESPONSE=$(curl -s -X POST "$BASE_URL/subscriptions" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "topic_name": "人工智能技术",
    "topic_description": "探索AI技术的最新发展"
  }')

echo "创建订阅响应: $CREATE_SUBSCRIPTION_RESPONSE"

# 提取订阅ID
SUBSCRIPTION_ID=$(echo $CREATE_SUBSCRIPTION_RESPONSE | grep -o '"id":[0-9]*' | cut -d':' -f2)
echo "订阅ID: $SUBSCRIPTION_ID"

# 等待模拟文章生成
echo "等待模拟文章生成..."
sleep 3

# 2. 获取订阅下的文章
echo "2. 获取订阅下的文章..."
ARTICLES_RESPONSE=$(curl -s -X GET "$BASE_URL/subscriptions/$SUBSCRIPTION_ID/articles" \
  -H "Authorization: Bearer $TOKEN")

echo "文章列表响应: $ARTICLES_RESPONSE"

# 3. 创建包含新字段的文章（需要管理员权限）
echo "3. 创建包含新字段的文章..."
CREATE_ARTICLE_RESPONSE=$(curl -s -X POST "$BASE_URL/subscription_articles" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "subscription_id": '$SUBSCRIPTION_ID',
    "title": "测试文章 - 新字段功能",
    "summary": "这是一篇测试新字段功能的文章",
    "content": "文章内容...",
    "author": "测试作者",
    "published_at": "2024-01-01T00:00:00Z",
    "article_url": "https://example.com/test",
    "key_points": "1. 测试重点提炼功能\n2. 验证新字段的正确性\n3. 确保API正常工作",
    "journal_name": "测试期刊",
    "read_count": 100,
    "citation_count": 5,
    "rating": 8.5
  }')

echo "创建文章响应: $CREATE_ARTICLE_RESPONSE"

# 4. 获取所有订阅文章
echo "4. 获取所有订阅文章..."
ALL_ARTICLES_RESPONSE=$(curl -s -X GET "$BASE_URL/subscriptions/articles" \
  -H "Authorization: Bearer $TOKEN")

echo "所有文章响应: $ALL_ARTICLES_RESPONSE"

echo "=== 测试完成 ==="
echo "请检查响应中是否包含新字段: key_points, journal_name, read_count, citation_count, rating"
