GOCMD=go
GOBUILD=$(GOCMD) build
BINARY_NAME=objs/srs-sip
MAIN_PATH=main/main.go
VUE_DIR=html/NextGB

default: build

build:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)

clean:
	rm -f $(BINARY_NAME)
	rm -rf $(VUE_DIR)/dist
	rm -rf $(VUE_DIR)/node_modules

run:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)
	./$(BINARY_NAME)

install:
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)
	mv $(BINARY_NAME) /usr/local/bin	

vue-install:
	cd $(VUE_DIR) && npm install

vue-build:
	cd $(VUE_DIR) && npm run build

vue-dev:
	cd $(VUE_DIR) && npm run dev

all: build vue-build

# 测试相关
test:
	$(GOCMD) test -v ./...

test-coverage:
	$(GOCMD) test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

test-short:
	$(GOCMD) test -short -v ./...

# 代码质量检查
lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

# 格式化代码
fmt:
	$(GOCMD) fmt ./...
	goimports -w .

# 检查代码
vet:
	$(GOCMD) vet ./...

# CI 相关
ci: lint test build

# 帮助信息
help:
	@echo "可用的 make 命令："
	@echo "  make build          - 构建后端二进制文件"
	@echo "  make clean          - 清理构建产物"
	@echo "  make run            - 构建并运行"
	@echo "  make test           - 运行所有测试"
	@echo "  make test-coverage  - 运行测试并生成覆盖率报告"
	@echo "  make lint           - 运行代码质量检查"
	@echo "  make lint-fix       - 运行代码质量检查并自动修复"
	@echo "  make fmt            - 格式化代码"
	@echo "  make vet            - 运行 go vet"
	@echo "  make vue-install    - 安装前端依赖"
	@echo "  make vue-build      - 构建前端"
	@echo "  make vue-dev        - 启动前端开发服务器"
	@echo "  make all            - 构建前后端"
	@echo "  make ci             - 运行 CI 流程 (lint + test + build)"

.PHONY: clean vue-install vue-build vue-dev all test test-coverage test-short lint lint-fix fmt vet ci help
