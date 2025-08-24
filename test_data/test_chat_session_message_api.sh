#!/bin/bash

# 聊天会话和消息API测试脚本
# 基于两表结构的API测试

# 配置
BASE_URL="http://localhost:3000"
ACCESS_TOKEN="your_access_token_here"
USER_ID="1"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查服务是否运行
check_service() {
    print_info "检查服务是否运行..."
    if curl -s "$BASE_URL/api/status" > /dev/null; then
        print_success "服务正在运行"
        return 0
    else
        print_error "服务未运行，请先启动服务"
        return 1
    fi
}

# 测试会话管理API
test_session_apis() {
    print_info "开始测试会话管理API..."
    
    # 1. 创建会话
    print_info "1. 测试创建会话"
    SESSION_RESPONSE=$(curl -s -X POST "$BASE_URL/api/chat_sessions/" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "New-Api-User: $USER_ID" \
        -H "Content-Type: application/json" \
        -d '{
            "topic": "API测试会话",
            "model": "gpt-3.5-turbo",
            "channel_id": 1
        }')
    
    if echo "$SESSION_RESPONSE" | grep -q "success.*true"; then
        print_success "创建会话成功"
        SESSION_ID=$(echo "$SESSION_RESPONSE" | grep -o '"session_id":"[^"]*"' | cut -d'"' -f4)
        print_info "会话ID: $SESSION_ID"
    else
        print_error "创建会话失败: $SESSION_RESPONSE"
        return 1
    fi
    
    # 2. 获取会话详情
    print_info "2. 测试获取会话详情"
    GET_SESSION_RESPONSE=$(curl -s -X GET "$BASE_URL/api/chat_sessions/$SESSION_ID" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "New-Api-User: $USER_ID")
    
    if echo "$GET_SESSION_RESPONSE" | grep -q "success.*true"; then
        print_success "获取会话详情成功"
    else
        print_error "获取会话详情失败: $GET_SESSION_RESPONSE"
    fi
    
    # 3. 更新会话
    print_info "3. 测试更新会话"
    UPDATE_SESSION_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/chat_sessions/$SESSION_ID" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "New-Api-User: $USER_ID" \
        -H "Content-Type: application/json" \
        -d '{
            "topic": "API测试会话-已更新"
        }')
    
    if echo "$UPDATE_SESSION_RESPONSE" | grep -q "success.*true"; then
        print_success "更新会话成功"
    else
        print_error "更新会话失败: $UPDATE_SESSION_RESPONSE"
    fi
    
    # 4. 获取用户会话列表
    print_info "4. 测试获取用户会话列表"
    LIST_SESSIONS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/chat_sessions/user/sessions?page=1&size=10" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "New-Api-User: $USER_ID")
    
    if echo "$LIST_SESSIONS_RESPONSE" | grep -q "success.*true"; then
        print_success "获取会话列表成功"
    else
        print_error "获取会话列表失败: $LIST_SESSIONS_RESPONSE"
    fi
    
    # 5. 获取用户会话统计
    print_info "5. 测试获取用户会话统计"
    STATS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/chat_sessions/user/stats" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "New-Api-User: $USER_ID")
    
    if echo "$STATS_RESPONSE" | grep -q "success.*true"; then
        print_success "获取会话统计成功"
    else
        print_error "获取会话统计失败: $STATS_RESPONSE"
    fi
    
    # 6. 获取模型使用统计
    print_info "6. 测试获取模型使用统计"
    MODEL_STATS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/chat_sessions/user/model_stats" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "New-Api-User: $USER_ID")
    
    if echo "$MODEL_STATS_RESPONSE" | grep -q "success.*true"; then
        print_success "获取模型统计成功"
    else
        print_error "获取模型统计失败: $MODEL_STATS_RESPONSE"
    fi
    
    # 7. 搜索会话
    print_info "7. 测试搜索会话"
    SEARCH_RESPONSE=$(curl -s -X GET "$BASE_URL/api/chat_sessions/search?keyword=API测试&page=1&size=10" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "New-Api-User: $USER_ID")
    
    if echo "$SEARCH_RESPONSE" | grep -q "success.*true"; then
        print_success "搜索会话成功"
    else
        print_error "搜索会话失败: $SEARCH_RESPONSE"
    fi
    
    return 0
}

