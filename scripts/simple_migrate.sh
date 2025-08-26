#!/bin/bash

# 简单的数据库迁移脚本
# 使用方法: ./scripts/simple_migrate.sh

# 配置变量 - 请根据实际情况修改
DB_PATH="./one-api.db"
BACKUP_DIR="./backups"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

echo "=== 简单数据库迁移脚本 ==="

# 检查sqlite3是否安装
if ! command -v sqlite3 &> /dev/null; then
    echo "错误: sqlite3 未安装，请先安装 sqlite3"
    echo "Ubuntu/Debian: sudo apt-get install sqlite3"
    echo "CentOS/RHEL: sudo yum install sqlite"
    echo "macOS: brew install sqlite3"
    exit 1
fi

# 1. 检查数据库文件是否存在
if [ ! -f "$DB_PATH" ]; then
    echo "错误: 数据库文件不存在: $DB_PATH"
    echo "请检查数据库路径是否正确"
    exit 1
fi

# 2. 创建备份目录
mkdir -p $BACKUP_DIR

# 3. 备份数据库
BACKUP_FILE="$BACKUP_DIR/one-api_backup_$TIMESTAMP.db"
cp "$DB_PATH" "$BACKUP_FILE"
echo "数据库已备份到: $BACKUP_FILE"

# 4. 检查字段是否已存在
echo "检查现有字段..."
EXISTING_FIELDS=$(sqlite3 "$DB_PATH" "SELECT name FROM pragma_table_info('subscription_articles') WHERE name IN ('key_points', 'journal_name', 'read_count', 'citation_count', 'rating');")

if [ ! -z "$EXISTING_FIELDS" ]; then
    echo "发现已存在的字段:"
    echo "$EXISTING_FIELDS"
    echo "这些字段将被跳过..."
fi

# 5. 执行迁移
echo "开始执行迁移..."

# 添加重点提炼字段
echo "添加 key_points 字段..."
sqlite3 "$DB_PATH" "ALTER TABLE subscription_articles ADD COLUMN key_points TEXT;" 2>/dev/null || echo "key_points 字段已存在，跳过"

# 添加期刊名称字段
echo "添加 journal_name 字段..."
sqlite3 "$DB_PATH" "ALTER TABLE subscription_articles ADD COLUMN journal_name VARCHAR(200);" 2>/dev/null || echo "journal_name 字段已存在，跳过"

# 添加阅读次数字段
echo "添加 read_count 字段..."
sqlite3 "$DB_PATH" "ALTER TABLE subscription_articles ADD COLUMN read_count INTEGER DEFAULT 0;" 2>/dev/null || echo "read_count 字段已存在，跳过"

# 添加引用次数字段
echo "添加 citation_count 字段..."
sqlite3 "$DB_PATH" "ALTER TABLE subscription_articles ADD COLUMN citation_count INTEGER DEFAULT 0;" 2>/dev/null || echo "citation_count 字段已存在，跳过"

# 添加评分的字段
echo "添加 rating 字段..."
sqlite3 "$DB_PATH" "ALTER TABLE subscription_articles ADD COLUMN rating DECIMAL(3,1) DEFAULT 0.0;" 2>/dev/null || echo "rating 字段已存在，跳过"

# 6. 验证迁移结果
echo "验证迁移结果..."
NEW_FIELDS=$(sqlite3 "$DB_PATH" "SELECT name FROM pragma_table_info('subscription_articles') WHERE name IN ('key_points', 'journal_name', 'read_count', 'citation_count', 'rating');")

echo "成功添加的字段:"
echo "$NEW_FIELDS"

# 7. 显示表结构
echo "当前 subscription_articles 表结构:"
sqlite3 "$DB_PATH" ".schema subscription_articles"

echo "=== 迁移完成 ==="
echo "备份文件: $BACKUP_FILE"
echo "如果迁移出现问题，可以使用以下命令恢复:"
echo "cp $BACKUP_FILE $DB_PATH"
