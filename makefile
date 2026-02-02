.PHONY: swag wire gen all clean help

# 生成 Swagger API 文档
swag:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/server/main.go -o ./docs --parseDependency --parseInternal --exclude runtime,vendor/test --generatedTime

# 生成 Wire 依赖注入代码
wire:
	@echo "Generating Wire dependency injection code..."
	cd cmd && wire

# 一次性生成所有代码（swag + wire）
gen: wire swag
	@echo "All code generation completed!"

# 清理生成的文件
clean:
	@echo "Cleaning generated files..."
	rm -rf docs/
	rm -f cmd/wire_gen.go

# 显示帮助信息
help:
	@echo "Available commands:"
	@echo "  make swag    - Generate Swagger API documentation"
	@echo "  make wire    - Generate Wire dependency injection code"
	@echo "  make gen     - Generate all (wire + swag)"
	@echo "  make clean   - Clean generated files"
	@echo "  make help    - Show this help message"

# 默认目标
all: gen
