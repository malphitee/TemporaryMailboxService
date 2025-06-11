# 临时邮箱系统 Makefile (Windows 版本)
# 简化版本，兼容Windows环境

# 变量定义
APP_NAME := temp-mailbox-service
BACKEND_DIR := backend
BIN_DIR := $(BACKEND_DIR)/bin

.PHONY: help
help: ## 显示帮助信息
	@echo "临时邮箱系统 - 可用命令:"
	@echo ""
	@echo "开发命令:"
	@echo "  setup          初始化开发环境"
	@echo "  dev            启动开发服务器（热重载）"
	@echo "  run            直接运行服务器"
	@echo "  worker         运行后台任务服务"
	@echo ""
	@echo "构建命令:"
	@echo "  build          构建所有二进制文件"
	@echo "  build-server   构建服务器二进制文件"
	@echo "  build-worker   构建后台任务二进制文件"
	@echo ""
	@echo "测试命令:"
	@echo "  test           运行所有测试"
	@echo "  test-coverage  运行测试并生成覆盖率报告"
	@echo ""
	@echo "数据库命令:"
	@echo "  db-migrate     执行数据库迁移"
	@echo "  db-seed        填充测试数据"
	@echo "  db-reset       重置数据库"
	@echo ""
	@echo "代码质量:"
	@echo "  fmt            格式化代码"
	@echo "  vet            运行静态分析"
	@echo "  tidy           整理依赖"
	@echo ""
	@echo "Docker命令:"
	@echo "  docker-up      启动Docker开发环境"
	@echo "  docker-down    停止Docker开发环境"
	@echo ""
	@echo "工具命令:"
	@echo "  gen-env        生成环境配置文件"
	@echo "  clean          清理构建文件"
	@echo "  install-tools  安装开发工具"

##@ 开发命令

.PHONY: setup
setup: ## 初始化开发环境
	@echo "🔧 初始化开发环境..."
	@cd $(BACKEND_DIR) && go mod download
	@cd $(BACKEND_DIR) && go mod tidy
	@echo "✅ 开发环境初始化完成"

.PHONY: dev
dev: ## 启动开发服务器（热重载）
	@echo "🚀 启动开发服务器..."
	@cd $(BACKEND_DIR) && air

.PHONY: run
run: ## 直接运行服务器
	@echo "🚀 启动服务器..."
	@cd $(BACKEND_DIR) && go run cmd/server/main.go

.PHONY: worker
worker: ## 运行后台任务服务
	@echo "⚙️ 启动后台任务服务..."
	@cd $(BACKEND_DIR) && go run cmd/worker/main.go

##@ 构建命令

.PHONY: build
build: build-server build-worker ## 构建所有二进制文件

.PHONY: build-server
build-server: ## 构建服务器二进制文件
	@echo "🔨 构建服务器..."
	@if not exist "$(BIN_DIR)" mkdir "$(BIN_DIR)"
	@cd $(BACKEND_DIR) && go build -o bin/server.exe cmd/server/main.go
	@echo "✅ 服务器构建完成: $(BIN_DIR)/server.exe"

.PHONY: build-worker
build-worker: ## 构建后台任务二进制文件
	@echo "🔨 构建后台任务服务..."
	@if not exist "$(BIN_DIR)" mkdir "$(BIN_DIR)"
	@cd $(BACKEND_DIR) && go build -o bin/worker.exe cmd/worker/main.go
	@echo "✅ 后台任务服务构建完成: $(BIN_DIR)/worker.exe"

##@ 测试命令

.PHONY: test
test: ## 运行所有测试
	@echo "🧪 运行测试..."
	@cd $(BACKEND_DIR) && go test ./... -v

.PHONY: test-coverage
test-coverage: ## 运行测试并生成覆盖率报告
	@echo "📊 生成测试覆盖率报告..."
	@cd $(BACKEND_DIR) && go test ./... -coverprofile=coverage.out
	@cd $(BACKEND_DIR) && go tool cover -html=coverage.out -o coverage.html
	@echo "✅ 覆盖率报告生成完成: $(BACKEND_DIR)/coverage.html"

##@ 数据库命令

.PHONY: db-migrate
db-migrate: ## 执行数据库迁移
	@echo "🗄️ 执行数据库迁移..."
	@cd $(BACKEND_DIR) && go run scripts/migrate.go

.PHONY: db-seed
db-seed: ## 填充测试数据
	@echo "🌱 填充测试数据..."
	@cd $(BACKEND_DIR) && go run scripts/seed.go

.PHONY: db-reset
db-reset: ## 重置数据库（删除并重新迁移）
	@echo "🔄 重置数据库..."
	@if exist "$(BACKEND_DIR)\dev.db" del "$(BACKEND_DIR)\dev.db"
	@$(MAKE) db-migrate
	@$(MAKE) db-seed

##@ 代码质量

.PHONY: fmt
fmt: ## 格式化代码
	@echo "💅 格式化代码..."
	@cd $(BACKEND_DIR) && go fmt ./...
	@echo "✅ 代码格式化完成"

.PHONY: vet
vet: ## 运行go vet检查
	@echo "🔍 运行静态分析..."
	@cd $(BACKEND_DIR) && go vet ./...

.PHONY: tidy
tidy: ## 整理依赖
	@echo "📦 整理依赖..."
	@cd $(BACKEND_DIR) && go mod tidy

##@ Docker命令

.PHONY: docker-up
docker-up: ## 启动Docker开发环境
	@echo "🐳 启动Docker开发环境..."
	@docker-compose -f docker-compose.dev.yml up -d

.PHONY: docker-down
docker-down: ## 停止Docker开发环境
	@echo "🐳 停止Docker开发环境..."
	@docker-compose -f docker-compose.dev.yml down

.PHONY: docker-logs
docker-logs: ## 查看Docker日志
	@docker-compose -f docker-compose.dev.yml logs -f

##@ 工具命令

.PHONY: install-tools
install-tools: ## 安装开发工具
	@echo "🔧 安装开发工具..."
	@go install github.com/cosmtrek/air@latest
	@echo "✅ 开发工具安装完成"

.PHONY: gen-env
gen-env: ## 生成环境配置文件
	@echo "📝 生成环境配置文件..."
	@if not exist ".env.dev" ( \
		copy ".env.example" ".env.dev" >nul 2>&1 && \
		echo "✅ 创建 .env.dev 文件" \
	) else ( \
		echo "⚠️  .env.dev 文件已存在" \
	)

.PHONY: clean
clean: ## 清理构建文件
	@echo "🧹 清理构建文件..."
	@if exist "$(BIN_DIR)" rmdir /s /q "$(BIN_DIR)" 2>nul
	@if exist "$(BACKEND_DIR)\coverage.out" del "$(BACKEND_DIR)\coverage.out" 2>nul
	@if exist "$(BACKEND_DIR)\coverage.html" del "$(BACKEND_DIR)\coverage.html" 2>nul
	@echo "✅ 清理完成"

##@ 快速命令组合

.PHONY: quick-start
quick-start: gen-env install-tools setup ## 快速启动开发环境
	@echo "🎉 开发环境已准备就绪！"
	@echo "💡 提示:"
	@echo "   - 使用 'make docker-up' 启动数据库等服务"
	@echo "   - 使用 'make dev' 启动热重载开发服务器"
	@echo "   - 使用 'make run' 直接运行服务器"

.PHONY: quick-test
quick-test: fmt vet test ## 快速测试（格式化 + 检查 + 测试）

# 默认目标
.DEFAULT_GOAL := help 