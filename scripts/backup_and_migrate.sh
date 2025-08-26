#!/bin/bash

# 数据库备份和迁移脚本
# 使用方法: ./scripts/backup_and_migrate.sh

# 配置变量
DB_PATH="./data/one-api.db"  # 根据实际数据库路径调整
BACKUP_DIR="./backups"
MIGRATION_SCRIPT="./scripts/migrate_subscription_articles.sql"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

echo "=== 数据库备份和迁移脚本 ==="

# 1. 创建备份目录
echo "1. 创建备份目录..."
mkdir -p $BACKUP_DIR

# 2. 备份数据库
echo "2. 备份数据库..."
BACKUP_FILE="$BACKUP_DIR/one-api_backup_$TIMESTAMP.db"
cp "$DB_PATH" "$BACKUP_FILE"
echo "数据库已备份到: $BACKUP_FILE"

# 3. 检查数据库文件是否存在
if [ ! -f "$DB_PATH" ]; then
    echo "错误: 数据库文件不存在: $DB_PATH"
    echo "请检查数据库路径是否正确"
    exit 1
fi

# 4. 检查迁移脚本是否存在
if [ ! -f "$MIGRATION_SCRIPT" ]; then
    echo "错误: 迁移脚本不存在: $MIGRATION_SCRIPT"
    exit 1
fi

# 5. 执行迁移
echo "3. 执行数据库迁移..."
sqlite3 "$DB_PATH" < "$MIGRATION_SCRIPT"

# 6. 检查迁移结果
echo "4. 检查迁移结果..."
sqlite3 "$DB_PATH" ".schema subscription_articles" | grep -E "(key_points|journal_name|read_count|citation_count|rating)"

echo "=== 迁移完成 ==="
echo "备份文件: $BACKUP_FILE"
echo "如果迁移出现问题，可以使用以下命令恢复:"
echo "cp $BACKUP_FILE $DB_PATH"
