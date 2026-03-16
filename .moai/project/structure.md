# TechPulse Blog — Structure Document

> 生成日: 2026-03-16
> ソース基盤: コードベース自動分析

---

## 1. プロジェクト構造

```
techpulse-blog/
├── cmd/
│   └── server/
│       ├── main.go              # エントリポイント（Echo起動）
│       └── main_test.go         # サーバーテスト
├── internal/
│   ├── handler/                 # HTTPハンドラー（リクエスト/レスポンス）
│   ├── service/                 # ビジネスロジック
│   ├── repository/              # データアクセス（SQLite、ファイルシステム）
│   ├── model/                   # データモデル（構造体）
│   ├── markdown/                # マークダウンパース（goldmark）
│   └── search/                  # 全文検索（Bleve）
├── templates/                   # templテンプレート
├── static/
│   ├── css/                     # Tailwind CSS
│   ├── js/                      # HTMX等
│   └── img/                     # 画像
├── content/
│   ├── posts/                   # 公開記事（マークダウン）
│   ├── drafts/                  # 下書き
│   ├── inbox/                   # Obsidian受信箱
│   ├── knowledge/               # ナレッジベース
│   └── config/                  # コンテンツ設定
├── deploy/                      # デプロイ設定
├── .claude/                     # Claude Code設定
│   ├── rules/                   # コーディングルール、Gitワークフロー
│   └── CLAUDE.md                # プロジェクトルール
├── .moai/                       # MoAI-ADK設定
│   ├── config/sections/         # 設定ファイル群
│   ├── project/                 # プロジェクト文書（このファイル）
│   └── specs/                   # SPEC文書
├── .github/
│   ├── workflows/ci.yml         # CI設定
│   └── ISSUE_TEMPLATE/          # Issueテンプレート
├── go.mod                       # Go依存関係
├── Makefile                     # ビルドコマンド
├── .air.toml                    # ホットリロード設定
└── .golangci.yml                # 静的解析設定
```

---

## 2. レイヤー依存性アーキテクチャ

```
┌──────────────────────────────────────────┐
│  cmd/server/main.go                       │
│  └─ エントリポイント、DI構成              │
├──────────────────────────────────────────┤
│  internal/handler/                        │
│  └─ HTTPハンドラー、ルーティング          │
├──────────────────────────────────────────┤
│  internal/service/                        │
│  └─ ビジネスロジック、ユースケース        │
├──────────────────────────────────────────┤
│  internal/repository/                     │
│  └─ データアクセス（SQLite、ファイルI/O） │
│                                           │
│  internal/markdown/                       │
│  └─ マークダウンパース（goldmark）        │
│                                           │
│  internal/search/                         │
│  └─ 全文検索（Bleve）                     │
├──────────────────────────────────────────┤
│  internal/model/                          │
│  └─ データモデル（構造体のみ）            │
└──────────────────────────────────────────┘
     ↑ 依存性方向（上位 → 下位のみ許可）
```

### 依存性マトリックス

| import ↓ \ from → | model | repository | markdown | search | service | handler |
|--------------------|-------|------------|----------|--------|---------|---------|
| **model** | - | X | X | X | X | X |
| **repository** | O | - | X | X | X | X |
| **markdown** | O | X | - | X | X | X |
| **search** | O | X | X | - | X | X |
| **service** | O | O | O | O | - | X |
| **handler** | O | X | X | X | O | - |

---

## 3. パッケージ別責務

### cmd/server（エントリポイント）

- Echoサーバー初期化
- ミドルウェア設定（Logger, Recover, Gzip）
- DI構成（Repository → Service → Handler）
- グレースフルシャットダウン

### internal/model（データモデル）

- Post、Tag、Category構造体
- 標準ライブラリのみ依存
- ドメインエラー定義

### internal/repository（データアクセス）

- PostRepository インターフェース + SQLite実装
- CRUD + ページネーション
- Prepared Statement必須

### internal/service（ビジネスロジック）

- PostService（CRUD、キャッシュ）
- Repositoryインターフェース依存（DI）
- エラーラッピング必須

### internal/handler（HTTPハンドラー）

- Echo Handlerメソッド
- HTTPステータスコード決定はここのみ
- Service依存

### internal/markdown（マークダウン）

- goldmarkベースパーサー
- frontmatter抽出
- HTMLキャッシュ

### internal/search（全文検索）

- Bleveインデックス管理
- 検索クエリ実行

---

## 4. テスト構造

```
cmd/server/
    main_test.go          # サーバー初期化、ヘルスチェック
internal/
    model/
        *_test.go         # モデルバリデーション
    repository/
        *_test.go         # DB操作テスト
    service/
        *_test.go         # ビジネスロジックテスト
    handler/
        *_test.go         # HTTPハンドラーテスト
    markdown/
        *_test.go         # パーステスト
```

テストパターン:
- テーブルドリブンテスト必須
- カバレッジ目標: 80%以上
- `go test -race` でデータ競合検出

---

## 5. 設定・インフラ

```
.claude/                    # Claude Code設定
├── rules/                  # CODE_STYLE, CODING_RULES, GIT_WORKFLOW, etc.
└── CLAUDE.md               # プロジェクトルール

.moai/                      # MoAI-ADK設定
├── config/sections/        # user, language, quality, project, workflow, etc.
├── project/                # tech.md, structure.md, product.md
└── specs/                  # SPEC文書

.github/workflows/
└── ci.yml                  # コード品質 + テスト
```
