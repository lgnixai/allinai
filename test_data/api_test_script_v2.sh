#!/bin/bash

# 用户API测试脚本 v2.0
# 包含新的登录流程（自动生成access_token）和所有用户相关接口

# 配置
BASE_URL="http://localhost:3000"
TURNSTILE_TOKEN="your_turnstile_token_here"

# 测试用户信息
TEST_PHONE="13800138000"
TEST_PASSWORD="12345678"
TEST_DISPLAY_NAME="API测试用户"
TEST_SCHOOL="测试大学"
TEST_COLLEGE="测试学院"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 全局变量
ACCESS_TOKEN=""
USER_ID=""

# 测试函数
test_api() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo -e "${BLUE}测试: $description${NC}"
    echo -e "${YELLOW}请求: $method $BASE_URL$endpoint${NC}"
    
    if [ -n "$data" ]; then
        echo -e "${YELLOW}数据: $data${NC}"
    fi
    
    # 执行请求
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "$BASE_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$BASE_URL$endpoint")
    fi
    
    # 分离响应体和状态码
    http_code=$(echo "$response" | tail -n1)
    response_body=$(echo "$response" | head -n -1)
    
    echo -e "${YELLOW}状态码: $http_code${NC}"
    echo -e "${YELLOW}响应: $response_body${NC}"
    
    # 检查响应
    if echo "$response_body" | grep -q '"success":true'; then
        echo -e "${GREEN}✓ 成功${NC}"
        
        # 如果是登录请求，提取access_token和user_id
        if [ "$endpoint" = "/api/user/login" ]; then
            ACCESS_TOKEN=$(echo "$response_body" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
            USER_ID=$(echo "$response_body" | grep -o '"id":[0-9]*' | cut -d':' -f2)
            echo -e "${GREEN}提取到 access_token: $ACCESS_TOKEN${NC}"
            echo -e "${GREEN}提取到 user_id: $USER_ID${NC}"
        fi
        
        return 0
    else
        echo -e "${RED}✗ 失败${NC}"
        return 1
    fi
}

# 带认证的API测试函数
test_api_with_auth() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo -e "${BLUE}测试: $description${NC}"
    echo -e "${YELLOW}请求: $method $BASE_URL$endpoint${NC}"
    echo -e "${YELLOW}认证: Authorization=$ACCESS_TOKEN, New-Api-User=$USER_ID${NC}"
    
    if [ -n "$data" ]; then
        echo -e "${YELLOW}数据: $data${NC}"
    fi
    
    # 执行请求
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" \
            -H "Authorization: $ACCESS_TOKEN" \
            -H "New-Api-User: $USER_ID" \
            "$BASE_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" \
            -H "Content-Type: application/json" \
            -H "Authorization: $ACCESS_TOKEN" \
            -H "New-Api-User: $USER_ID" \
            -d "$data" \
            "$BASE_URL$endpoint")
    fi
    
    # 分离响应体和状态码
    http_code=$(echo "$response" | tail -n1)
    response_body=$(echo "$response" | head -n -1)
    
    echo -e "${YELLOW}状态码: $http_code${NC}"
    echo -e "${YELLOW}响应: $response_body${NC}"
    
    # 检查响应
    if echo "$response_body" | grep -q '"success":true'; then
        echo -e "${GREEN}✓ 成功${NC}"
        return 0
    else
        echo -e "${RED}✗ 失败${NC}"
        return 1
    fi
}

