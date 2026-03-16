# コードスタイルルール

> **関連文書**: [CODING_RULES.md](./CODING_RULES.md) | [PR_CHECKLIST.md](./PR_CHECKLIST.md)

---

## 1. コメント言語

- **全てのコードコメントは日本語で作成**
- 関数/構造体/インターフェースドキュメントコメントを含む
- インラインコメントを含む
- 例外: 外部ライブラリインターフェース実装時は英文維持可能

```go
// PostServiceはブログポストのCRUD操作を処理します。
//
// マークダウンパース、frontmatter抽出、検索インデキシングを含みます。
type PostService struct {
    // ポスト保存所
    repo PostRepository
    // 全文検索エンジン
    search *bleve.Index
}

// GetBySlugはスラッグでポストを照会します。
// ポストがない場合はErrPostNotFoundを返します。
func (s *PostService) GetBySlug(slug string) (*Post, error) {
    // キャッシュから先に照会
    if cached, ok := s.cache.Get(slug); ok {
        return cached, nil
    }
    return s.repo.FindBySlug(slug)
}
```

---

## 2. ネーミング規則

| 対象 | 規則 | 例 |
|------|------|------|
| パッケージ | lowercase | `handler`, `service`, `model` |
| 構造体/インターフェース | PascalCase | `PostService`, `SearchEngine` |
| 公開関数/メソッド | PascalCase | `GetBySlug`, `ParseMarkdown` |
| 非公開関数/メソッド | camelCase | `parseContent`, `buildIndex` |
| 定数 | PascalCaseまたはcamelCase | `MaxPageSize`, `defaultTimeout` |
| ファイル | snake_case | `post_handler.go`, `feed_service.go` |
| テストファイル | snake_case + _test | `post_handler_test.go` |

---

## 3. Goコードフォーマット

- `gofmt` / `goimports` 必須適用
- `golangci-lint` 静的解析通過

---

## 4. 品質ルール

### ハードコーディング禁止

```go
// 違反
if pageSize > 20 { ... }

// 準拠
const MaxPageSize = 20
if pageSize > MaxPageSize { ... }
```

### エラー処理必須

```go
// 違反
result, _ := doSomething()

// 準拠
result, err := doSomething()
if err != nil {
    return fmt.Errorf("操作失敗: %w", err)
}
```

### 公開API文書化

```go
// 違反 - ドキュメントなし
func ParseMarkdown(content []byte) (*Post, error) { ... }

// 準拠
// ParseMarkdownはマークダウンコンテンツをパースしてPost構造体に変換します。
//
// frontmatter(YAML)と本文を分離して処理します。
// Mermaidコードブロックはクライアントレンダリング用に保存されます。
func ParseMarkdown(content []byte) (*Post, error) { ... }
```
