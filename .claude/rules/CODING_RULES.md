# コーディングルール

> **関連文書**: [CODE_STYLE.md](./CODE_STYLE.md) | [PR_CHECKLIST.md](./PR_CHECKLIST.md)

TechPulse Blogプロジェクトのアーキテクチャ、パターン、パフォーマンス、エラー処理ルールを定義します。

---

## 1. アーキテクチャルール

### 1.1 レイヤー分離原則

```
cmd/server/main.go          ← エントリポイント
    |
internal/handler/            ← HTTPハンドラー（リクエスト/レスポンス）
    |
internal/service/            ← ビジネスロジック
    |
internal/repository/         ← データアクセス（SQLite、ファイルシステム）
    |
internal/model/              ← データモデル（構造体）
```

| 階層 | パッケージ | 許可依存性 |
|------|--------|-------------|
| Model | `internal/model` | 標準ライブラリのみ |
| Repository | `internal/repository` | model |
| Service | `internal/service` | model, repository |
| Handler | `internal/handler` | model, service |
| Markdown | `internal/markdown` | model, goldmark |
| Search | `internal/search` | model, bleve |

### 1.2 依存性ルール

- **下位 → 上位依存禁止**: modelがhandlerをimportしてはならない
- **インターフェース基盤依存性注入**: serviceはrepositoryインターフェースに依存
- **循環依存禁止**: パッケージ間循環import不可（Goコンパイラがブロック）

```go
// 準拠 - インターフェース基盤
type PostRepository interface {
    FindBySlug(slug string) (*model.Post, error)
    FindAll(opts ListOptions) ([]*model.Post, error)
}

type PostService struct {
    repo PostRepository  // インターフェース依存
}
```

---

## 2. エラー処理ルール

### 2.1 エラーラッピング必須

```go
// 違反
return err

// 準拠
return fmt.Errorf("ポスト照会失敗 (slug=%s): %w", slug, err)
```

### 2.2 カスタムエラータイプ

```go
// ドメインエラー定義
var (
    ErrPostNotFound = errors.New("ポストが見つかりません")
    ErrInvalidSlug  = errors.New("無効なスラッグです")
)
```

### 2.3 HTTPエラーレスポンス

```go
// handlerレイヤーでのみHTTPステータスコード決定
func (h *PostHandler) GetPost(c echo.Context) error {
    post, err := h.service.GetBySlug(slug)
    if err != nil {
        if errors.Is(err, service.ErrPostNotFound) {
            return c.Render(http.StatusNotFound, "404", nil)
        }
        return c.Render(http.StatusInternalServerError, "500", nil)
    }
    return c.Render(http.StatusOK, "post", post)
}
```

---

## 3. パフォーマンスルール

### 3.1 マークダウンパースキャッシュ

- パース済みHTMLはメモリキャッシュに保管
- ファイル変更検知時キャッシュ無効化
- 最大キャッシュサイズ設定

### 3.2 データベース

- SQLite WALモード使用（読み取り性能向上）
- インデックス: slug, date, tags, category
- クエリにLIMIT/OFFSET必須

### 3.3 HTTPパフォーマンス

- 静的ファイルキャッシュヘッダー設定
- gzip圧縮ミドルウェア適用
- templコンポーネントレンダリング（ストリーミング）

---

## 4. テストルール

### 4.1 テスト構造

```
internal/
    service/
        post_service.go
        post_service_test.go    ← 同一パッケージにテスト
    handler/
        post_handler.go
        post_handler_test.go
```

### 4.2 テーブルドリブンテスト

```go
func TestParseMarkdown(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    *Post
        wantErr bool
    }{
        {
            name:  "正常なポスト",
            input: "---\ntitle: test\n---\n# Hello",
            want:  &Post{Title: "test"},
        },
        {
            name:    "frontmatterなし",
            input:   "# No frontmatter",
            wantErr: true,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseMarkdown([]byte(tt.input))
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseMarkdown() error = %v, wantErr %v", err, tt.wantErr)
            }
            if !tt.wantErr && got.Title != tt.want.Title {
                t.Errorf("ParseMarkdown() title = %v, want %v", got.Title, tt.want.Title)
            }
        })
    }
}
```

---

## 5. Quick Reference

| 項目 | 準拠 | 違反 |
|------|------|------|
| 依存性方向 | handler → service → repo → model | model → handler |
| エラー処理 | `fmt.Errorf("...: %w", err)` | `return err`（ラッピングなし） |
| エラー無視 | 明示的 `_ =` + コメント | 暗黙的無視 |
| HTTPステータス | handlerでのみ決定 | serviceでHTTPコード返却 |
| SQLクエリ | Prepared statement | 文字列連結 |
| テスト | テーブルドリブン | 単一ケース |

### 便利なコマンド

```bash
# コードフォーマット
gofmt -w .
goimports -w .

# 静的解析
golangci-lint run

# テスト実行
go test ./...

# テストカバレッジ
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# templコード生成
templ generate

# ビルド
go build -o bin/server ./cmd/server
```
