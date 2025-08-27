const axios = require('axios');

// 测试配置
const BASE_URL = 'http://localhost:9999';
const TEST_PHONE = '17629726688';
const TEST_VERIFICATION_CODE = '1111';

async function debugTest() {
    console.log('🔍 调试 is_first_use 字段问题...\n');
    
    try {
        // 1. 发送登录请求
        console.log('📱 发送登录请求...');
        const response = await axios.post(`${BASE_URL}/api/user/login`, {
            phone: TEST_PHONE,
            phone_verification_code: TEST_VERIFICATION_CODE
        });
        
        console.log('✅ 登录请求成功');
        console.log('📊 响应状态:', response.status);
        
        const { success, data } = response.data;
        
        if (success && data) {
            console.log('\n📄 完整响应数据:');
            console.log(JSON.stringify(response.data, null, 2));
            
            console.log('\n🔍 字段分析:');
            console.log('   is_first_use 字段存在:', data.hasOwnProperty('is_first_use'));
            console.log('   is_first_use 值:', data.is_first_use);
            console.log('   is_first_use 类型:', typeof data.is_first_use);
            
            // 检查所有字段
            console.log('\n📋 所有返回字段:');
            Object.keys(data).forEach(key => {
                console.log(`   ${key}: ${data[key]} (${typeof data[key]})`);
            });
            
            // 检查是否有其他相关字段
            const relatedFields = ['first_use', 'firstUse', 'isFirstUse', 'first_use_flag'];
            relatedFields.forEach(field => {
                if (data.hasOwnProperty(field)) {
                    console.log(`   ⚠️ 发现相关字段 ${field}: ${data[field]}`);
                }
            });
            
        } else {
            console.log('❌ 登录失败:', response.data.message);
        }
        
    } catch (error) {
        console.error('❌ 测试失败:', error.message);
        if (error.response) {
            console.error('📊 错误响应:', error.response.data);
        }
    }
}

debugTest();
