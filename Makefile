.PHONY: build run test lint fmt clean dev

# ビルド設定
BINARY_NAME=server
BUILD_DIR=bin
CMD_DIR=cmd/server

# ビルド
build:
	@echo "ビルド中..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)

# 実行
run: build
	@echo "サーバー起動..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# 開発モード（air使用）
dev:
	@air

# テスト
test:
	go test -race -cover ./...

# テストカバレッジ
coverage:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# リント
lint:
	golangci-lint run

# フォーマット
fmt:
	gofmt -w .
	goimports -w .

# templ生成
templ:
	templ generate

# クリーン
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
