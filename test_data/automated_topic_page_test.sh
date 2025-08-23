#!/bin/bash

# 话题页面自动化测试脚本
# 测试话题页面的完整功能，包括前端页面和后端API

# 配置
BASE_URL="http://localhost:4000"
FRONTEND_URL="http://localhost:5173"
ACCESS_TOKEN="your_access_token_here"
USER_ID="7"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 测试计数器
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
    ((PASSED_TESTS++))
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
    ((FAILED_TESTS++))
}

print_test() {
    echo -e "${PURPLE}[TEST]${NC} $1"
    ((TOTAL_TESTS++))
}

print_header() {
    echo -e "${CYAN}================================${NC}"
    echo -e "${CYAN}  话题页面自动化测试开始        ${NC}"
    echo -e "${CYAN}================================${NC}"
}

print_footer() {
    echo -e "${CYAN}================================${NC}"
    echo -e "${CYAN}  测试结果统计                  ${NC}"
    echo -e "${CYAN}================================${NC}"
    echo -e "总测试数: ${TOTAL_TESTS}"
    echo -e "通过: ${GREEN}${PASSED_TESTS}${NC}"
    echo -e "失败: ${RED}${FAILED_TESTS}${NC}"
    echo -e "成功率: $((PASSED_TESTS * 100 / TOTAL_TESTS))%"
}

# 检查服务是否运行
check_services() {
    print_info "检查服务状态..."
    
    # 检查后端服务
    print_test "检查后端服务"
    if curl -s "$BASE_URL/api/status" > /dev/null; then
        print_success "后端服务正在运行"
    else
        print_error "后端服务未运行"
        return 1
    fi
    
    # 检查前端服务
    print_test "检查前端服务"
    if curl -s "$FRONTEND_URL" > /dev/null; then
        print_success "前端服务正在运行"
    else
        print_warning "前端服务未运行，跳过前端测试"
        FRONTEND_AVAILABLE=false
    fi
}

# 测试用户认证
test_authentication() {
    print_info "测试用户认证..."
    
    print_test "测试用户登录"
    LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/api/user/login" \
        -H "Content-Type: application/json" \
        -d '{
            "phone": "13800138001",
            "password": "12345678"
        }')
    
    if echo "$LOGIN_RESPONSE" | grep -q "success.*true"; then
        print_success "用户登录成功"
        # 提取access_token
        ACCESS_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
        if [ -n "$ACCESS_TOKEN" ]; then
            print_success "获取到access_token: ${ACCESS_TOKEN:0:20}..."
        else
            print_error "未获取到access_token"
            return 1
        fi
    else
        print_error "用户登录失败: $LOGIN_RESPONSE"
        return 1
    fi
}

# 测试会话管理API
test_session_apis() {
    print_info "测试会话管理API..."
    
    # 创建测试会话
    print_test "创建测试会话"
    SESSION_RESPONSE=$(curl -s -X POST "$BASE_URL/api/chat_sessions/" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "UserID: $USER_ID" \
        -H "Content-Type: application/json" \
        -d '{
            "topic": "自动化测试会话",
            "model": "gpt-3.5-turbo",
            "channel_id": 1
        }')
    
    if echo "$SESSION_RESPONSE" | grep -q "success.*true"; then
        print_success "创建测试会话成功"
        SESSION_ID=$(echo "$SESSION_RESPONSE" | grep -o '"session_id":"[^"]*"' | cut -d'"' -f4)
        print_info "会话ID: $SESSION_ID"
    else
        print_error "创建测试会话失败: $SESSION_RESPONSE"
        return 1
    fi
    
    # 获取会话列表
    print_test "获取会话列表"
    LIST_RESPONSE=$(curl -s -X GET "$BASE_URL/api/chat_sessions/user/sessions?page=1&size=10" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "UserID: $USER_ID")
    
    if echo "$LIST_RESPONSE" | grep -q "success.*true"; then
        print_success "获取会话列表成功"
        SESSION_COUNT=$(echo "$LIST_RESPONSE" | grep -o '"total":[0-9]*' | cut -d':' -f2)
        print_info "总会话数: $SESSION_COUNT"
    else
        print_error "获取会话列表失败: $LIST_RESPONSE"
    fi
    
    # 获取会话详情
    print_test "获取会话详情"
    DETAIL_RESPONSE=$(curl -s -X GET "$BASE_URL/api/chat_sessions/$SESSION_ID" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "UserID: $USER_ID")
    
    if echo "$DETAIL_RESPONSE" | grep -q "success.*true"; then
        print_success "获取会话详情成功"
    else
        print_error "获取会话详情失败: $DETAIL_RESPONSE"
    fi
    
    # 搜索会话
    print_test "搜索会话"
    SEARCH_RESPONSE=$(curl -s -X GET "$BASE_URL/api/chat_sessions/search?keyword=自动化测试&page=1&size=10" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "UserID: $USER_ID")
    
    if echo "$SEARCH_RESPONSE" | grep -q "success.*true"; then
        print_success "搜索会话成功"
    else
        print_error "搜索会话失败: $SEARCH_RESPONSE"
    fi
    
    # 获取用户统计
    print_test "获取用户统计"
    STATS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/chat_sessions/user/stats" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "UserID: $USER_ID")
    
    if echo "$STATS_RESPONSE" | grep -q "success.*true"; then
        print_success "获取用户统计成功"
    else
        print_error "获取用户统计失败: $STATS_RESPONSE"
    fi
}

