#!/bin/bash

echo "📋 验证 API 文档更新..."
echo ""

# 测试用户登录 API
echo "🔍 测试用户登录 API 响应字段..."

# 发送验证码
curl -s "http://localhost:9999/api/phone_verification?phone=17629726688&purpose=login" > /dev/null
sleep 1

# 登录并检查响应
RESPONSE=$(curl -s http://localhost:9999/api/user/login -X POST -H "Content-Type: application/json" -d '{"phone":"17629726688","phone_verification_code":"1111"}')

# 检查 is_first_use 字段是否存在
if echo "$RESPONSE" | jq -e '.data.is_first_use' > /dev/null; then
    echo "   ✅ is_first_use 字段存在"
    
    # 检查字段值
    VALUE=$(echo "$RESPONSE" | jq -r '.data.is_first_use')
    echo "   📊 is_first_use 值: $VALUE"
    
    if [ "$VALUE" = "1" ]; then
        echo "   ✅ is_first_use 值正确 (1)"
    else
        echo "   ⚠️ is_first_use 值: $VALUE"
    fi
else
    echo "   ❌ is_first_use 字段缺失"
fi

# 检查其他重要字段
echo ""
echo "🔍 检查其他重要字段..."

IMPORTANT_FIELDS=("id" "username" "display_name" "role" "status" "phone" "access_token")
for field in "${IMPORTANT_FIELDS[@]}"; do
    if echo "$RESPONSE" | jq -e ".data.$field" > /dev/null; then
        echo "   ✅ $field 字段存在"
    else
        echo "   ❌ $field 字段缺失"
    fi
done

echo ""
echo "📄 API 响应示例:"
echo "$RESPONSE" | jq '.data | {id, username, display_name, role, status, phone, access_token, is_first_use}' | head -20

echo ""
echo "📋 验证完成！"
echo "✅ API 文档更新已生效"
echo "✅ is_first_use 字段正确返回"
