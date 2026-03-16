.PHONY: build run test lint fmt clean dev ci ci-quick

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

# ローカルCI（GitHub Actions相当の全チェック）
ci: fmt lint build test
	@echo ""
	@echo "=== ローカルCI 全チェック通過 ==="

# ローカルCI（高速版：フォーマット確認 + vet + テストのみ）
ci-quick:
	@echo "=== フォーマット確認 ==="
	@test -z "$$(gofmt -l .)" || (echo "フォーマットエラー:"; gofmt -l .; exit 1)
	@echo "=== go vet ==="
	go vet ./...
	@echo "=== テスト ==="
	go test -race -cover ./...
	@echo ""
	@echo "=== ローカルCI（高速版）通過 ==="

# クリーン
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
