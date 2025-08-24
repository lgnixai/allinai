-- 修复数据库迁移错误
PRAGMA foreign_keys=OFF;

-- 删除所有可能存在的临时表
DROP TABLE IF EXISTS messages__temp;
DROP TABLE IF EXISTS topics__temp;
DROP TABLE IF EXISTS subscriptions__temp;
DROP TABLE IF EXISTS subscription_articles__temp;
DROP TABLE IF EXISTS users__temp;

-- 删除现有表（如果存在）
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS topics;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS subscription_articles;

-- 重新创建topics表
CREATE TABLE topics (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  topic_name TEXT NOT NULL,
  model TEXT DEFAULT 'gpt-3.5-turbo',
  channel_id INTEGER DEFAULT 1,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  status INTEGER DEFAULT 1
);

-- 重新创建messages表
CREATE TABLE messages (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  topic_id INTEGER NOT NULL,
  role TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  status INTEGER DEFAULT 1
);

-- 重新创建subscriptions表
CREATE TABLE subscriptions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    topic_name VARCHAR(100) NOT NULL,
    topic_description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    status INTEGER DEFAULT 1
);

-- 重新创建subscription_articles表
CREATE TABLE subscription_articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    subscription_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    author VARCHAR(100),
    published_at DATETIME,
    article_url VARCHAR(500),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    status INTEGER DEFAULT 1
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_topics_user_id ON topics(user_id);
CREATE INDEX IF NOT EXISTS idx_topics_status ON topics(status);
CREATE INDEX IF NOT EXISTS idx_topics_created_at ON topics(created_at);
CREATE INDEX IF NOT EXISTS idx_messages_topic_id ON messages(topic_id);
CREATE INDEX IF NOT EXISTS idx_messages_status ON messages(status);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);

CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_topic_name ON subscriptions(topic_name);
CREATE INDEX IF NOT EXISTS idx_subscriptions_status ON subscriptions(status);
CREATE INDEX IF NOT EXISTS idx_subscription_articles_subscription_id ON subscription_articles(subscription_id);
CREATE INDEX IF NOT EXISTS idx_subscription_articles_published_at ON subscription_articles(published_at);
CREATE INDEX IF NOT EXISTS idx_subscription_articles_status ON subscription_articles(status);

-- 添加唯一约束
CREATE UNIQUE INDEX IF NOT EXISTS idx_subscriptions_user_topic ON subscriptions(user_id, topic_name);

PRAGMA foreign_keys=ON;
