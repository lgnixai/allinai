-- 删除用户表中的密码字段
-- 这个迁移脚本将删除用户表中的 password 和 original_password 字段

-- 删除 password 字段
ALTER TABLE users DROP COLUMN IF EXISTS password;

-- 删除 original_password 字段（如果存在）
ALTER TABLE users DROP COLUMN IF EXISTS original_password;

-- 更新用户模型，确保手机号是唯一且必填的
ALTER TABLE users ALTER COLUMN phone SET NOT NULL;
ALTER TABLE users ADD CONSTRAINT users_phone_unique UNIQUE (phone);

-- 添加索引以提高查询性能
CREATE INDEX IF NOT EXISTS idx_users_phone ON users(phone);
