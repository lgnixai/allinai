-- 用户API测试数据生成脚本 v2.0
-- 包含新的用户字段：学校、学院、手机号，以及access_token

-- 清理现有测试数据
DELETE FROM users WHERE phone LIKE '138%' OR phone LIKE '139%' OR phone LIKE '137%';
DELETE FROM tokens WHERE name LIKE '%测试%';

-- 插入测试用户数据
INSERT INTO users (
    username, 
    password, 
    display_name, 
    role, 
    status, 
    email, 
    phone, 
    school, 
    college, 
    group_name, 
    quota, 
    used_quota, 
    request_count, 
    aff_code, 
    aff_count, 
    aff_quota, 
    aff_history_quota, 
    inviter_id, 
    access_token,
    created_time,
    updated_time
) VALUES 
-- 普通用户1
('user_8000', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVEFDa', '测试用户1', 1, 1, 'test1@example.com', '13800138000', '清华大学', '计算机学院', 'default', 1000000, 50000, 100, 'TEST001', 2, 50000, 100000, 0, 'test_access_token_001', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- 普通用户2
('user_8001', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVEFDa', '测试用户2', 1, 1, 'test2@example.com', '13800138001', '北京大学', '信息学院', 'default', 500000, 25000, 50, 'TEST002', 1, 25000, 50000, 1, 'test_access_token_002', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- 普通用户3
('user_8002', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVEFDa', '测试用户3', 1, 1, 'test3@example.com', '13800138002', '复旦大学', '软件学院', 'vip', 2000000, 100000, 200, 'TEST003', 0, 0, 0, 0, 'test_access_token_003', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- 管理员用户
('admin_8003', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVEFDa', '管理员用户', 2, 1, 'admin@example.com', '13800138003', '浙江大学', '管理学院', 'admin', 5000000, 200000, 500, 'ADMIN001', 5, 100000, 200000, 0, 'admin_access_token_001', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- 超级管理员
('root_8004', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVEFDa', '超级管理员', 3, 1, 'root@example.com', '13800138004', '上海交通大学', '电子学院', 'root', 10000000, 500000, 1000, 'ROOT001', 10, 200000, 500000, 0, 'root_access_token_001', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- 禁用用户
('disabled_8005', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVEFDa', '禁用用户', 1, 2, 'disabled@example.com', '13800138005', '武汉大学', '物理学院', 'default', 100000, 0, 0, 'DISABLED001', 0, 0, 0, 0, 'disabled_access_token_001', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- VIP用户
('vip_8006', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVEFDa', 'VIP用户', 1, 1, 'vip@example.com', '13800138006', '华中科技大学', '机械学院', 'vip', 3000000, 150000, 300, 'VIP001', 3, 75000, 150000, 2, 'vip_access_token_001', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- 新注册用户（无access_token）
('new_8007', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVEFDa', '新注册用户', 1, 1, 'new@example.com', '13800138007', '西安交通大学', '化学学院', 'default', 100000, 0, 0, 'NEW001', 0, 0, 0, 0, NULL, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- 企业用户
('enterprise_8008', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVEFDa', '企业用户', 1, 1, 'enterprise@example.com', '13800138008', '北京理工大学', '经济学院', 'enterprise', 5000000, 300000, 600, 'ENTERPRISE001', 8, 150000, 300000, 0, 'enterprise_access_token_001', UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),

-- 学生用户
('student_8009', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVEFDa', '学生用户', 1, 1, 'student@example.com', '13800138009', '南京大学', '文学院', 'student', 200000, 10000, 20, 'STUDENT001', 0, 0, 0, 0, 'student_access_token_001', UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

-- 插入测试令牌数据
INSERT INTO tokens (
    user_id,
    name,
    key,
    created_time,
    accessed_time,
    expired_time,
    remain_quota,
    unlimited_quota,
    model_limits_enabled,
    group_name
) VALUES 
-- 用户1的令牌
(LAST_INSERT_ID()-9, '测试用户1的初始令牌', 'test_token_key_001', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), -1, 500000, 1, 0, 'auto'),

-- 用户2的令牌
(LAST_INSERT_ID()-8, '测试用户2的初始令牌', 'test_token_key_002', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), -1, 250000, 1, 0, 'auto'),

-- 用户3的令牌
(LAST_INSERT_ID()-7, '测试用户3的初始令牌', 'test_token_key_003', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), -1, 1000000, 1, 0, 'auto'),

-- 管理员令牌
(LAST_INSERT_ID()-6, '管理员用户的初始令牌', 'admin_token_key_001', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), -1, 2500000, 1, 0, 'auto'),

-- 超级管理员令牌
(LAST_INSERT_ID()-5, '超级管理员的初始令牌', 'root_token_key_001', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), -1, 5000000, 1, 0, 'auto'),

-- VIP用户令牌
(LAST_INSERT_ID()-3, 'VIP用户的初始令牌', 'vip_token_key_001', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), -1, 1500000, 1, 0, 'auto'),

-- 企业用户令牌
(LAST_INSERT_ID()-1, '企业用户的初始令牌', 'enterprise_token_key_001', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), -1, 2500000, 1, 0, 'auto'),

-- 学生用户令牌
(LAST_INSERT_ID(), '学生用户的初始令牌', 'student_token_key_001', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), -1, 100000, 1, 0, 'auto');

-- 显示插入的用户数据
SELECT 
    id,
    username,
    display_name,
    phone,
    school,
    college,
    role,
    status,
    group_name,
    access_token,
    quota,
    used_quota,
    aff_code
FROM users 
WHERE phone LIKE '138%' 
ORDER BY id;

-- 显示插入的令牌数据
SELECT 
    t.id,
    t.user_id,
    u.username,
    t.name,
    t.key,
    t.remain_quota,
    t.unlimited_quota
FROM tokens t
JOIN users u ON t.user_id = u.id
WHERE u.phone LIKE '138%'
ORDER BY t.id;

