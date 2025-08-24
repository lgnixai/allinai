#!/bin/bash

# 自动部署脚本
# 使用方法: ./deploy.sh [branch_name]

set -e  # 遇到错误立即退出

# 配置变量
PROJECT_DIR="/path/to/your/project"  # 修改为您的项目路径
BRANCH=${1:-main}  # 默认使用main分支
LOG_FILE="/var/log/deploy.log"
PID_FILE="/tmp/new-api.pid"

# 日志函数
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

# 检查项目目录
if [ ! -d "$PROJECT_DIR" ]; then
    log "错误: 项目目录不存在: $PROJECT_DIR"
    exit 1
fi

cd "$PROJECT_DIR"

log "开始部署分支: $BRANCH"

# 1. 拉取最新代码
log "拉取最新代码..."
git fetch origin
git reset --hard origin/$BRANCH
git clean -fd

# 2. 更新Go依赖
log "更新Go依赖..."
go mod download
go mod tidy

# 3. 构建应用
log "构建应用..."
go build -o main .

# 4. 构建前端
log "构建前端..."
cd web
npm install
npm run build
cd ..

# 5. 停止现有服务
log "停止现有服务..."
if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if kill -0 "$PID" 2>/dev/null; then
        kill "$PID"
        sleep 5
        if kill -0 "$PID" 2>/dev/null; then
            log "强制停止服务..."
            kill -9 "$PID"
        fi
    fi
    rm -f "$PID_FILE"
fi

# 6. 启动新服务
log "启动新服务..."
nohup ./main > /var/log/new-api.log 2>&1 &
echo $! > "$PID_FILE"

# 7. 等待服务启动
log "等待服务启动..."
sleep 10

# 8. 检查服务状态
if curl -f http://localhost:4000/api/status > /dev/null 2>&1; then
    log "部署成功！服务已启动"
else
    log "警告: 服务可能未正常启动，请检查日志"
    exit 1
fi

log "部署完成"
