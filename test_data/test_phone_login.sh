#!/bin/bash

# 测试手机号验证码登录功能
# 这个脚本测试新的登录系统，使用手机号+验证码而不是密码

BASE_URL="http://localhost:3000"
API_BASE_URL="http://localhost:3000/api"

echo "=== 测试手机号验证码登录功能 ==="

# 1. 测试发送验证码
echo "1. 测试发送验证码..."
PHONE="13800138000"
RESPONSE=$(curl -s -X GET "$API_BASE_URL/phone_verification?phone=$PHONE")
echo "发送验证码响应: $RESPONSE"

# 2. 测试登录（需要验证码）
echo ""
echo "2. 测试登录（需要验证码）..."
echo "注意：这个测试需要真实的验证码，请手动输入验证码进行测试"
echo "登录API: POST $API_BASE_URL/user/login"
echo "请求体: {\"phone\": \"$PHONE\", \"phone_verification_code\": \"验证码\"}"

# 3. 测试注册（需要验证码）
echo ""
echo "3. 测试注册（需要验证码）..."
echo "注册API: POST $API_BASE_URL/user/register"
echo "请求体: {\"phone\": \"新手机号\", \"phone_verification_code\": \"验证码\", \"display_name\": \"测试用户\"}"

echo ""
echo "=== 测试完成 ==="
echo "请手动测试以下功能："
echo "1. 访问 $BASE_URL/login 查看新的登录界面"
echo "2. 访问 $BASE_URL/register 查看新的注册界面"
echo "3. 验证密码字段已被移除，只有手机号和验证码字段"
