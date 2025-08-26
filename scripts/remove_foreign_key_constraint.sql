-- 去掉 subscriptions 表的外键约束
-- 执行时间: 2024年

-- 注意：SQLite 不支持直接删除外键约束，需要重建表

-- 1. 创建临时表
CREATE TABLE subscriptions_temp (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    create_user_id INTEGER NOT NULL,
    topic_name VARCHAR(100) NOT NULL,
    topic_description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    status INTEGER DEFAULT 1
);

-- 2. 复制数据
INSERT INTO subscriptions_temp 
SELECT id, create_user_id, topic_name, topic_description, created_at, updated_at, status 
FROM subscriptions;

-- 3. 删除原表
DROP TABLE subscriptions;

-- 4. 重命名临时表
ALTER TABLE subscriptions_temp RENAME TO subscriptions;

-- 5. 重新创建索引
CREATE INDEX idx_subscriptions_create_user_id ON subscriptions(create_user_id);
CREATE INDEX idx_subscriptions_topic_name ON subscriptions(topic_name);
CREATE INDEX idx_subscriptions_status ON subscriptions(status);
CREATE UNIQUE INDEX idx_subscriptions_create_user_topic ON subscriptions(create_user_id, topic_name);

-- 6. 验证表结构
.schema subscriptions
