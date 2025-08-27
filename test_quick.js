const axios = require('axios');

// 测试配置
const BASE_URL = 'http://localhost:9999';
const TEST_PHONE = '17629726688';
const TEST_VERIFICATION_CODE = '1111';

async function quickTest() {
    console.log('🚀 快速测试 is_first_use 字段修复效果...\n');
    
    try {
        const response = await axios.post(`${BASE_URL}/api/user/login`, {
            phone: TEST_PHONE,
            phone_verification_code: TEST_VERIFICATION_CODE
        });
        
        const { success, data } = response.data;
        
        if (success && data) {
            console.log('✅ 登录成功');
            console.log('📊 is_first_use 字段存在:', data.hasOwnProperty('is_first_use'));
            console.log('📊 is_first_use 值:', data.is_first_use);
            
            if (data.is_first_use === 1) {
                console.log('🎉 修复成功！is_first_use 值正确 (1)');
            } else {
                console.log('⚠️ is_first_use 值仍不正确:', data.is_first_use);
            }
        } else {
            console.log('❌ 登录失败:', response.data.message);
        }
        
    } catch (error) {
        console.error('❌ 测试失败:', error.message);
    }
}

quickTest();