# 测试消息管理API
test_message_apis() {
    print_info "测试消息管理API..."
    
    # 创建用户消息
    print_test "创建用户消息"
    USER_MESSAGE_ID="auto-test-user-$(date +%s)"
    USER_MSG_RESPONSE=$(curl -s -X POST "$BASE_URL/api/chat_messages/" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "UserID: $USER_ID" \
        -H "Content-Type: application/json" \
        -d "{
            \"session_id\": \"$SESSION_ID\",
            \"message_id\": \"$USER_MESSAGE_ID\",
            \"role\": \"user\",
            \"content\": \"这是一个自动化测试消息\",
            \"tokens\": 15,
            \"cost\": 0.0003,
            \"status\": 1,
            \"error_msg\": \"\"
        }")
    
    if echo "$USER_MSG_RESPONSE" | grep -q "success.*true"; then
        print_success "创建用户消息成功"
    else
        print_error "创建用户消息失败: $USER_MSG_RESPONSE"
    fi
    
    # 创建AI回复消息
    print_test "创建AI回复消息"
    ASSISTANT_MESSAGE_ID="auto-test-assistant-$(date +%s)"
    ASSISTANT_MSG_RESPONSE=$(curl -s -X POST "$BASE_URL/api/chat_messages/" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "UserID: $USER_ID" \
        -H "Content-Type: application/json" \
        -d "{
            \"session_id\": \"$SESSION_ID\",
            \"message_id\": \"$ASSISTANT_MESSAGE_ID\",
            \"role\": \"assistant\",
            \"content\": \"这是AI的自动化测试回复消息\",
            \"tokens\": 25,
            \"cost\": 0.0005,
            \"status\": 1,
            \"error_msg\": \"\"
        }")
    
    if echo "$ASSISTANT_MSG_RESPONSE" | grep -q "success.*true"; then
        print_success "创建AI回复消息成功"
    else
        print_error "创建AI回复消息失败: $ASSISTANT_MSG_RESPONSE"
    fi
    
    # 获取会话消息列表
    print_test "获取会话消息列表"
    MESSAGES_RESPONSE=$(curl -s -X GET "$BASE_URL/api/chat_messages/session/$SESSION_ID?page=1&size=20" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "UserID: $USER_ID")
    
    if echo "$MESSAGES_RESPONSE" | grep -q "success.*true"; then
        print_success "获取会话消息列表成功"
        MESSAGE_COUNT=$(echo "$MESSAGES_RESPONSE" | grep -o '"total":[0-9]*' | cut -d':' -f2)
        print_info "会话消息数: $MESSAGE_COUNT"
    else
        print_error "获取会话消息列表失败: $MESSAGES_RESPONSE"
    fi
    
    # 搜索消息
    print_test "搜索消息"
    SEARCH_MSG_RESPONSE=$(curl -s -X GET "$BASE_URL/api/chat_messages/search?keyword=自动化测试&page=1&size=20" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "UserID: $USER_ID")
    
    if echo "$SEARCH_MSG_RESPONSE" | grep -q "success.*true"; then
        print_success "搜索消息成功"
    else
        print_error "搜索消息失败: $SEARCH_MSG_RESPONSE"
    fi
    
    # 获取消息统计
    print_test "获取消息统计"
    MSG_STATS_RESPONSE=$(curl -s -X GET "$BASE_URL/api/chat_messages/session/$SESSION_ID/stats" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -H "UserID: $USER_ID")
    
    if echo "$MSG_STATS_RESPONSE" | grep -q "success.*true"; then
        print_success "获取消息统计成功"
    else
        print_error "获取消息统计失败: $MSG_STATS_RESPONSE"
    fi
}

