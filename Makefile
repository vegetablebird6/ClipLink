.PHONY: build run test help

# 默认目标
all: build

# 编译项目
build:
	go build -o bin/cliplink ./cmd/main.go

# 运行项目
run:
	go run ./cmd/main.go

# 测试
test:
	go test -v ./cmd/... ./internal/...

# 帮助信息
help:
	@echo "可用的命令:"
	@echo "  make build      - 编译项目"
	@echo "  make run        - 运行项目"
	@echo "  make test       - 运行测试"
	@echo "  make help       - 显示此帮助信息" 
