FRONTEND_DIR = ./web
BACKEND_DIR = .

.PHONY: all build-frontend start-backend dev hot-reload

all: build-frontend start-backend

build-frontend:
	@echo "Building frontend..."
	@cd $(FRONTEND_DIR) && bun install && DISABLE_ESLINT_PLUGIN='true' VITE_REACT_APP_VERSION=$(cat VERSION) bun run build

start-backend:
	@echo "Starting backend dev server..."
	@cd $(BACKEND_DIR) && go run main.go &

# 热重载开发模式
dev: hot-reload

hot-reload:
	@echo "🚀 启动热重载开发模式..."
	@echo "📝 使用 air 进行热重载开发"
	@echo "🔄 修改代码后会自动重启服务器"
	@echo ""
	@cd $(BACKEND_DIR) && ./dev.sh
