-- 添加用户表新字段：学校、学院、手机号
ALTER TABLE users ADD COLUMN school VARCHAR(100) DEFAULT '' COMMENT '学校';
ALTER TABLE users ADD COLUMN college VARCHAR(100) DEFAULT '' COMMENT '学院';
ALTER TABLE users ADD COLUMN phone VARCHAR(20) DEFAULT '' COMMENT '手机号';

-- 为手机号字段添加索引以提高查询性能
CREATE INDEX idx_users_phone ON users(phone);
