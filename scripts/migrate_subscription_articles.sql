-- 为 subscription_articles 表添加新字段的迁移脚本
-- 执行时间: 2024年

-- 添加重点提炼字段
ALTER TABLE subscription_articles ADD COLUMN key_points TEXT;

-- 添加期刊名称字段
ALTER TABLE subscription_articles ADD COLUMN journal_name VARCHAR(200);

-- 添加阅读次数字段
ALTER TABLE subscription_articles ADD COLUMN read_count INTEGER DEFAULT 0;

-- 添加引用次数字段
ALTER TABLE subscription_articles ADD COLUMN citation_count INTEGER DEFAULT 0;

-- 添加评分的字段
ALTER TABLE subscription_articles ADD COLUMN rating DECIMAL(3,1) DEFAULT 0.0;

-- 为新增字段创建索引（可选，根据查询需求决定）
-- CREATE INDEX idx_subscription_articles_journal_name ON subscription_articles(journal_name);
-- CREATE INDEX idx_subscription_articles_read_count ON subscription_articles(read_count);
-- CREATE INDEX idx_subscription_articles_citation_count ON subscription_articles(citation_count);
-- CREATE INDEX idx_subscription_articles_rating ON subscription_articles(rating);
