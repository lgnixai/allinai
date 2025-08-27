const axios = require('axios');

// 测试配置
const BASE_URL = 'http://localhost:3000'; // 根据实际情况调整
const TEST_PHONE = '17629726688'; // 测试手机号
const TEST_VERIFICATION_CODE = '123456'; // 测试验证码

// 测试用例
async function testFirstUseLogin() {
    console.log('🚀 开始测试首次登录 is_first_use 字段问题...\n');
    
    try {
        // 1. 发送登录请求
        console.log('📱 发送登录请求...');
        const loginResponse = await axios.post(`${BASE_URL}/api/user/login`, {
            phone: TEST_PHONE,
            phone_verification_code: TEST_VERIFICATION_CODE
        });
        
        console.log('✅ 登录请求成功');
        console.log('📊 响应状态:', loginResponse.status);
        console.log('📄 响应数据:', JSON.stringify(loginResponse.data, null, 2));
        
        // 2. 检查响应结构
        const { success, message, data } = loginResponse.data;
        
        if (!success) {
            console.log('❌ 登录失败:', message);
            return;
        }
        
        // 3. 检查是否包含 is_first_use 字段
        if (data && data.hasOwnProperty('is_first_use')) {
            console.log('✅ 响应中包含 is_first_use 字段');
            console.log('📊 is_first_use 值:', data.is_first_use);
            
            if (data.is_first_use === 1) {
                console.log('✅ is_first_use 值正确 (1)');
            } else if (data.is_first_use === 0) {
                console.log('❌ is_first_use 值错误 (0) - 首次登录应该是 1');
            } else {
                console.log('⚠️ is_first_use 值异常:', data.is_first_use);
            }
        } else {
            console.log('❌ 响应中缺少 is_first_use 字段');
            console.log('🔍 当前响应字段:', Object.keys(data || {}));
        }
        
        // 4. 检查其他重要字段
        const expectedFields = ['id', 'username', 'phone', 'role', 'status'];
        const missingFields = expectedFields.filter(field => !data || !data.hasOwnProperty(field));
        
        if (missingFields.length > 0) {
            console.log('⚠️ 缺少其他重要字段:', missingFields);
        } else {
            console.log('✅ 其他重要字段都存在');
        }
        
    } catch (error) {
        console.error('❌ 测试失败:', error.message);
        if (error.response) {
            console.error('📊 错误响应:', error.response.data);
        }
    }
}

// 运行测试
testFirstUseLogin();
