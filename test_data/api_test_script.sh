#!/bin/bash

# API测试脚本
# 用于测试用户相关的API接口

BASE_URL="http://localhost:3000"
TURNSTILE_TOKEN="test_token"

echo "=== 用户API测试脚本 ==="
echo "基础URL: $BASE_URL"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试函数
test_api() {
    local name="$1"
    local method="$2"
    local url="$3"
    local data="$4"
    local expected_status="$5"
    
    echo -e "${YELLOW}测试: $name${NC}"
    echo "URL: $url"
    echo "方法: $method"
    if [ ! -z "$data" ]; then
        echo "数据: $data"
    fi
    echo ""
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "$url")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" -H "Content-Type: application/json" -d "$data" "$url")
    fi
    
    # 分离响应体和状态码
    body=$(echo "$response" | head -n -1)
    status=$(echo "$response" | tail -n 1)
    
    echo "响应状态: $status"
    echo "响应内容: $body"
    
    if [ "$status" = "$expected_status" ]; then
        echo -e "${GREEN}✓ 测试通过${NC}"
    else
        echo -e "${RED}✗ 测试失败 (期望: $expected_status, 实际: $status)${NC}"
    fi
    echo "----------------------------------------"
    echo ""
}

# 1. 发送手机验证码
echo "1. 测试发送手机验证码"
test_api "发送手机验证码" "GET" "$BASE_URL/api/phone_verification?phone=13800138001&turnstile=$TURNSTILE_TOKEN" "" "200"

# 2. 用户注册
echo "2. 测试用户注册"
register_data='{
  "phone": "13800138001",
  "phone_verification_code": "1111",
  "password": "12345678",
  "display_name": "测试用户",
  "school": "测试大学",
  "college": "计算机学院"
}'
test_api "用户注册" "POST" "$BASE_URL/api/user/register?turnstile=$TURNSTILE_TOKEN" "$register_data" "200"

# 3. 用户登录
echo "3. 测试用户登录"
login_data='{
  "phone": "13800138001",
  "password": "12345678"
}'
test_api "用户登录" "POST" "$BASE_URL/api/user/login?turnstile=$TURNSTILE_TOKEN" "$login_data" "200"

# 4. 发送密码重置验证码
echo "4. 测试发送密码重置验证码"
test_api "发送密码重置验证码" "GET" "$BASE_URL/api/reset_password?phone=13800138001&turnstile=$TURNSTILE_TOKEN" "" "200"

# 5. 重置密码
echo "5. 测试重置密码"
reset_data='{
  "phone": "13800138001",
  "token": "1111"
}'
test_api "重置密码" "POST" "$BASE_URL/api/user/reset" "$reset_data" "200"

# 6. 使用新密码登录
echo "6. 测试使用新密码登录"
new_login_data='{
  "phone": "13800138001",
  "password": "87654321"
}'
test_api "使用新密码登录" "POST" "$BASE_URL/api/user/login?turnstile=$TURNSTILE_TOKEN" "$new_login_data" "200"

# 7. 获取用户信息（需要认证）
echo "7. 测试获取用户信息（需要认证）"
echo "注意: 这个测试需要有效的访问令牌"
test_api "获取用户信息" "GET" "$BASE_URL/api/user/self" "" "401"

# 8. 更新用户信息（需要认证）
echo "8. 测试更新用户信息（需要认证）"
update_data='{
  "username": "user_8001",
  "display_name": "更新后的显示名称",
  "school": "更新后的大学",
  "college": "更新后的学院",
  "phone": "13800138001"
}'
test_api "更新用户信息" "PUT" "$BASE_URL/api/user/self" "$update_data" "401"

# 9. 用户登出
echo "9. 测试用户登出"
test_api "用户登出" "GET" "$BASE_URL/api/user/logout" "" "200"

echo "=== 测试完成 ==="
echo ""
echo "注意事项:"
echo "1. 需要先启动服务器"
echo "2. 确保数据库已正确配置"
echo "3. 验证码目前固定为 '1111'"
echo "4. 需要有效的访问令牌才能测试认证接口"
echo "5. 建议使用Postman进行更详细的测试"
