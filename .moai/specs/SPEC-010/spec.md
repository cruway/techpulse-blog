# SPEC-010: SQLiteリポジトリ実装

> Issue: #10
> Phase: 1 (ブログエンジン MVP)
> Wave: 2 (コア機能)
> サイズ: M
> 方法論: TDD（新規機能）
> 依存: SPEC-008（データモデル）, SPEC-009（マークダウンパーサー）

---

## 1. 概要

SQLiteを使用したPostRepositoryを実装する。WALモード、インデックス設定、
CRUD操作、ページネーション、マークダウンファイルからのDB同期を含む。
CGO不要のmodernc.org/sqliteドライバを使用する。

---

## 2. スコープ

### IN

- `internal/repository/post_repository.go` — PostRepositoryインターフェース
- `internal/repository/sqlite.go` — SQLite実装（接続、マイグレーション）
- `internal/repository/sqlite_post.go` — PostRepository SQLite実装
- `internal/repository/sqlite_post_test.go` — テスト（インメモリSQLite）

### OUT

- コンテンツ同期ロジック（#11 ポストサービス）
- 全文検索Bleve連携（Phase 4）

---

## 3. 技術設計

### 3.1 依存関係

```go
import (
    "database/sql"
    _ "modernc.org/sqlite"  // ドライバ名: "sqlite"
)

// 接続例
db, err := sql.Open("sqlite", dbPath)
```

**選定理由**: Oracle Cloud Always Free環境でCGOビルドが困難。modernc.org/sqliteは純Go実装。
**ドライバ登録名**: `"sqlite"`（`_ "modernc.org/sqlite"` のimportで自動登録）

### 3.2 パッケージ構造

```
internal/repository/
├── post_repository.go      # PostRepositoryインターフェース
├── sqlite.go               # DB接続、マイグレーション
├── sqlite_post.go          # PostRepository SQLite実装
└── sqlite_post_test.go     # テスト
```

### 3.3 PostRepositoryインターフェース

```go
type PostRepository interface {
    FindBySlug(slug string) (*model.Post, error)
    FindAll(opts model.ListOptions) (*model.PageResult[*model.Post], error)
    FindByTag(tag string, opts model.ListOptions) (*model.PageResult[*model.Post], error)
    FindByCategory(category string, opts model.ListOptions) (*model.PageResult[*model.Post], error)
    AllTags() ([]model.Tag, error)
    AllCategories() ([]model.Category, error)
    Upsert(post *model.Post) error
    Delete(slug string) error
}
```

**設計判断**:
- `PageResult[*model.Post]` を返却 → Service層でのページネーション計算不要
- `Upsert` → ファイル同期時にINSERT or UPDATE
- `Delete` 追加 → ファイル削除時のDB同期

### 3.4 テーブルスキーマ

```sql
CREATE TABLE IF NOT EXISTS posts (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    title       TEXT NOT NULL,
    slug        TEXT NOT NULL UNIQUE,
    content     TEXT NOT NULL,
    html        TEXT NOT NULL,
    excerpt     TEXT NOT NULL DEFAULT '',
    date        DATETIME NOT NULL,
    category    TEXT NOT NULL DEFAULT '',
    status      TEXT NOT NULL DEFAULT 'draft',
    source_url  TEXT NOT NULL DEFAULT '',
    mermaid     BOOLEAN NOT NULL DEFAULT 0,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS post_tags (
    post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    tag     TEXT NOT NULL,
    PRIMARY KEY (post_id, tag)
);

CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug);
CREATE INDEX IF NOT EXISTS idx_posts_date ON posts(date DESC);
CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status);
CREATE INDEX IF NOT EXISTS idx_posts_category ON posts(category);
CREATE INDEX IF NOT EXISTS idx_post_tags_tag ON post_tags(tag);
```

**設計判断**:
- タグを別テーブル（post_tags）に正規化 → タグ別検索のパフォーマンス
- `slug` にUNIQUE制約 → Upsert時のON CONFLICTで使用
- WALモード → 読み取り並行性向上

### 3.5 DB初期化

```go
func NewDB(dbPath string) (*sql.DB, error)  // WALモード + マイグレーション実行
func NewTestDB() (*sql.DB, error)           // インメモリDB（テスト用）
```

### 3.6 Prepared Statement

全クエリでPrepared Statementを使用（SQL injection防止 — CODING_RULES準拠）。

---

## 4. TDD計画

### RED Phase

| テストケース | 検証内容 |
|-------------|---------|
| TestUpsert_新規作成 | Insert動作確認 |
| TestUpsert_更新 | 同一slugで更新 |
| TestFindBySlug_正常系 | slug検索 |
| TestFindBySlug_存在しない | ErrPostNotFound |
| TestFindAll_ページネーション | Page/PageSize正しく動作 |
| TestFindAll_ステータスフィルタ | publishedのみ取得 |
| TestFindByTag | タグ別検索 |
| TestFindByCategory | カテゴリ別検索 |
| TestAllTags | タグ一覧 + カウント |
| TestAllCategories | カテゴリ一覧 + カウント |
| TestDelete | 削除 + カスケード（post_tags） |

---

## 5. 品質ゲート

| ゲート | 基準 |
|--------|------|
| テスト通過 | 全パス |
| カバレッジ | 80%以上 |
| make ci | 全チェック通過 |
| SQL injection | Prepared Statement 100% |

---

## 6. 受入条件

- [ ] PostRepositoryインターフェースが定義されている
- [ ] SQLite WALモードで初期化される
- [ ] CRUD操作（Upsert, FindBySlug, FindAll, Delete）が動作する
- [ ] ページネーション（Page, PageSize, TotalItems, TotalPages）が正確
- [ ] タグ別・カテゴリ別検索が動作する
- [ ] AllTags/AllCategoriesがカウント付きで返却される
- [ ] テストはインメモリSQLiteで実行（外部依存なし）
- [ ] make ci 全チェック通過