# 测试前端页面访问
test_frontend_access() {
    if [ "$FRONTEND_AVAILABLE" = false ]; then
        print_warning "跳过前端测试（前端服务未运行）"
        return 0
    fi
    
    print_info "测试前端页面访问..."
    
    # 测试话题页面访问
    print_test "测试话题页面访问"
    FRONTEND_RESPONSE=$(curl -s -I "$FRONTEND_URL/console/topic" | head -1)
    
    if echo "$FRONTEND_RESPONSE" | grep -q "200\|302"; then
        print_success "话题页面可以访问"
    else
        print_error "话题页面访问失败: $FRONTEND_RESPONSE"
    fi
    
    # 测试主页访问
    print_test "测试主页访问"
    MAIN_RESPONSE=$(curl -s -I "$FRONTEND_URL" | head -1)
    
    if echo "$MAIN_RESPONSE" | grep -q "200"; then
        print_success "主页可以访问"
    else
        print_error "主页访问失败: $MAIN_RESPONSE"
    fi
}

# 测试数据库连接
test_database_connection() {
    print_info "测试数据库连接..."
    
    print_test "检查数据库文件"
    if [ -f "sqlite.db" ]; then
        print_success "数据库文件存在"
        
        # 检查表是否存在
        print_test "检查chat_sessions表"
        SESSIONS_TABLE=$(sqlite3 sqlite.db ".tables" | grep "chat_sessions")
        if [ -n "$SESSIONS_TABLE" ]; then
            print_success "chat_sessions表存在"
        else
            print_error "chat_sessions表不存在"
        fi
        
        print_test "检查chat_messages表"
        MESSAGES_TABLE=$(sqlite3 sqlite.db ".tables" | grep "chat_messages")
        if [ -n "$MESSAGES_TABLE" ]; then
            print_success "chat_messages表存在"
        else
            print_error "chat_messages表不存在"
        fi
        
        # 检查数据
        print_test "检查会话数据"
        SESSION_COUNT=$(sqlite3 sqlite.db "SELECT COUNT(*) FROM chat_sessions;" 2>/dev/null)
        if [ -n "$SESSION_COUNT" ] && [ "$SESSION_COUNT" -ge 0 ]; then
            print_success "会话数据正常，共 $SESSION_COUNT 条记录"
        else
            print_error "会话数据异常"
        fi
        
        print_test "检查消息数据"
        MESSAGE_COUNT=$(sqlite3 sqlite.db "SELECT COUNT(*) FROM chat_messages;" 2>/dev/null)
        if [ -n "$MESSAGE_COUNT" ] && [ "$MESSAGE_COUNT" -ge 0 ]; then
            print_success "消息数据正常，共 $MESSAGE_COUNT 条记录"
        else
            print_error "消息数据异常"
        fi
    else
        print_error "数据库文件不存在"
    fi
}

# 清理测试数据
cleanup_test_data() {
    print_info "清理测试数据..."
    
    if [ -n "$SESSION_ID" ]; then
        print_test "删除测试会话"
        DELETE_RESPONSE=$(curl -s -X DELETE "$BASE_URL/api/chat_sessions/$SESSION_ID" \
            -H "Authorization: Bearer $ACCESS_TOKEN" \
            -H "UserID: $USER_ID")
        
        if echo "$DELETE_RESPONSE" | grep -q "success.*true"; then
            print_success "测试会话删除成功"
        else
            print_warning "测试会话删除失败: $DELETE_RESPONSE"
        fi
    fi
}

# 主测试函数
main() {
    print_header
    
    # 初始化变量
    FRONTEND_AVAILABLE=true
    SESSION_ID=""
    
    # 运行测试
    check_services || exit 1
    test_authentication || exit 1
    test_session_apis || exit 1
    test_message_apis || exit 1
    test_frontend_access
    test_database_connection
    
    # 清理测试数据
    cleanup_test_data
    
    # 显示测试结果
    print_footer
    
    # 返回结果
    if [ $FAILED_TESTS -eq 0 ]; then
        print_success "所有测试通过！"
        exit 0
    else
        print_error "有 $FAILED_TESTS 个测试失败"
        exit 1
    fi
}

# 运行主函数
main "$@"

