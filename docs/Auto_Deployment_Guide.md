# 自动部署指南

## 概述

本文档介绍如何配置自动部署系统，实现Git提交后自动部署到服务器的功能。

## 方案选择

### 方案一：Git Hooks + 远程脚本（推荐）
- **优点**：简单直接，无需额外服务
- **适用场景**：单服务器部署，直接Git推送

### 方案二：Webhook + Flask服务器
- **优点**：支持GitHub/GitLab webhook，更灵活
- **适用场景**：使用GitHub/GitLab托管代码

### 方案三：Docker Compose + Webhook
- **优点**：容器化部署，易于管理
- **适用场景**：生产环境，需要容器化部署

## 方案一：Git Hooks 配置

### 1. 服务器端配置

#### 1.1 设置Git仓库
```bash
# 在服务器上创建裸仓库
mkdir -p /var/git/new-api.git
cd /var/git/new-api.git
git init --bare

# 设置post-receive hook
cp scripts/post-receive hooks/
chmod +x hooks/post-receive
```

#### 1.2 配置部署脚本
```bash
# 修改scripts/deploy.sh中的路径
PROJECT_DIR="/var/www/new-api"  # 修改为实际项目路径

# 设置执行权限
chmod +x scripts/deploy.sh
```

#### 1.3 配置本地Git远程
```bash
# 在本地项目中添加服务器远程仓库
git remote add production user@server:/var/git/new-api.git

# 推送代码
git push production main
```

### 2. 使用systemd管理服务

#### 2.1 安装服务
```bash
# 复制服务文件
sudo cp scripts/new-api.service /etc/systemd/system/

# 修改服务文件中的路径
sudo nano /etc/systemd/system/new-api.service

# 重新加载systemd
sudo systemctl daemon-reload

# 启用服务
sudo systemctl enable new-api

# 启动服务
sudo systemctl start new-api
```

#### 2.2 服务管理命令
```bash
# 查看服务状态
sudo systemctl status new-api

# 重启服务
sudo systemctl restart new-api

# 查看日志
sudo journalctl -u new-api -f
```

## 方案二：Webhook 配置

### 1. 启动Webhook服务器

#### 1.1 安装依赖
```bash
# 安装Python和Flask
sudo apt update
sudo apt install python3 python3-pip
pip3 install flask
```

#### 1.2 配置环境变量
```bash
# 设置webhook密钥
export WEBHOOK_SECRET="your_secret_key"
export PROJECT_DIR="/var/www/new-api"
```

#### 1.3 启动服务
```bash
# 直接运行
python3 scripts/webhook-server.py

# 或使用systemd
sudo cp scripts/webhook-server.service /etc/systemd/system/
sudo systemctl enable webhook-server
sudo systemctl start webhook-server
```

### 2. 配置GitHub/GitLab Webhook

#### 2.1 GitHub配置
1. 进入仓库设置 → Webhooks
2. 添加新的webhook
3. 设置URL：`http://your-server:8080/webhook`
4. 设置Secret：与服务器上的`WEBHOOK_SECRET`一致
5. 选择事件：仅选择"Push events"

#### 2.2 GitLab配置
1. 进入仓库设置 → Webhooks
2. 设置URL：`http://your-server:8080/webhook`
3. 设置Secret Token：与服务器上的`WEBHOOK_SECRET`一致
4. 选择触发事件：仅选择"Push events"

## 方案三：Docker Compose 配置

### 1. 环境准备
```bash
# 创建环境变量文件
cat > .env << EOF
SESSION_SECRET=your_session_secret
WEBHOOK_SECRET=your_webhook_secret
EOF
```

### 2. 启动服务
```bash
# 启动所有服务
docker-compose -f docker-compose.prod.yml up -d

# 查看服务状态
docker-compose -f docker-compose.prod.yml ps

# 查看日志
docker-compose -f docker-compose.prod.yml logs -f
```

### 3. 配置Webhook
按照方案二的步骤配置GitHub/GitLab webhook，URL为：
`http://your-server:8080/webhook`

## 监控和日志

### 1. 查看部署日志
```bash
# 部署脚本日志
tail -f /var/log/deploy.log

# Git hook日志
tail -f /var/log/git-hook.log

# Webhook日志
tail -f /var/log/webhook.log

# 应用日志
tail -f /var/log/new-api.log
```

### 2. 健康检查
```bash
# 检查应用状态
curl http://localhost:4000/api/status

# 检查webhook服务器状态
curl http://localhost:8080/health
```

## 故障排除

### 1. 常见问题

#### 问题：部署脚本权限不足
```bash
# 解决方案
chmod +x scripts/deploy.sh
chown www-data:www-data scripts/deploy.sh
```

#### 问题：Git hook不触发
```bash
# 检查hook文件权限
chmod +x hooks/post-receive

# 检查Git仓库配置
git config --list
```

#### 问题：Webhook签名验证失败
```bash
# 检查密钥配置
echo $WEBHOOK_SECRET

# 检查GitHub/GitLab webhook配置
```

#### 问题：服务启动失败
```bash
# 检查端口占用
netstat -tlnp | grep :4000

# 检查日志
journalctl -u new-api -f
```

### 2. 调试模式

#### 启用详细日志
```bash
# 修改deploy.sh，添加调试信息
set -x  # 在脚本开头添加
```

#### 手动测试部署
```bash
# 手动执行部署脚本
./scripts/deploy.sh main
```

## 安全考虑

### 1. 网络安全
- 使用HTTPS进行webhook通信
- 配置防火墙，只开放必要端口
- 使用强密钥进行webhook签名验证

### 2. 权限控制
- 使用专用用户运行服务
- 限制文件系统访问权限
- 定期更新密钥和密码

### 3. 监控告警
- 配置部署失败告警
- 监控服务健康状态
- 设置日志轮转和清理

## 最佳实践

### 1. 部署策略
- 使用蓝绿部署或滚动更新
- 配置回滚机制
- 在部署前进行健康检查

### 2. 备份策略
- 定期备份数据库
- 备份配置文件
- 测试恢复流程

### 3. 性能优化
- 使用CDN加速静态资源
- 配置数据库连接池
- 启用缓存机制

## 总结

选择合适的自动部署方案取决于您的具体需求：

- **简单部署**：使用Git Hooks
- **团队协作**：使用Webhook + GitHub/GitLab
- **生产环境**：使用Docker Compose

无论选择哪种方案，都要确保：
1. 配置正确的权限和安全设置
2. 建立完善的监控和日志系统
3. 制定故障恢复计划
4. 定期测试部署流程
