# 线上数据库迁移部署指南

## 概述

本指南将帮助您在线上环境中安全地执行数据库迁移，为 `subscription_articles` 表添加新字段。

## 准备工作

### 1. 确认环境
- 确认服务器上已安装 `sqlite3` 命令行工具
- 确认有足够的磁盘空间进行备份
- 确认有数据库文件的读写权限

### 2. 检查依赖
```bash
# 检查 sqlite3 是否安装
sqlite3 --version

# 如果没有安装，请根据系统安装：
# Ubuntu/Debian:
sudo apt-get update && sudo apt-get install sqlite3

# CentOS/RHEL:
sudo yum install sqlite

# macOS:
brew install sqlite3
```

## 迁移方法

### 方法一：使用自动化脚本（推荐）

#### 1. 上传脚本到服务器
```bash
# 将以下文件上传到服务器：
# - scripts/simple_migrate.sh
# - scripts/migrate_subscription_articles.sql
```

#### 2. 修改配置
编辑 `scripts/simple_migrate.sh`，修改数据库路径：
```bash
# 找到这一行并修改为实际的数据库路径
DB_PATH="./data/one-api.db"  # 修改为实际路径
```

#### 3. 执行迁移
```bash
# 给脚本执行权限
chmod +x scripts/simple_migrate.sh

# 执行迁移
./scripts/simple_migrate.sh
```

### 方法二：手动执行SQL

#### 1. 备份数据库
```bash
# 创建备份目录
mkdir -p backups

# 备份数据库（替换为实际路径）
cp /path/to/your/database.db backups/database_backup_$(date +%Y%m%d_%H%M%S).db
```

#### 2. 执行迁移SQL
```bash
# 方法A：直接执行SQL文件
sqlite3 /path/to/your/database.db < scripts/migrate_subscription_articles.sql

# 方法B：逐条执行SQL语句
sqlite3 /path/to/your/database.db << EOF
ALTER TABLE subscription_articles ADD COLUMN key_points TEXT;
ALTER TABLE subscription_articles ADD COLUMN journal_name VARCHAR(200);
ALTER TABLE subscription_articles ADD COLUMN read_count INTEGER DEFAULT 0;
ALTER TABLE subscription_articles ADD COLUMN citation_count INTEGER DEFAULT 0;
ALTER TABLE subscription_articles ADD COLUMN rating DECIMAL(3,1) DEFAULT 0.0;
EOF
```

#### 3. 验证迁移
```bash
# 检查表结构
sqlite3 /path/to/your/database.db ".schema subscription_articles"

# 检查新字段是否存在
sqlite3 /path/to/your/database.db "SELECT name FROM pragma_table_info('subscription_articles') WHERE name IN ('key_points', 'journal_name', 'read_count', 'citation_count', 'rating');"
```

## 部署步骤

### 1. 停止应用程序
```bash
# 停止当前运行的应用
sudo systemctl stop your-app-service
# 或者
pkill -f your-app-process
```

### 2. 执行数据库迁移
```bash
# 使用自动化脚本
./scripts/simple_migrate.sh

# 或者手动执行
# ... (参考上面的手动执行SQL部分)
```

### 3. 更新应用程序
```bash
# 上传新的应用程序文件
# 重新编译（如果需要）
go build -o your-app

# 或者直接替换二进制文件
```

### 4. 重启应用程序
```bash
# 启动应用
sudo systemctl start your-app-service
# 或者
./your-app

# 检查应用状态
sudo systemctl status your-app-service
```

### 5. 验证部署
```bash
# 检查应用日志
sudo journalctl -u your-app-service -f

# 测试API接口
curl -X GET "http://your-server:port/api/subscriptions/articles" \
  -H "Authorization: Bearer your-token"
```

## 常见问题处理

### 1. 字段已存在错误
如果遇到 "duplicate column name" 错误，说明字段已经存在，可以忽略此错误。

### 2. 权限问题
```bash
# 确保有数据库文件的读写权限
sudo chown your-user:your-group /path/to/database.db
sudo chmod 644 /path/to/database.db
```

### 3. 磁盘空间不足
```bash
# 检查磁盘空间
df -h

# 清理不必要的文件
sudo apt-get clean  # Ubuntu/Debian
sudo yum clean all  # CentOS/RHEL
```

### 4. 应用启动失败
```bash
# 检查应用日志
tail -f /var/log/your-app.log

# 检查数据库连接
sqlite3 /path/to/database.db "SELECT COUNT(*) FROM subscription_articles;"
```

## 回滚方案

如果迁移出现问题，可以快速回滚：

### 1. 停止应用
```bash
sudo systemctl stop your-app-service
```

### 2. 恢复数据库
```bash
# 使用备份文件恢复
cp backups/database_backup_TIMESTAMP.db /path/to/your/database.db
```

### 3. 重启应用
```bash
sudo systemctl start your-app-service
```

## 监控和验证

### 1. 监控指标
- 应用启动状态
- API响应时间
- 数据库连接状态
- 错误日志

### 2. 功能验证
```bash
# 测试创建文章（需要管理员权限）
curl -X POST "http://your-server:port/api/subscription_articles" \
  -H "Authorization: Bearer admin-token" \
  -H "Content-Type: application/json" \
  -d '{
    "subscription_id": 1,
    "title": "测试文章",
    "key_points": "测试重点提炼",
    "journal_name": "测试期刊",
    "read_count": 100,
    "citation_count": 5,
    "rating": 8.5
  }'

# 测试获取文章列表
curl -X GET "http://your-server:port/api/subscriptions/articles" \
  -H "Authorization: Bearer your-token"
```

## 安全注意事项

1. **备份重要**: 始终在迁移前备份数据库
2. **权限控制**: 确保只有授权用户能访问数据库文件
3. **日志记录**: 记录所有迁移操作
4. **测试环境**: 建议先在测试环境验证迁移脚本
5. **监控**: 迁移后密切监控应用状态

## 联系支持

如果在迁移过程中遇到问题，请：
1. 保存错误日志
2. 记录执行步骤
3. 联系技术支持团队
