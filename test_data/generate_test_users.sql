-- 测试用户数据生成脚本
-- 用于生成测试用户数据

-- 清理现有测试数据
DELETE FROM users WHERE phone LIKE '138%' AND phone != '13800138000';

-- 生成测试用户数据
INSERT INTO users (
    username,
    password,
    display_name,
    phone,
    email,
    school,
    college,
    role,
    status,
    quota,
    used_quota,
    request_count,
    `group`,
    aff_code,
    aff_count,
    aff_quota,
    aff_history_quota,
    inviter_id,
    created_at,
    updated_at
) VALUES 
-- 普通用户
('user_8001', '$2a$10$hashed_password_here', '张三', '13800138001', 'zhangsan@test.com', '清华大学', '计算机科学与技术学院', 1, 1, 1000000, 50000, 100, 'default', 'ABCD', 0, 0, 0, 0, NOW(), NOW()),
('user_8002', '$2a$10$hashed_password_here', '李四', '13800138002', 'lisi@test.com', '北京大学', '信息科学技术学院', 1, 1, 2000000, 150000, 200, 'default', 'EFGH', 2, 100000, 50000, 0, NOW(), NOW()),
('user_8003', '$2a$10$hashed_password_here', '王五', '13800138003', 'wangwu@test.com', '复旦大学', '软件学院', 1, 1, 500000, 75000, 50, 'default', 'IJKL', 1, 50000, 25000, 0, NOW(), NOW()),
('user_8004', '$2a$10$hashed_password_here', '赵六', '13800138004', 'zhaoliu@test.com', '上海交通大学', '电子信息与电气工程学院', 1, 1, 1500000, 300000, 300, 'default', 'MNOP', 0, 0, 0, 0, NOW(), NOW()),
('user_8005', '$2a$10$hashed_password_here', '钱七', '13800138005', 'qianqi@test.com', '浙江大学', '计算机学院', 1, 1, 800000, 120000, 80, 'default', 'QRST', 3, 150000, 75000, 0, NOW(), NOW()),

-- 管理员用户
('admin_001', '$2a$10$hashed_password_here', '管理员A', '13800138006', 'admin1@test.com', '系统管理', '技术部', 2, 1, 5000000, 500000, 1000, 'admin', 'UVWX', 10, 500000, 250000, 0, NOW(), NOW()),
('admin_002', '$2a$10$hashed_password_here', '管理员B', '13800138007', 'admin2@test.com', '系统管理', '运营部', 2, 1, 3000000, 200000, 800, 'admin', 'YZAB', 5, 200000, 100000, 0, NOW(), NOW()),

-- 超级管理员
('root_admin', '$2a$10$hashed_password_here', '超级管理员', '13800138008', 'root@test.com', '系统管理', '技术部', 3, 1, 10000000, 1000000, 2000, 'admin', 'CDEF', 20, 1000000, 500000, 0, NOW(), NOW()),

-- 禁用用户
('disabled_user', '$2a$10$hashed_password_here', '禁用用户', '13800138009', 'disabled@test.com', '测试大学', '测试学院', 1, 0, 100000, 0, 0, 'default', 'GHIJ', 0, 0, 0, 0, NOW(), NOW()),

-- 高额度用户
('vip_user', '$2a$10$hashed_password_here', 'VIP用户', '13800138010', 'vip@test.com', 'VIP大学', 'VIP学院', 1, 1, 10000000, 2000000, 500, 'vip', 'KLMN', 15, 1000000, 500000, 0, NOW(), NOW());

-- 生成测试令牌数据
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
    `group`,
    status
) VALUES 
(1, '张三的测试令牌', 'sk-test-token-001', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), -1, 500000, 0, 0, 'default', 1),
(2, '李四的测试令牌', 'sk-test-token-002', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), -1, 1000000, 0, 0, 'default', 1),
(3, '王五的测试令牌', 'sk-test-token-003', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), -1, 200000, 0, 0, 'default', 1),
(4, '赵六的测试令牌', 'sk-test-token-004', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), -1, 800000, 0, 0, 'default', 1),
(5, '钱七的测试令牌', 'sk-test-token-005', UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), -1, 400000, 0, 0, 'default', 1);

-- 验证插入结果
SELECT 
    id,
    username,
    display_name,
    phone,
    school,
    college,
    role,
    status,
    quota,
    used_quota,
    aff_code
FROM users 
WHERE phone LIKE '138%' 
ORDER BY id;

-- 显示生成的令牌
SELECT 
    t.id,
    t.name,
    t.key,
    u.display_name,
    t.remain_quota,
    t.status
FROM tokens t
JOIN users u ON t.user_id = u.id
WHERE t.key LIKE 'sk-test-token-%'
ORDER BY t.id;