# 主测试流程
main() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}用户API测试脚本 v2.0${NC}"
    echo -e "${BLUE}================================${NC}"
    echo ""
    
    # 1. 发送手机验证码
    test_api "GET" "/api/phone_verification?phone=$TEST_PHONE&turnstile=$TURNSTILE_TOKEN" "" "发送手机验证码"
    echo ""
    
    # 2. 用户注册
    register_data="{\"phone\":\"$TEST_PHONE\",\"password\":\"$TEST_PASSWORD\",\"display_name\":\"$TEST_DISPLAY_NAME\",\"school\":\"$TEST_SCHOOL\",\"college\":\"$TEST_COLLEGE\",\"phone_verification_code\":\"1111\"}"
    test_api "POST" "/api/user/register" "$register_data" "用户注册"
    echo ""
    
    # 3. 用户登录（自动生成access_token）
    login_data="{\"phone\":\"$TEST_PHONE\",\"password\":\"$TEST_PASSWORD\"}"
    test_api "POST" "/api/user/login" "$login_data" "用户登录（自动生成access_token）"
    echo ""
    
    # 检查是否成功获取到认证信息
    if [ -z "$ACCESS_TOKEN" ] || [ -z "$USER_ID" ]; then
        echo -e "${RED}错误: 未能获取到access_token或user_id，无法继续测试${NC}"
        exit 1
    fi
    
    # 4. 获取用户信息
    test_api_with_auth "GET" "/api/user/self" "" "获取用户信息"
    echo ""
    
    # 5. 更新用户信息
    update_data="{\"display_name\":\"更新后的显示名称\",\"school\":\"更新后的大学\",\"college\":\"更新后的学院\"}"
    test_api_with_auth "PUT" "/api/user/self" "$update_data" "更新用户信息"
    echo ""
    
    # 6. 修改密码
    password_data="{\"original_password\":\"$TEST_PASSWORD\",\"password\":\"newpassword123\"}"
    test_api_with_auth "PUT" "/api/user/self" "$password_data" "修改密码"
    echo ""
    
    # 7. 生成新的访问令牌
    test_api_with_auth "GET" "/api/user/token" "" "生成新的访问令牌"
    echo ""
    
    # 8. 获取邀请码
    test_api_with_auth "GET" "/api/user/aff" "" "获取邀请码"
    echo ""
    
    # 9. 测试密码重置流程
    echo -e "${BLUE}测试密码重置流程${NC}"
    
    # 9.1 发送密码重置验证码
    test_api "GET" "/api/reset_password?phone=$TEST_PHONE&turnstile=$TURNSTILE_TOKEN" "" "发送密码重置验证码"
    echo ""
    
    # 9.2 验证重置验证码
    verify_data="{\"phone\":\"$TEST_PHONE\",\"token\":\"1111\"}"
    test_api "POST" "/api/user/verify_reset_code" "$verify_data" "验证重置验证码"
    echo ""
    
    # 9.3 重置密码
    reset_data="{\"phone\":\"$TEST_PHONE\",\"token\":\"1111\",\"password\":\"resetpassword123\"}"
    test_api "POST" "/api/user/reset_password" "$reset_data" "重置密码"
    echo ""
    
    # 10. 用户登出
    test_api "GET" "/api/user/logout" "" "用户登出"
    echo ""
    
    echo -e "${BLUE}================================${NC}"
    echo -e "${GREEN}所有测试完成！${NC}"
    echo -e "${BLUE}================================${NC}"
}

# 错误处理
error_handler() {
    echo -e "${RED}错误: $1${NC}"
    exit 1
}

# 检查依赖
check_dependencies() {
    if ! command -v curl &> /dev/null; then
        error_handler "curl 未安装，请先安装 curl"
    fi
    
    if ! command -v jq &> /dev/null; then
        echo -e "${YELLOW}警告: jq 未安装，JSON解析可能不完整${NC}"
    fi
}

# 显示使用说明
show_usage() {
    echo -e "${BLUE}使用说明:${NC}"
    echo "1. 确保API服务器正在运行"
    echo "2. 更新 BASE_URL 和 TURNSTILE_TOKEN 变量"
    echo "3. 运行脚本: ./api_test_script_v2.sh"
    echo ""
    echo -e "${YELLOW}注意事项:${NC}"
    echo "- 验证码已硬编码为 '1111'"
    echo "- 测试会创建新用户，请确保手机号可用"
    echo "- 测试完成后建议清理测试数据"
}

# 主程序
if [ "$1" = "--help" ] || [ "$1" = "-h" ]; then
    show_usage
    exit 0
fi

check_dependencies
main

