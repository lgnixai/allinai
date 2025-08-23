#!/bin/bash

echo "=================================="
echo "  聊天消息显示测试"
echo "=================================="

# 设置变量
BASE_URL="http://localhost:4000"
USER_ID="7"
ACCESS_TOKEN="GiLjdWi5bjR1uXu43k4kbXZnH1ZBqpMf"

echo "[INFO] 测试聊天消息显示功能..."

# 1. 获取现有会话列表
echo "[TEST] 获取会话列表"
SESSIONS_RESPONSE=$(curl -s -H "UserID: $USER_ID" -H "Authorization: Bearer $ACCESS_TOKEN" \
  "$BASE_URL/api/chat_sessions/user/sessions?page=1&size=50")

if echo "$SESSIONS_RESPONSE" | jq -e '.success' > /dev/null; then
    echo "[SUCCESS] 获取会话列表成功"
    
    # 获取第一个有消息的会话
    SESSION_ID=$(echo "$SESSIONS_RESPONSE" | jq -r '.data.data[0].session_id')
    SESSION_TOPIC=$(echo "$SESSIONS_RESPONSE" | jq -r '.data.data[0].topic')
    MESSAGE_COUNT=$(echo "$SESSIONS_RESPONSE" | jq -r '.data.data[0].message_count')
    
    echo "[INFO] 选择会话: $SESSION_TOPIC (ID: $SESSION_ID, 消息数: $MESSAGE_COUNT)"
    
    # 2. 获取该会话的消息
    echo "[TEST] 获取会话消息"
    MESSAGES_RESPONSE=$(curl -s -H "UserID: $USER_ID" -H "Authorization: Bearer $ACCESS_TOKEN" \
      "$BASE_URL/api/chat_messages/session/$SESSION_ID?page=1&size=100")
    
    if echo "$MESSAGES_RESPONSE" | jq -e '.success' > /dev/null; then
        echo "[SUCCESS] 获取消息成功"
        
        # 显示消息详情
        MESSAGE_COUNT_ACTUAL=$(echo "$MESSAGES_RESPONSE" | jq '.data.data | length')
        echo "[INFO] 实际消息数: $MESSAGE_COUNT_ACTUAL"
        
        # 显示前几条消息
        echo "[INFO] 消息内容预览:"
        echo "$MESSAGES_RESPONSE" | jq -r '.data.data[] | "  \(.role): \(.content)"' | head -6
        
        # 检查时间戳
        echo "[INFO] 时间戳检查:"
        echo "$MESSAGES_RESPONSE" | jq -r '.data.data[] | "  \(.message_id): created_time=\(.created_time)"' | head -4
        
    else
        echo "[ERROR] 获取消息失败"
        echo "$MESSAGES_RESPONSE"
        exit 1
    fi
    
else
    echo "[ERROR] 获取会话列表失败"
    echo "$SESSIONS_RESPONSE"
    exit 1
fi

echo "=================================="
echo "  测试完成"
echo "=================================="



