-- 更新用户表手机号字段约束
-- 将手机号字段设置为非空，并添加唯一约束

-- 先删除可能存在的重复手机号（如果有的话）
-- 注意：这个操作会删除重复的手机号记录，请在生产环境中谨慎使用
DELETE u1 FROM users u1
INNER JOIN users u2 
WHERE u1.id > u2.id 
AND u1.phone = u2.phone 
AND u1.phone != '';

-- 更新手机号字段约束
ALTER TABLE users MODIFY COLUMN phone VARCHAR(20) NOT NULL DEFAULT '';
ALTER TABLE users ADD UNIQUE INDEX idx_users_phone_unique (phone);

-- 删除之前的普通索引（如果存在）
DROP INDEX IF EXISTS idx_users_phone ON users;

-- 创建话题表
CREATE TABLE IF NOT EXISTS `topics` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `user_id` INTEGER NOT NULL,
  `topic_name` TEXT NOT NULL,
  `model` TEXT DEFAULT 'gpt-3.5-turbo',
  `channel_id` INTEGER DEFAULT 1,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `status` INTEGER DEFAULT 1
);

-- 创建消息表
CREATE TABLE IF NOT EXISTS `messages` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `topic_id` INTEGER NOT NULL,
  `role` TEXT NOT NULL,
  `content` TEXT NOT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `status` INTEGER DEFAULT 1,
  FOREIGN KEY (`topic_id`) REFERENCES `topics` (`id`) ON DELETE CASCADE
);

-- 创建索引
CREATE INDEX IF NOT EXISTS `idx_topics_user_id` ON `topics` (`user_id`);
CREATE INDEX IF NOT EXISTS `idx_topics_status` ON `topics` (`status`);
CREATE INDEX IF NOT EXISTS `idx_topics_created_at` ON `topics` (`created_at`);
CREATE INDEX IF NOT EXISTS `idx_messages_topic_id` ON `messages` (`topic_id`);
CREATE INDEX IF NOT EXISTS `idx_messages_status` ON `messages` (`status`);
CREATE INDEX IF NOT EXISTS `idx_messages_created_at` ON `messages` (`created_at`);
