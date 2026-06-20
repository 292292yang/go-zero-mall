# 配置项，按需修改
SERVICE_NAME=api
API_FILE=api/user.api
CONF_PATH=etc/$(SERVICE_NAME).yaml
BIN_NAME=$(SERVICE_NAME)

# Go 环境变量
GO=go
GOCTL=goctl

.PHONY: all api run build tidy fmt clean help

# 默认目标：打印帮助
all: help

# 1. 根据 .api 文件生成 go-zero 代码
api:
	$(GOCTL) api go -api $(API_FILE) -dir .

#model:
#	goctl model mysql ddl -src sql/product.sql -dir ./internal/model -cache

# 2. 下载/整理依赖
tidy:
	$(GO) mod tidy

# 3. 代码格式化
fmt:
	$(GOCTL) api format --dir .
	$(GO) fmt ./...

# 4. 本地直接运行服务
run: tidy fmt
	$(GO) run $(SERVICE_NAME).go -f $(CONF_PATH)

# 5. 编译生成二进制文件
build: tidy fmt
	CGO_ENABLED=0 $(GO) build -o $(BIN_NAME) $(SERVICE_NAME).go

# 6. 清理产物
clean:
	rm -f $(BIN_NAME)
	$(GO) clean

# 帮助文档
help:
	@echo "Usage:"
	@echo "  make api      - 基于api文件生成go-zero代码"
	@echo "  make tidy     - 拉取/整理go module依赖"
	@echo "  make fmt      - 格式化api与go代码"
	@echo "  make run      - 本地启动服务"
	@echo "  make build    - 静态编译二进制程序"
	@echo "  make clean    - 删除二进制、清理编译缓存"