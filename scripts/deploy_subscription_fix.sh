#!/bin/bash

# 订阅功能修复完整部署脚本
# 使用方法: ./scripts/deploy_subscription_fix.sh

# 配置变量 - 请根据实际情况修改
DB_PATH="./one-api.db"
BACKUP_DIR="./backups"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

echo "=== 订阅功能修复完整部署脚本 ==="

# 1. 检查环境
echo "1. 检查环境..."
if ! command -v sqlite3 &> /dev/null; then
    echo "错误: sqlite3 未安装，请先安装 sqlite3"
    exit 1
fi

# 2. 检查数据库文件
if [ ! -f "$DB_PATH" ]; then
    echo "错误: 数据库文件不存在: $DB_PATH"
    echo "请检查数据库路径是否正确"
    exit 1
fi

# 3. 创建备份目录
echo "2. 创建备份目录..."
mkdir -p $BACKUP_DIR

# 4. 备份数据库
echo "3. 备份数据库..."
BACKUP_FILE="$BACKUP_DIR/one-api_backup_$TIMESTAMP.db"
cp "$DB_PATH" "$BACKUP_FILE"
echo "数据库已备份到: $BACKUP_FILE"

# 5. 执行外键约束移除
echo "4. 执行外键约束移除..."
sqlite3 "$DB_PATH" < scripts/remove_foreign_key_constraint.sql

# 6. 验证表结构
echo "5. 验证表结构..."
echo "subscriptions 表结构:"
sqlite3 "$DB_PATH" ".schema subscriptions"

# 7. 检查新字段是否存在
echo "6. 检查 subscription_articles 表新字段..."
NEW_FIELDS=$(sqlite3 "$DB_PATH" "SELECT name FROM pragma_table_info('subscription_articles') WHERE name IN ('key_points', 'journal_name', 'read_count', 'citation_count', 'rating');")

if [ ! -z "$NEW_FIELDS" ]; then
    echo "新字段已存在:"
    echo "$NEW_FIELDS"
else
    echo "新字段不存在，执行字段添加..."
    sqlite3 "$DB_PATH" < scripts/migrate_subscription_articles.sql
fi

# 8. 验证新字段
echo "7. 验证新字段..."
FINAL_FIELDS=$(sqlite3 "$DB_PATH" "SELECT name FROM pragma_table_info('subscription_articles') WHERE name IN ('key_points', 'journal_name', 'read_count', 'citation_count', 'rating');")
echo "最终字段列表:"
echo "$FINAL_FIELDS"

# 9. 测试数据插入
echo "8. 测试数据插入..."
sqlite3 "$DB_PATH" << EOF
INSERT OR IGNORE INTO subscriptions (create_user_id, topic_name, topic_description, status) 
VALUES (0, '测试订阅', '测试订阅描述', 1);

INSERT OR IGNORE INTO subscription_articles (subscription_id, title, summary, content, author, key_points, journal_name, read_count, citation_count, rating, status)
VALUES (1, '测试文章', '测试文章概要', '测试文章内容', '测试作者', '测试重点提炼', '测试期刊', 100, 5, 8.5, 1);
EOF

echo "测试数据插入完成"

# 10. 验证测试数据
echo "9. 验证测试数据..."
SUBSCRIPTION_COUNT=$(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM subscriptions WHERE topic_name = '测试订阅';")
ARTICLE_COUNT=$(sqlite3 "$DB_PATH" "SELECT COUNT(*) FROM subscription_articles WHERE title = '测试文章';")

echo "测试订阅数量: $SUBSCRIPTION_COUNT"
echo "测试文章数量: $ARTICLE_COUNT"

# 11. 清理测试数据
echo "10. 清理测试数据..."
sqlite3 "$DB_PATH" << EOF
DELETE FROM subscription_articles WHERE title = '测试文章';
DELETE FROM subscriptions WHERE topic_name = '测试订阅';
EOF

echo "测试数据清理完成"

echo "=== 部署完成 ==="
echo "备份文件: $BACKUP_FILE"
echo ""
echo "下一步操作:"
echo "1. 重新编译应用程序: go build"
echo "2. 重启应用程序"
echo "3. 运行测试脚本: ./test_data/test_subscription_fix.sh"
echo ""
echo "如果出现问题，可以使用以下命令恢复:"
echo "cp $BACKUP_FILE $DB_PATH"
