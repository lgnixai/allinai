package model

import (
	"one-api/common"
)

// SetupSQLiteHooks 设置SQLite钩子，自动启用外键约束
func SetupSQLiteHooks() {
	if !common.UsingSQLite {
		return
	}

	// 在数据库连接建立后启用外键约束
	sqlDB, err := DB.DB()
	if err != nil {
		common.SysLog("failed to get sql.DB: " + err.Error())
		return
	}

	// 启用外键约束
	_, err = sqlDB.Exec("PRAGMA foreign_keys=ON;")
	if err != nil {
		common.SysLog("failed to enable foreign keys: " + err.Error())
		return
	}

	common.SysLog("SQLite foreign keys enabled successfully")
}

// CreateTablesWithForeignKeys 创建带外键约束的表（仅在新部署时使用）
func CreateTablesWithForeignKeys() error {
	if !common.UsingSQLite {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// 启用外键约束
	_, err = sqlDB.Exec("PRAGMA foreign_keys=ON;")
	if err != nil {
		return err
	}

	// 检查表是否已存在
	var count int
	err = sqlDB.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='topics'").Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		// 表已存在，不重新创建
		common.SysLog("tables already exist, skipping foreign key table creation")
		return nil
	}

	// 创建表的SQL语句（带外键约束）
	createTableSQLs := []string{
		`CREATE TABLE topics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			topic_name TEXT NOT NULL,
			model TEXT DEFAULT 'gpt-3.5-turbo',
			channel_id INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			status INTEGER DEFAULT 1,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			topic_id INTEGER NOT NULL,
			role TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			status INTEGER DEFAULT 1,
			FOREIGN KEY (topic_id) REFERENCES topics(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE subscriptions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			create_user_id INTEGER NOT NULL,
			topic_name VARCHAR(100) NOT NULL,
			topic_description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			status INTEGER DEFAULT 1,
			FOREIGN KEY (create_user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE user_subscriptions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			subscription_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			status INTEGER DEFAULT 1,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (subscription_id) REFERENCES subscriptions(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE subscription_articles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			subscription_id INTEGER NOT NULL,
			title VARCHAR(255) NOT NULL,
			summary TEXT,
			content TEXT,
			author VARCHAR(100),
			published_at DATETIME,
			article_url VARCHAR(500),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			status INTEGER DEFAULT 1,
			FOREIGN KEY (subscription_id) REFERENCES subscriptions(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE system_recommendations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			category VARCHAR(100),
			subscription_count INTEGER DEFAULT 0,
			article_count INTEGER DEFAULT 0,
			status INTEGER DEFAULT 1,
			sort_order INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	// 创建索引
	createIndexSQLs := []string{
		`CREATE INDEX idx_topics_user_id ON topics(user_id)`,
		`CREATE INDEX idx_topics_status ON topics(status)`,
		`CREATE INDEX idx_topics_created_at ON topics(created_at)`,
		`CREATE INDEX idx_messages_topic_id ON messages(topic_id)`,
		`CREATE INDEX idx_messages_status ON messages(status)`,
		`CREATE INDEX idx_messages_created_at ON messages(created_at)`,
		`CREATE INDEX idx_subscriptions_create_user_id ON subscriptions(create_user_id)`,
		`CREATE INDEX idx_subscriptions_topic_name ON subscriptions(topic_name)`,
		`CREATE INDEX idx_subscriptions_status ON subscriptions(status)`,
		`CREATE INDEX idx_subscription_articles_subscription_id ON subscription_articles(subscription_id)`,
		`CREATE INDEX idx_subscription_articles_published_at ON subscription_articles(published_at)`,
		`CREATE INDEX idx_subscription_articles_status ON subscription_articles(status)`,
		`CREATE UNIQUE INDEX idx_subscriptions_create_user_topic ON subscriptions(create_user_id, topic_name)`,
		`CREATE INDEX idx_user_subscriptions_user_id ON user_subscriptions(user_id)`,
		`CREATE INDEX idx_user_subscriptions_subscription_id ON user_subscriptions(subscription_id)`,
		`CREATE UNIQUE INDEX idx_user_subscriptions_user_subscription ON user_subscriptions(user_id, subscription_id)`,
		`CREATE INDEX idx_system_recommendations_category ON system_recommendations(category)`,
		`CREATE INDEX idx_system_recommendations_status ON system_recommendations(status)`,
		`CREATE INDEX idx_system_recommendations_sort_order ON system_recommendations(sort_order)`,
	}

	// 执行创建表的SQL
	for _, sql := range createTableSQLs {
		_, err = sqlDB.Exec(sql)
		if err != nil {
			common.SysLog("failed to create table: " + err.Error())
			return err
		}
	}

	// 执行创建索引的SQL
	for _, sql := range createIndexSQLs {
		_, err = sqlDB.Exec(sql)
		if err != nil {
			common.SysLog("failed to create index: " + err.Error())
			// 索引创建失败不影响主要功能，只记录日志
		}
	}

	common.SysLog("SQLite tables with foreign keys created successfully")
	return nil
}

// CheckAndFixForeignKeys 检查并修复外键约束
func CheckAndFixForeignKeys() error {
	if !common.UsingSQLite {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// 检查外键约束是否存在
	checkQueries := []struct {
		tableName string
		query     string
	}{
		{"topics", `SELECT COUNT(*) FROM pragma_foreign_key_list('topics') WHERE "table"='users' AND "from"='user_id'`},
		{"messages", `SELECT COUNT(*) FROM pragma_foreign_key_list('messages') WHERE "table"='topics' AND "from"='topic_id'`},
		{"subscriptions", `SELECT COUNT(*) FROM pragma_foreign_key_list('subscriptions') WHERE "table"='users' AND "from"='create_user_id'`},
		{"subscription_articles", `SELECT COUNT(*) FROM pragma_foreign_key_list('subscription_articles') WHERE "table"='subscriptions' AND "from"='subscription_id'`},
	}

	missingFKs := false
	for _, check := range checkQueries {
		var count int
		err := sqlDB.QueryRow(check.query).Scan(&count)
		if err != nil {
			// 表可能不存在，跳过
			continue
		}
		if count == 0 {
			common.SysLog("missing foreign key constraint detected for table: " + check.tableName)
			missingFKs = true
		}
	}

	if missingFKs {
		common.SysLog("foreign key constraints are missing, but tables already exist")
		common.SysLog("for new deployments, foreign keys will be automatically created")
		common.SysLog("for existing deployments, manual intervention may be required")
	}

	return nil
}
