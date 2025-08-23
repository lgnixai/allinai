-- 测试手机号登录功能的SQL脚本

-- 1. 检查用户表结构
DESCRIBE users;

-- 2. 检查手机号字段约束
SELECT 
    COLUMN_NAME, 
    IS_NULLABLE, 
    COLUMN_DEFAULT, 
    COLUMN_TYPE,
    COLUMN_KEY
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_NAME = 'users' 
    AND COLUMN_NAME = 'phone';

-- 3. 检查手机号唯一索引
SHOW INDEX FROM users WHERE Key_name LIKE '%phone%';

-- 4. 测试插入用户数据（手机号登录）
INSERT INTO users (
    username, 
    password, 
    display_name, 
    phone, 
    school, 
    college,
    role,
    status,
    quota
) VALUES (
    'user_8888',
    '$2a$10$hashed_password_here', -- 实际使用时需要正确的哈希密码
    '测试用户',
    '13800138000',
    '测试大学',
    '计算机学院',
    1,
    1,
    1000000
);

-- 5. 验证数据插入
SELECT id, username, display_name, phone, school, college FROM users WHERE phone = '13800138000';

-- 6. 测试手机号唯一性约束
-- 尝试插入重复手机号（应该失败）
INSERT INTO users (
    username, 
    password, 
    display_name, 
    phone, 
    role,
    status
) VALUES (
    'user_9999',
    '$2a$10$hashed_password_here',
    '重复手机号用户',
    '13800138000', -- 重复的手机号
    1,
    1
);

-- 7. 测试手机号格式验证
-- 尝试插入无效手机号（应该失败）
INSERT INTO users (
    username, 
    password, 
    display_name, 
    phone, 
    role,
    status
) VALUES (
    'user_invalid',
    '$2a$10$hashed_password_here',
    '无效手机号用户',
    '123', -- 无效的手机号格式
    1,
    1
);

-- 8. 测试手机号查询
SELECT * FROM users WHERE phone = '13800138000';

-- 9. 清理测试数据
DELETE FROM users WHERE phone = '13800138000';
DELETE FROM users WHERE username = 'user_9999';
DELETE FROM users WHERE username = 'user_invalid';

-- 10. 验证清理结果
SELECT COUNT(*) as remaining_users FROM users WHERE phone IN ('13800138000', '123');
