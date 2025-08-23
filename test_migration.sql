-- 测试用户表新字段的SQL脚本
-- 运行迁移后执行此脚本来验证字段是否正确添加

-- 检查表结构
DESCRIBE users;

-- 检查新字段是否存在
SELECT 
    COLUMN_NAME, 
    DATA_TYPE, 
    IS_NULLABLE, 
    COLUMN_DEFAULT, 
    COLUMN_COMMENT
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_NAME = 'users' 
    AND COLUMN_NAME IN ('school', 'college', 'phone');

-- 检查索引是否存在
SHOW INDEX FROM users WHERE Key_name = 'idx_users_phone';

-- 测试插入数据
INSERT INTO users (username, password, display_name, school, college, phone) 
VALUES ('test_user', 'hashed_password', 'Test User', '测试大学', '计算机学院', '13800138000');

-- 验证数据插入
SELECT id, username, school, college, phone FROM users WHERE username = 'test_user';

-- 清理测试数据
DELETE FROM users WHERE username = 'test_user';
