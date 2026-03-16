# SPEC-012: templテンプレート実装

> Issue: #12
> Phase: 1 (ブログエンジン MVP)
> Wave: 2 (コア機能)
> サイズ: M
> 方法論: TDD（新規機能）
> 依存: SPEC-008（データモデル）

---

## 1. 概要

templベースのHTMLテンプレートを実装する。共通レイアウト、記事一覧、
記事詳細、タグ・カテゴリページ、再利用コンポーネントを含む。
HTMX属性を組み込み、ページネーション等の動的UIをサーバーサイドで実現する。

---

## 2. スコープ

### IN

- `templates/layout.templ` — 共通レイアウト（HTML head, header, footer）
- `templates/components/` — 再利用コンポーネント（post_card, pagination, tag_list）
- `templates/pages/` — ページテンプレート（home, post_list, post_detail, tag_page, category_page, error）
- Mermaid.js / highlight.js CDN読み込み
- `templ generate` でコンパイル確認

### OUT

- Tailwind CSSスタイリング（#14 — 構造のみ、classは最小限）
- HTTPハンドラー連携（#13）

---

## 3. 技術設計

### 3.1 依存関係

```go
import "github.com/a-h/templ"
```

templ CLIで `.templ` → `_templ.go` コード生成。

### 3.2 パッケージ構造

```
templates/
├── layout.templ                # 共通レイアウト
├── components/
│   ├── post_card.templ         # 記事カード
│   ├── pagination.templ        # ページネーション（HTMX対応）
│   └── tag_list.templ          # タグ一覧（カウント付き）
└── pages/
    ├── home.templ              # トップページ（最新記事）
    ├── post_list.templ         # 記事一覧
    ├── post_detail.templ       # 記事詳細（Mermaid + コードハイライト）
    ├── tag_page.templ          # タグ別記事一覧
    ├── category_page.templ     # カテゴリ別記事一覧
    └── error.templ             # エラーページ（404, 500）
```

### 3.3 共通レイアウト

```go
templ Layout(title string) {
    <!DOCTYPE html>
    <html lang="ja">
    <head>
        <meta charset="UTF-8"/>
        <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
        <title>{ title } | TechPulse Blog</title>
        <!-- Tailwind CSS CDN（#14で本格導入） -->
        <!-- HTMX -->
        <script src="https://unpkg.com/htmx.org@2"></script>
    </head>
    <body>
        @Header()
        <main>
            { children... }
        </main>
        @Footer()
        <!-- Mermaid.js（条件付きロード） -->
        <!-- highlight.js（条件付きロード） -->
    </body>
    </html>
}
```

### 3.4 HTMX統合

ページネーションコンポーネントにHTMX属性を組み込み:
```go
templ Pagination(result model.PageResult[*model.Post], baseURL string) {
    <nav hx-boost="true">
        // 前ページ・次ページリンク
        // hx-get, hx-target, hx-swap で部分更新
    </nav>
}
```

### 3.5 記事詳細 — Mermaid/コードハイライト対応

```go
templ PostDetail(post *model.Post) {
    @Layout(post.Title) {
        <article>
            <h1>{ post.Title }</h1>
            <div>@templ.Raw(post.HTML)</div>  // パース済みHTML埋め込み
        </article>
        if post.Mermaid {
            <script src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"></script>
            <script>mermaid.initialize({startOnLoad: true});</script>
        }
    }
}
```

**設計判断**:
- `templ.Raw()` — パース済みHTMLをエスケープなしで埋め込み
- Mermaid.jsは`post.Mermaid == true`の場合のみロード（パフォーマンス）
- highlight.jsも同様に条件付きロード

### 3.6 エラーページ

```go
templ ErrorPage(code int, message string) {
    @Layout("エラー") {
        <h1>{ fmt.Sprintf("%d", code) }</h1>
        <p>{ message }</p>
    }
}
```

---

## 4. TDD計画

templテンプレートはコード生成のため、テストは生成後の`_templ.go`に対して実施:

### RED Phase

| テストケース | 検証内容 |
|-------------|---------|
| TestLayout_HTMLヘッダー | title, meta, HTMX scriptタグ |
| TestPostCard_レンダリング | タイトル, 日付, タグ表示 |
| TestPagination_複数ページ | 前/次リンク生成 |
| TestPagination_単一ページ | ナビゲーション非表示 |
| TestPostDetail_Mermaid | Mermaid script条件付きロード |
| TestPostDetail_MermaidなしL | script非ロード |
| TestErrorPage | ステータスコード表示 |

### テスト方式

```go
// templ コンポーネントをバッファにレンダリングしてHTML検証
var buf bytes.Buffer
err := component.Render(context.Background(), &buf)
html := buf.String()
// strings.Contains / regex で検証
```

---

## 5. 品質ゲート

| ゲート | 基準 |
|--------|------|
| templ generate | エラー0 |
| テスト通過 | 全パス |
| カバレッジ | 80%以上 |
| make ci | 全チェック通過 |

---

## 6. 受入条件

- [ ] 共通レイアウト（head, header, footer）が正しくレンダリングされる
- [ ] 記事カード、ページネーション、タグ一覧コンポーネントが動作する
- [ ] 記事一覧・詳細・タグ別・カテゴリ別ページが正しくレンダリングされる
- [ ] Mermaid.jsが条件付きでロードされる
- [ ] HTMX属性がページネーションに組み込まれている
- [ ] エラーページ（404, 500）が動作する
- [ ] `templ generate` でエラーなし
- [ ] make ci 全チェック通過
