#!/bin/bash

# 开发环境启动脚本 - 支持热重载

echo "🚀 启动 One-API 开发环境 (热重载模式)"

# 检查 air 是否已安装
if ! command -v air &> /dev/null; then
    echo "❌ air 未安装，正在安装..."
    go install github.com/air-verse/air@latest
fi

# 创建临时目录
mkdir -p tmp

# 检查环境变量文件
if [ ! -f .env ]; then
    echo "⚠️  未找到 .env 文件，请确保已配置环境变量"
fi

# 启动热重载开发服务器
echo "📝 使用 air 启动热重载开发服务器..."
echo "🔧 配置文件: .air.toml"
echo "📁 监听目录: 当前目录"
echo "🔄 热重载已启用 - 修改代码后会自动重启"
echo ""

air


