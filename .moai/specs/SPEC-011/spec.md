# SPEC-011: ポストサービス実装

> Issue: #11
> Phase: 1 (ブログエンジン MVP)
> Wave: 2 (コア機能)
> サイズ: S
> 方法論: TDD（新規機能）
> 依存: SPEC-010（SQLiteリポジトリ）, SPEC-009（マークダウンパーサー）

---

## 1. 概要

ポスト関連のビジネスロジックをService層に実装する。
Repositoryインターフェースに依存し、コンテンツ同期（ファイル→DB）、
一覧取得、詳細取得を提供する。

---

## 2. スコープ

### IN

- `internal/service/post_service.go` — PostService実装
- `internal/service/post_service_test.go` — テスト（モックリポジトリ）

### OUT

- HTTPハンドラー（#13）
- HTMLキャッシュ（MVP後に検討）

---

## 3. 技術設計

### 3.1 パッケージ構造

```
internal/service/
├── post_service.go       # PostService
└── post_service_test.go  # テスト
```

### 3.2 PostService構造体

```go
type PostService struct {
    repo       repository.PostRepository  // インターフェース依存（DI）
    parser     *markdown.Parser
    contentDir string                     // content/posts/ ディレクトリパス
}

func NewPostService(repo repository.PostRepository, parser *markdown.Parser, contentDir string) *PostService
```

**設計判断**:
- Repositoryインターフェース依存 → テスト時にモック差し替え
- Parser注入 → テスト時に制御可能
- contentDir注入 → テスト時に一時ディレクトリ使用

### 3.3 公開メソッド

```go
// GetBySlug はスラッグでポストを取得します。
func (s *PostService) GetBySlug(slug string) (*model.Post, error)

// List はポスト一覧を取得します。
func (s *PostService) List(opts model.ListOptions) (*model.PageResult[*model.Post], error)

// ListByTag はタグ別ポスト一覧を取得します。
func (s *PostService) ListByTag(tag string, opts model.ListOptions) (*model.PageResult[*model.Post], error)

// ListByCategory はカテゴリ別ポスト一覧を取得します。
func (s *PostService) ListByCategory(category string, opts model.ListOptions) (*model.PageResult[*model.Post], error)

// GetAllTags は全タグをカウント付きで取得します。
func (s *PostService) GetAllTags() ([]model.Tag, error)

// GetAllCategories は全カテゴリをカウント付きで取得します。
func (s *PostService) GetAllCategories() ([]model.Category, error)

// SyncContent はcontent/posts/のマークダウンファイルをDBに同期します。
func (s *PostService) SyncContent() error
```

### 3.4 SyncContentロジック

```
1. contentDirの*.mdファイルを走査
2. 各ファイルをparser.ParseFileでPost構造体に変換
3. ファイル名からSlugを生成（拡張子除去、ケバブケース）
4. repo.Upsertで保存
5. DBの全slugを取得
6. ファイルに対応しないslugをrepo.Deleteで削除（孤立レコード除去）
```

**Slug生成規則**: `hello-world.md` → `hello-world`

**削除同期**: ファイルが削除された場合、DB上の対応レコードも削除される。
これによりファイルシステムとDBの整合性を保証する。

### 3.5 ListOptionsバリデーション

Service層でListOptions.Validate()を呼び出し、デフォルト値を設定:
- Page未設定 → 1
- PageSize未設定 → DefaultPageSize (10)

### 3.6 エラーラッピング

全メソッドでエラーラッピング必須（CODING_RULES準拠）:
```go
return nil, fmt.Errorf("ポスト一覧取得失敗: %w", err)
```

---

## 4. TDD計画

### RED Phase

| テストケース | 検証内容 |
|-------------|---------|
| TestGetBySlug_正常系 | モックリポジトリから取得 |
| TestGetBySlug_存在しない | ErrPostNotFound伝播 |
| TestList_正常系 | 一覧取得 + PageResult |
| TestList_デフォルトオプション | Page=0→1, PageSize=0→10 |
| TestListByTag | タグ別一覧 |
| TestListByCategory | カテゴリ別一覧 |
| TestGetAllTags | タグ一覧 |
| TestGetAllCategories | カテゴリ一覧 |
| TestSyncContent_新規追加 | ファイル同期（一時ディレクトリ + 実Parser） |
| TestSyncContent_削除同期 | DBにあるがファイルにないレコードが削除される |

### モックリポジトリ

```go
type mockPostRepository struct {
    posts map[string]*model.Post
    // 各メソッドのスタブ
}
```

---

## 5. 品質ゲート

| ゲート | 基準 |
|--------|------|
| テスト通過 | 全パス |
| カバレッジ | 80%以上 |
| make ci | 全チェック通過 |
| 逆方向依存 | handler/repository をimportしない |

---

## 6. 受入条件

- [ ] PostServiceが全公開メソッドを実装している
- [ ] Repositoryインターフェース依存（DI）
- [ ] SyncContentでcontent/posts/*.mdをDB同期できる
- [ ] ListOptionsのデフォルト値設定が動作する
- [ ] エラーラッピングが全メソッドに適用されている
- [ ] モックリポジトリを使用したテストが全パス
- [ ] make ci 全チェック通過
