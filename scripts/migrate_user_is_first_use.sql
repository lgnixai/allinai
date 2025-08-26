-- 为 users 表添加 is_first_use 字段（是否首次使用）
-- 1 表示首次使用（默认），0 表示非首次使用

ALTER TABLE users ADD COLUMN is_first_use INTEGER DEFAULT 1;
