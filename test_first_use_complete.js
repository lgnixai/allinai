const axios = require('axios');
const sqlite3 = require('sqlite3').verbose();
const path = require('path');

// 测试配置
const BASE_URL = 'http://localhost:3000'; // 根据实际情况调整
const TEST_PHONE = '17629726688'; // 测试手机号
const TEST_VERIFICATION_CODE = '123456'; // 测试验证码
const DB_PATH = path.join(__dirname, 'one-api.db'); // 数据库路径

// 数据库操作函数
function queryDatabase(phone) {
    return new Promise((resolve, reject) => {
        const db = new sqlite3.Database(DB_PATH);
        
        const query = `
            SELECT id, username, phone, is_first_use, role, status, school, college 
            FROM users 
            WHERE phone = ?
        `;
        
        db.get(query, [phone], (err, row) => {
            db.close();
            if (err) {
                reject(err);
            } else {
                resolve(row);
            }
        });
    });
}

// 更新用户 is_first_use 字段
function updateUserFirstUse(phone, isFirstUse) {
    return new Promise((resolve, reject) => {
        const db = new sqlite3.Database(DB_PATH);
        
        const query = `UPDATE users SET is_first_use = ? WHERE phone = ?`;
        
        db.run(query, [isFirstUse, phone], function(err) {
            db.close();
            if (err) {
                reject(err);
            } else {
                resolve(this.changes);
            }
        });
    });
}

// 测试用例
async function testFirstUseComplete() {
    console.log('🚀 开始完整测试首次登录 is_first_use 字段问题...\n');
    
    try {
        // 1. 检查数据库中的用户信息
        console.log('🔍 检查数据库中的用户信息...');
        const dbUser = await queryDatabase(TEST_PHONE);
        
        if (!dbUser) {
            console.log('❌ 数据库中未找到该用户');
            return;
        }
        
        console.log('✅ 数据库用户信息:');
        console.log('   ID:', dbUser.id);
        console.log('   用户名:', dbUser.username);
        console.log('   手机号:', dbUser.phone);
        console.log('   数据库中的 is_first_use:', dbUser.is_first_use);
        console.log('   角色:', dbUser.role);
        console.log('   状态:', dbUser.status);
        console.log('   学校:', dbUser.school);
        console.log('   学院:', dbUser.college);
        console.log('');
        
        // 2. 如果数据库中的 is_first_use 不是 1，则更新为 1
        if (dbUser.is_first_use !== 1) {
            console.log('🔄 更新数据库中的 is_first_use 为 1...');
            const updateResult = await updateUserFirstUse(TEST_PHONE, 1);
            console.log('✅ 更新完成，影响行数:', updateResult);
            
            // 重新查询确认更新
            const updatedUser = await queryDatabase(TEST_PHONE);
            console.log('✅ 更新后的 is_first_use:', updatedUser.is_first_use);
            console.log('');
        }
        
        // 3. 发送登录请求
        console.log('📱 发送登录请求...');
        const loginResponse = await axios.post(`${BASE_URL}/api/user/login`, {
            phone: TEST_PHONE,
            phone_verification_code: TEST_VERIFICATION_CODE
        });
        
        console.log('✅ 登录请求成功');
        console.log('📊 响应状态:', loginResponse.status);
        
        // 4. 检查响应结构
        const { success, message, data } = loginResponse.data;
        
        if (!success) {
            console.log('❌ 登录失败:', message);
            return;
        }
        
        console.log('📄 响应数据:');
        console.log(JSON.stringify(loginResponse.data, null, 2));
        console.log('');
        
        // 5. 检查是否包含 is_first_use 字段
        if (data && data.hasOwnProperty('is_first_use')) {
            console.log('✅ 响应中包含 is_first_use 字段');
            console.log('📊 API 返回的 is_first_use 值:', data.is_first_use);
            console.log('📊 数据库中的 is_first_use 值:', dbUser.is_first_use);
            
            if (data.is_first_use === 1) {
                console.log('✅ is_first_use 值正确 (1)');
            } else if (data.is_first_use === 0) {
                console.log('❌ is_first_use 值错误 (0) - 首次登录应该是 1');
            } else {
                console.log('⚠️ is_first_use 值异常:', data.is_first_use);
            }
            
            // 检查数据库和API返回的值是否一致
            if (data.is_first_use === dbUser.is_first_use) {
                console.log('✅ 数据库和API返回的 is_first_use 值一致');
            } else {
                console.log('❌ 数据库和API返回的 is_first_use 值不一致');
            }
        } else {
            console.log('❌ 响应中缺少 is_first_use 字段');
            console.log('🔍 当前响应字段:', Object.keys(data || {}));
        }
        
        // 6. 检查其他重要字段
        const expectedFields = ['id', 'username', 'phone', 'role', 'status', 'school', 'college'];
        const missingFields = expectedFields.filter(field => !data || !data.hasOwnProperty(field));
        
        if (missingFields.length > 0) {
            console.log('⚠️ 缺少其他重要字段:', missingFields);
        } else {
            console.log('✅ 其他重要字段都存在');
        }
        
        // 7. 测试总结
        console.log('\n📋 测试总结:');
        console.log('   数据库检查: ✅');
        console.log('   登录请求: ✅');
        console.log('   字段完整性: ' + (data && data.hasOwnProperty('is_first_use') ? '✅' : '❌'));
        console.log('   字段值正确性: ' + (data && data.is_first_use === 1 ? '✅' : '❌'));
        
    } catch (error) {
        console.error('❌ 测试失败:', error.message);
        if (error.response) {
            console.error('📊 错误响应:', error.response.data);
        }
    }
}

// 运行测试
testFirstUseComplete();
