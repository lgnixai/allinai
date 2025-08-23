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
