package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 配置
	dbPath := "./data/one-api.db"
	backupDir := "./backups"

	fmt.Println("=== 数据库迁移工具 ===")

	// 1. 检查数据库文件是否存在
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Fatalf("数据库文件不存在: %s", dbPath)
	}

	// 2. 创建备份目录
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		log.Fatalf("创建备份目录失败: %v", err)
	}

	// 3. 备份数据库
	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(backupDir, fmt.Sprintf("one-api_backup_%s.db", timestamp))

	backupData, err := os.ReadFile(dbPath)
	if err != nil {
		log.Fatalf("读取数据库文件失败: %v", err)
	}

	if err := os.WriteFile(backupPath, backupData, 0644); err != nil {
		log.Fatalf("备份数据库失败: %v", err)
	}
	fmt.Printf("数据库已备份到: %s\n", backupPath)

	// 4. 连接数据库
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	defer db.Close()

	// 5. 检查字段是否已存在
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('subscription_articles') WHERE name IN ('key_points', 'journal_name', 'read_count', 'citation_count', 'rating')").Scan(&count)
	if err != nil {
		log.Fatalf("检查字段失败: %v", err)
	}

	if count == 5 {
		fmt.Println("所有新字段已存在，无需迁移")
		return
	}

	// 6. 执行迁移
	fmt.Println("开始执行迁移...")

	// 分割SQL语句并执行
	statements := []string{
		"ALTER TABLE subscription_articles ADD COLUMN key_points TEXT;",
		"ALTER TABLE subscription_articles ADD COLUMN journal_name VARCHAR(200);",
		"ALTER TABLE subscription_articles ADD COLUMN read_count INTEGER DEFAULT 0;",
		"ALTER TABLE subscription_articles ADD COLUMN citation_count INTEGER DEFAULT 0;",
		"ALTER TABLE subscription_articles ADD COLUMN rating DECIMAL(3,1) DEFAULT 0.0;",
	}

	for i, stmt := range statements {
		fmt.Printf("执行语句 %d/5: %s\n", i+1, stmt)
		_, err := db.Exec(stmt)
		if err != nil {
			// 如果字段已存在，忽略错误
			if err.Error() != "duplicate column name: key_points" &&
				err.Error() != "duplicate column name: journal_name" &&
				err.Error() != "duplicate column name: read_count" &&
				err.Error() != "duplicate column name: citation_count" &&
				err.Error() != "duplicate column name: rating" {
				log.Printf("执行语句失败: %v", err)
			} else {
				fmt.Printf("字段已存在，跳过: %s\n", stmt)
			}
		}
	}

	// 7. 验证迁移结果
	fmt.Println("验证迁移结果...")
	rows, err := db.Query("SELECT name FROM pragma_table_info('subscription_articles') WHERE name IN ('key_points', 'journal_name', 'read_count', 'citation_count', 'rating')")
	if err != nil {
		log.Fatalf("验证迁移失败: %v", err)
	}
	defer rows.Close()

	var fields []string
	for rows.Next() {
		var field string
		if err := rows.Scan(&field); err != nil {
			log.Fatalf("读取字段名失败: %v", err)
		}
		fields = append(fields, field)
	}

	fmt.Printf("成功添加的字段: %v\n", fields)
	fmt.Println("=== 迁移完成 ===")
	fmt.Printf("备份文件: %s\n", backupPath)
	fmt.Println("如果迁移出现问题，可以使用以下命令恢复:")
	fmt.Printf("cp %s %s\n", backupPath, dbPath)
}
