# 新系统部署指南（自动外键约束版本）

## 概述

本指南适用于在新操作系统上部署带有自动外键约束功能的API服务。系统现在会自动处理数据库表创建和外键约束，无需手动干预。

## 1. 环境准备

### 安装Go语言环境
```bash
# 下载并安装Go (推荐1.21+版本)
wget https://golang.org/dl/go1.21.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 验证安装
go version
```

### 安装SQLite (如果使用SQLite)
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install sqlite3

# CentOS/RHEL
sudo yum install sqlite

# macOS
brew install sqlite3
```

### 安装Redis (可选，用于缓存)
```bash
# Ubuntu/Debian
sudo apt install redis-server

# CentOS/RHEL
sudo yum install redis

# macOS
brew install redis
```

## 2. 项目部署

### 克隆项目
```bash
git clone <your-repository-url>
cd new-api
```

### 安装依赖
```bash
go mod download
```

### 配置环境变量
创建 `.env` 文件：
```bash
# 数据库配置
SQL_DSN=one-api.db

# 服务器配置
PORT=9999
SYNC_FREQUENCY=60

# Redis配置 (可选)
REDIS_CONN_STRING=localhost:6379
REDIS_DB=0

# 其他配置
DEBUG=true
LOG_SQL_DSN=
```

## 3. 自动数据库初始化

### 首次运行
```bash
# 启动应用程序，会自动创建数据库和表
PORT=9999 go run .
```

**系统会自动执行以下操作：**
1. ✅ 创建SQLite数据库文件
2. ✅ 启用外键约束 (`PRAGMA foreign_keys=ON`)
3. ✅ 创建所有基础表（users, channels, tokens等）
4. ✅ 创建带外键约束的特殊表：
   - `topics` (外键: user_id → users.id)
   - `messages` (外键: topic_id → topics.id)
   - `subscriptions` (外键: user_id → users.id)
   - `subscription_articles` (外键: subscription_id → subscriptions.id)
5. ✅ 创建所有必要的索引
6. ✅ 验证外键约束完整性

### 验证部署
```bash
# 检查API状态
curl http://localhost:9999/api/status

# 检查数据库表
sqlite3 one-api.db ".tables"

# 检查外键约束
sqlite3 one-api.db "PRAGMA foreign_keys;"
sqlite3 one-api.db ".schema topics"
```

## 4. 生产环境部署

### 编译应用程序
```bash
# 编译为可执行文件
go build -o one-api .

# 或者交叉编译 (Linux)
GOOS=linux GOARCH=amd64 go build -o one-api .
```

### 使用systemd服务 (Linux)
创建服务文件 `/etc/systemd/system/one-api.service`：
```ini
[Unit]
Description=One API Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/your/new-api
ExecStart=/path/to/your/new-api/one-api
Environment=PORT=9999
Environment=SQL_DSN=one-api.db
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启动服务：
```bash
sudo systemctl daemon-reload
sudo systemctl enable one-api
sudo systemctl start one-api
sudo systemctl status one-api
```

### 使用Docker部署
创建 `Dockerfile`：
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o one-api .

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite
WORKDIR /root/
COPY --from=builder /app/one-api .
EXPOSE 9999
CMD ["./one-api"]
```

构建和运行：
```bash
docker build -t one-api .
docker run -d -p 9999:9999 -v $(pwd)/data:/root/data one-api
```

## 5. 外键约束功能

### 自动级联删除
系统现在支持完整的级联删除功能：

- **删除用户** → 自动删除该用户的所有topics和subscriptions
- **删除topic** → 自动删除该topic下的所有messages
- **删除subscription** → 自动删除该subscription下的所有articles

### 数据完整性保护
- 确保引用的数据存在
- 防止孤立数据
- 自动维护数据一致性

## 6. 验证部署

### 检查服务状态
```bash
# 检查端口是否监听
netstat -tlnp | grep 9999

# 测试API
curl http://localhost:9999/api/status
```

### 检查数据库
```bash
# 检查数据库文件
ls -la one-api.db

# 检查表结构
sqlite3 one-api.db ".tables"
sqlite3 one-api.db ".schema topics"
sqlite3 one-api.db ".schema messages"
```

### 测试外键约束
```bash
# 启用外键约束
sqlite3 one-api.db "PRAGMA foreign_keys=ON;"

# 创建测试数据
sqlite3 one-api.db "INSERT INTO users (username, display_name, role, status, quota) VALUES ('testuser', 'Test User', 1, 1, 1000000);"
sqlite3 one-api.db "INSERT INTO topics (user_id, topic_name, model) VALUES (1, 'Test Topic', 'gpt-3.5-turbo');"
sqlite3 one-api.db "INSERT INTO messages (topic_id, role, content) VALUES (1, 'user', 'Hello');"

# 测试级联删除
sqlite3 one-api.db "DELETE FROM topics WHERE id = 1;"
sqlite3 one-api.db "SELECT COUNT(*) FROM messages WHERE topic_id = 1;"
# 应该返回 0，表示messages被级联删除

# 清理测试数据
sqlite3 one-api.db "DELETE FROM users WHERE username = 'testuser';"
```

## 7. 故障排除

### 常见问题

1. **外键约束未启用**
   ```bash
   # 检查外键状态
   sqlite3 one-api.db "PRAGMA foreign_keys;"
   # 应该返回 1
   ```

2. **表创建失败**
   ```bash
   # 检查日志
   tail -f app.log
   
   # 重新创建数据库
   rm -f one-api.db
   PORT=9999 go run .
   ```

3. **端口被占用**
   ```bash
   # 查找占用进程
   lsof -i :9999
   
   # 杀死进程
   lsof -ti:9999 | xargs kill -9
   ```

## 8. 优势总结

### 自动化特性
- ✅ **自动数据库创建**：无需手动执行SQL脚本
- ✅ **自动外键约束**：自动启用和创建外键约束
- ✅ **自动索引创建**：自动创建性能优化索引
- ✅ **自动数据完整性**：确保数据一致性

### 部署简化
- ✅ **零手动干预**：新系统部署时无需手动处理外键
- ✅ **自动迁移**：数据库结构自动升级
- ✅ **向后兼容**：支持现有数据库的升级

### 功能完整性
- ✅ **完整的外键约束**：所有关系都有外键保护
- ✅ **级联删除**：自动维护数据一致性
- ✅ **性能优化**：自动创建必要的索引

## 9. 监控和维护

### 日志监控
```bash
# 查看应用日志
tail -f app.log

# 查看systemd日志
sudo journalctl -u one-api -f
```

### 数据库维护
```bash
# 数据库完整性检查
sqlite3 one-api.db "PRAGMA integrity_check;"

# 外键约束检查
sqlite3 one-api.db "PRAGMA foreign_key_check;"
```

### 备份策略
```bash
# 创建备份脚本
cat > backup.sh << 'EOF'
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
cp one-api.db "backup/one-api_$DATE.db"
find backup/ -name "*.db" -mtime +7 -delete
EOF

chmod +x backup.sh

# 添加到crontab
echo "0 2 * * * /path/to/backup.sh" | crontab -
```

---

**注意**：本版本完全自动化了外键约束的处理，在新系统部署时无需任何手动干预。系统会自动检测并创建所有必要的数据库结构和约束。