# 测试消息管理API
test_message_apis() {
    print_info "开始测试消息管理API..."
    
    # 获取一个会话ID用于测试
    SESSION_RESPONSE=$(curl -s -X GET "$BASE_URL/api/chat_sessions/user/sessions?page=1&size=1" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "New-Api-User: $USER_ID")
    
    SESSION_ID=$(echo "$SESSION_RESPONSE" | grep -o '"session_id":"[^"]*"' | head -1 | cut -d'"' -f4)
    
    if [ -z "$SESSION_ID" ]; then
        print_error "无法获取会话ID，请先创建会话"
        return 1
    fi
    
    print_info "使用会话ID: $SESSION_ID"
    
    # 1. 创建用户消息
    print_info "1. 测试创建用户消息"
    USER_MESSAGE_ID="test-user-msg-$(date +%s)"
    CREATE_USER_MSG_RESPONSE=$(curl -s -X POST "$BASE_URL/api/chat_messages/" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "New-Api-User: $USER_ID" \
        -H "Content-Type: application/json" \
        -d "{
            \"session_id\": \"$SESSION_ID\",
            \"message_id\": \"$USER_MESSAGE_ID\",
            \"role\": \"user\",
            \"content\": \"这是一条测试用户消息\",
            \"tokens\": 15,
            \"cost\": 0.0003,
            \"status\": 1,
            \"error_msg\": \"\"
        }")
    
    if echo "$CREATE_USER_MSG_RESPONSE" | grep -q "success.*true"; then
        print_success "创建用户消息成功"
    else
        print_error "创建用户消息失败: $CREATE_USER_MSG_RESPONSE"
    fi
    
    # 2. 创建助手消息
    print_info "2. 测试创建助手消息"
    ASSISTANT_MESSAGE_ID="test-assistant-msg-$(date +%s)"
    CREATE_ASSISTANT_MSG_RESPONSE=$(curl -s -X POST "$BASE_URL/api/chat_messages/" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "New-Api-User: $USER_ID" \
        -H "Content-Type: application/json" \
        -d "{
            \"session_id\": \"$SESSION_ID\",
            \"message_id\": \"$ASSISTANT_MESSAGE_ID\",
            \"role\": \"assistant\",
            \"content\": \"这是一条测试助手回复消息\",
            \"tokens\": 25,
            \"cost\": 0.0005,
            \"status\": 1,
            \"error_msg\": \"\"
        }")
    
    if echo "$CREATE_ASSISTANT_MSG_RESPONSE" | grep -q "success.*true"; then
        print_success "创建助手消息成功"
        MESSAGE_ID=$(echo "$CREATE_ASSISTANT_MSG_RESPONSE" | grep -o '"message_id":"[^"]*"' | cut -d'"' -f4)
    else
        print_error "创建助手消息失败: $CREATE_ASSISTANT_MSG_RESPONSE"
        return 1
    fi
    
    # 3. 获取消息详情
    print_info "3. 测试获取消息详情"
    GET_MSG_RESPONSE=$(curl -s -X GET "$BASE_URL/api/chat_messages/$MESSAGE_ID" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "New-Api-User: $USER_ID")
    
    if echo "$GET_MSG_RESPONSE" | grep -q "success.*true"; then
        print_success "获取消息详情成功"
    else
        print_error "获取消息详情失败: $GET_MSG_RESPONSE"
    fi
    
    # 4. 更新消息
    print_info "4. 测试更新消息"
    UPDATE_MSG_RESPONSE=$(curl -s -X PUT "$BASE_URL/api/chat_messages/$MESSAGE_ID" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "New-Api-User: $USER_ID" \
        -H "Content-Type: application/json" \
        -d '{
            "content": "这是一条更新后的测试助手回复消息",
            "tokens": 30,
            "cost": 0.0006,
            "status": 1,
            "error_msg": ""
        }')
    
    if echo "$UPDATE_MSG_RESPONSE" | grep -q "success.*true"; then
        print_success "更新消息成功"
    else
        print_error "更新消息失败: $UPDATE_MSG_RESPONSE"
    fi
    
    # 5. 获取会话消息列表
    print_info "5. 测试获取会话消息列表"
    LIST_MSGS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/chat_messages/session/$SESSION_ID?page=1&size=20" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "New-Api-User: $USER_ID")
    
    if echo "$LIST_MSGS_RESPONSE" | grep -q "success.*true"; then
        print_success "获取会话消息列表成功"
    else
        print_error "获取会话消息列表失败: $LIST_MSGS_RESPONSE"
    fi
    
    # 6. 获取用户所有消息
    print_info "6. 测试获取用户所有消息"
    USER_MSGS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/chat_messages/user/messages?page=1&size=20" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "New-Api-User: $USER_ID")
    
    if echo "$USER_MSGS_RESPONSE" | grep -q "success.*true"; then
        print_success "获取用户所有消息成功"
    else
        print_error "获取用户所有消息失败: $USER_MSGS_RESPONSE"
    fi
    
    # 7. 搜索消息
    print_info "7. 测试搜索消息"
    SEARCH_MSGS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/chat_messages/search?keyword=测试&page=1&size=20" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "New-Api-User: $USER_ID")
    
    if echo "$SEARCH_MSGS_RESPONSE" | grep -q "success.*true"; then
        print_success "搜索消息成功"
    else
        print_error "搜索消息失败: $SEARCH_MSGS_RESPONSE"
    fi
    
    # 8. 获取会话消息统计
    print_info "8. 测试获取会话消息统计"
    MSG_STATS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/chat_messages/session/$SESSION_ID/stats" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "New-Api-User: $USER_ID")
    
    if echo "$MSG_STATS_RESPONSE" | grep -q "success.*true"; then
        print_success "获取会话消息统计成功"
    else
        print_error "获取会话消息统计失败: $MSG_STATS_RESPONSE"
    fi
    
    # 9. 获取用户消息统计
    print_info "9. 测试获取用户消息统计"
    USER_MSG_STATS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/chat_messages/user/stats" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "New-Api-User: $USER_ID")
    
    if echo "$USER_MSG_STATS_RESPONSE" | grep -q "success.*true"; then
        print_success "获取用户消息统计成功"
    else
        print_error "获取用户消息统计失败: $USER_MSG_STATS_RESPONSE"
    fi
    
    # 10. 删除消息
    print_info "10. 测试删除消息"
    DELETE_MSG_RESPONSE=$(curl -s -X DELETE "$BASE_URL/api/chat_messages/$MESSAGE_ID" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "New-Api-User: $USER_ID")
    
    if echo "$DELETE_MSG_RESPONSE" | grep -q "success.*true"; then
        print_success "删除消息成功"
    else
        print_error "删除消息失败: $DELETE_MSG_RESPONSE"
    fi
    
    return 0
}

# 主函数
main() {
    print_info "开始测试聊天会话和消息API..."
    
    # 检查服务
    if ! check_service; then
        exit 1
    fi
    
    # 检查认证信息
    if [ "$ACCESS_TOKEN" = "your_access_token_here" ]; then
        print_warning "请先设置正确的ACCESS_TOKEN"
        print_info "请编辑脚本中的ACCESS_TOKEN变量"
        exit 1
    fi
    
    # 测试会话API
    if test_session_apis; then
        print_success "会话管理API测试完成"
    else
        print_error "会话管理API测试失败"
        exit 1
    fi
    
    # 测试消息API
    if test_message_apis; then
        print_success "消息管理API测试完成"
    else
        print_error "消息管理API测试失败"
        exit 1
    fi
    
    print_success "所有API测试完成！"
}

# 运行主函数
main "$@"





