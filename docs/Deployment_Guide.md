# 部署指南

## 概述

本文档介绍如何在新环境中正确部署 One-API 应用，确保所有功能模块正常工作。

## 环境要求

- Go 1.19+
- SQLite 或 MySQL 或 PostgreSQL
- Redis (可选，用于缓存)

## 部署步骤

### 1. 克隆代码

```bash
git clone <repository-url>
cd new-api
```

### 2. 安装依赖

```bash
# 后端依赖
go mod download

# 前端依赖
cd web
npm install
cd ..
```

### 3. 构建应用

```bash
# 构建后端
go build -o main .

# 构建前端
cd web
npm run build
cd ..
```

### 4. 配置环境变量

创建 `.env` 文件：

```env
# 数据库配置
SQL_DSN=local  # 使用 SQLite
# 或者使用 MySQL: SQL_DSN=user:password@tcp(localhost:3306)/database_name

# Redis 配置 (可选)
REDIS_CONN_STRING=redis://localhost:6379

# 其他配置
PORT=4000
SESSION_SECRET=your_random_secret_string
```

### 5. 启动应用

```bash
./main
```

## 数据库迁移

应用启动时会自动执行数据库迁移，创建以下表：

### 核心表
- `users` - 用户表
- `channels` - 渠道表
- `tokens` - 令牌表
- `options` - 配置表
- `logs` - 日志表

### 新增功能表
- `subscriptions` - 订阅表
- `subscription_articles` - 订阅文章表
- `topics` - 话题表
- `messages` - 消息表

## 验证部署

### 1. 检查应用状态

访问 `http://localhost:4000/api/status` 确认应用正常运行。

### 2. 检查数据库表

使用数据库客户端检查以下表是否存在：

```sql
-- SQLite
.tables

-- MySQL
SHOW TABLES;

-- PostgreSQL
\dt
```

确保以下表存在：
- `subscriptions`
- `subscription_articles`
- `topics`
- `messages`

### 3. 测试功能

1. **用户管理**：注册/登录用户
2. **订阅功能**：创建和管理订阅
3. **话题功能**：创建话题和发送消息

## 常见问题

### 问题：订阅和话题相关的表不存在

**原因**：数据库迁移没有正确执行。

**解决方案**：
1. 确保应用版本包含最新的数据库迁移代码
2. 重启应用，让迁移重新执行
3. 检查日志中是否有迁移错误信息

### 问题：前端页面显示错误

**原因**：前端构建文件缺失或配置错误。

**解决方案**：
1. 确保 `web/dist` 目录存在且包含构建文件
2. 重新构建前端：`cd web && npm run build`
3. 检查 `main.go` 中的 `embed` 指令是否正确

### 问题：API 接口返回 404

**原因**：路由配置错误或中间件问题。

**解决方案**：
1. 检查 `router/api-router.go` 中的路由配置
2. 确认中间件配置正确
3. 查看应用日志中的错误信息

## 生产环境部署

### Docker 部署

使用提供的 `docker-compose.yml`：

```bash
docker-compose up -d
```

### 手动部署

1. 构建生产版本
2. 配置生产环境变量
3. 使用进程管理器（如 systemd）管理应用
4. 配置反向代理（如 Nginx）

## 监控和维护

### 日志监控

- 应用日志：`./logs/`
- 错误日志：检查 `ERROR_LOG_ENABLED` 配置

### 数据库维护

- 定期备份数据库
- 监控数据库性能
- 清理过期日志数据

### 性能优化

- 启用 Redis 缓存
- 配置数据库连接池
- 优化查询性能

## 技术支持

如遇到部署问题，请：

1. 检查应用日志
2. 确认环境配置
3. 参考本文档的常见问题部分
4. 提交 Issue 到项目仓库
