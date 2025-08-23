-- SQLite 数据库迁移脚本：删除用户表中的密码字段
-- 这个迁移脚本将删除用户表中的 password 和 original_password 字段

-- 由于 SQLite 不支持直接删除列，我们需要重新创建表
-- 首先备份现有数据
CREATE TABLE users_backup AS SELECT * FROM users;

-- 删除原表
DROP TABLE users;

-- 重新创建表，不包含密码字段
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(12) UNIQUE NOT NULL,
    display_name VARCHAR(20),
    role INTEGER DEFAULT 1,
    status INTEGER DEFAULT 1,
    email VARCHAR(50),
    github_id VARCHAR(50),
    oidc_id VARCHAR(50),
    wechat_id VARCHAR(50),
    telegram_id VARCHAR(50),
    access_token CHAR(32) UNIQUE,
    quota INTEGER DEFAULT 0,
    used_quota INTEGER DEFAULT 0,
    request_count INTEGER DEFAULT 0,
    "group" VARCHAR(64) DEFAULT 'default',
    aff_code VARCHAR(32) UNIQUE,
    aff_count INTEGER DEFAULT 0,
    aff_quota INTEGER DEFAULT 0,
    aff_history INTEGER DEFAULT 0,
    inviter_id INTEGER,
    deleted_at DATETIME,
    linux_do_id VARCHAR(50),
    setting TEXT,
    remark VARCHAR(255),
    stripe_customer VARCHAR(64),
    school VARCHAR(100),
    college VARCHAR(100),
    phone VARCHAR(20) NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 恢复数据（排除密码字段）
INSERT INTO users (
    id, username, display_name, role, status, email, github_id, oidc_id, 
    wechat_id, telegram_id, access_token, quota, used_quota, request_count, 
    "group", aff_code, aff_count, aff_quota, aff_history, inviter_id, 
    deleted_at, linux_do_id, setting, remark, stripe_customer, school, 
    college, phone, created_at, updated_at
) SELECT 
    id, username, display_name, role, status, email, github_id, oidc_id, 
    wechat_id, telegram_id, access_token, quota, used_quota, request_count, 
    "group", aff_code, aff_count, aff_quota, aff_history, inviter_id, 
    deleted_at, linux_do_id, setting, remark, stripe_customer, school, 
    college, phone, created_at, updated_at
FROM users_backup;

-- 删除备份表
DROP TABLE users_backup;

-- 添加索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_aff_code ON users(aff_code);
CREATE INDEX IF NOT EXISTS idx_users_access_token ON users(access_token);
