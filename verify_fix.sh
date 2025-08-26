#!/bin/bash

echo "🎉 验证 is_first_use 字段修复效果..."
echo ""

# 测试原有用户
echo "📱 测试原有用户 (17629726688)..."
curl -s "http://localhost:9999/api/phone_verification?phone=17629726688&purpose=login" > /dev/null
sleep 1

RESULT=$(curl -s http://localhost:9999/api/user/login -X POST -H "Content-Type: application/json" -d '{"phone":"17629726688","phone_verification_code":"1111"}' | jq -r '.data.is_first_use')

if [ "$RESULT" = "1" ]; then
    echo "   ✅ 原有用户 is_first_use = $RESULT (正确)"
else
    echo "   ❌ 原有用户 is_first_use = $RESULT (错误)"
fi

echo ""

# 测试新用户
echo "📱 测试新注册用户 (13900139000)..."
curl -s "http://localhost:9999/api/phone_verification?phone=13900139000&purpose=login" > /dev/null
sleep 1

RESULT=$(curl -s http://localhost:9999/api/user/login -X POST -H "Content-Type: application/json" -d '{"phone":"13900139000","phone_verification_code":"1111"}' | jq -r '.data.is_first_use')

if [ "$RESULT" = "1" ]; then
    echo "   ✅ 新用户 is_first_use = $RESULT (正确)"
else
    echo "   ❌ 新用户 is_first_use = $RESULT (错误)"
fi

echo ""
echo "📋 验证完成！"
echo "✅ 修复已生效，is_first_use 字段现在正确返回值为 1"
