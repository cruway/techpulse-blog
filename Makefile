.PHONY: build run test lint fmt clean dev ci ci-quick coverage setup hooks coverage-check

# ビルド設定
BINARY_NAME=server
BUILD_DIR=bin
CMD_DIR=cmd/server
COVERAGE_THRESHOLD=80

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

# テストカバレッジ（HTMLレポート生成）
coverage:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# カバレッジ閾値チェック（パッケージ単位で検証、cmd/除外）
coverage-check:
	@echo "=== カバレッジ閾値チェック ($(COVERAGE_THRESHOLD)%) ==="
	@FAIL=0; \
	for pkg in $$(go test -race -cover ./... 2>&1 | grep -v 'no test files' | grep -v 'cmd/'); do \
		COV=$$(echo "$$pkg" | grep -oE '[0-9]+\.[0-9]+%' | sed 's/%//'); \
		PKG=$$(echo "$$pkg" | awk '{print $$2}'); \
		if [ -n "$$COV" ] && [ -n "$$PKG" ]; then \
			if [ $$(echo "$$COV < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
				echo "  NG: $$PKG — $$COV% < $(COVERAGE_THRESHOLD)%"; \
				FAIL=1; \
			else \
				echo "  OK: $$PKG — $$COV%"; \
			fi; \
		fi; \
	done; \
	if [ $$FAIL -eq 1 ]; then echo ""; echo "カバレッジ閾値未達のパッケージがあります"; exit 1; fi
	@echo ""
	@echo "全パッケージ カバレッジOK"

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

# ローカルCI（GitHub Actions相当の全チェック + カバレッジ閾値）
ci: fmt lint build test coverage-check
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

# 開発ツール一括インストール
setup:
	@echo "=== 開発ツールインストール ==="
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/air-verse/air@latest
	go install github.com/a-h/templ/cmd/templ@latest
	@echo ""
	@echo "=== pre-commitフックインストール ==="
	@$(MAKE) hooks
	@echo ""
	@echo "=== セットアップ完了 ==="

# pre-commitフックインストール
hooks:
	@mkdir -p .git/hooks
	@echo '#!/bin/sh' > .git/hooks/pre-commit
	@echo '# MoAI-ADK ローカルCI pre-commitフック' >> .git/hooks/pre-commit
	@echo 'echo "=== pre-commit: ローカルCI実行中 ==="' >> .git/hooks/pre-commit
	@echo 'make ci-quick' >> .git/hooks/pre-commit
	@echo 'if [ $$? -ne 0 ]; then' >> .git/hooks/pre-commit
	@echo '  echo "pre-commit: チェック失敗。コミットを中止します。"' >> .git/hooks/pre-commit
	@echo '  exit 1' >> .git/hooks/pre-commit
	@echo 'fi' >> .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "pre-commitフックをインストールしました"

# クリーン
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
