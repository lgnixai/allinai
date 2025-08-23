-- 创建用户订阅表
CREATE TABLE IF NOT EXISTS subscriptions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    topic_name VARCHAR(100) NOT NULL,
    topic_description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    status INTEGER DEFAULT 1, -- 1: 活跃, 0: 取消
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 创建订阅文章表
CREATE TABLE IF NOT EXISTS subscription_articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    subscription_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    author VARCHAR(100),
    published_at DATETIME,
    article_url VARCHAR(500),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    status INTEGER DEFAULT 1, -- 1: 正常, 0: 删除
    FOREIGN KEY (subscription_id) REFERENCES subscriptions(id) ON DELETE CASCADE
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_topic_name ON subscriptions(topic_name);
CREATE INDEX IF NOT EXISTS idx_subscriptions_status ON subscriptions(status);
CREATE INDEX IF NOT EXISTS idx_subscription_articles_subscription_id ON subscription_articles(subscription_id);
CREATE INDEX IF NOT EXISTS idx_subscription_articles_published_at ON subscription_articles(published_at);
CREATE INDEX IF NOT EXISTS idx_subscription_articles_status ON subscription_articles(status);

-- 添加唯一约束，确保用户不能重复订阅同一主题
CREATE UNIQUE INDEX IF NOT EXISTS idx_subscriptions_user_topic ON subscriptions(user_id, topic_name);
