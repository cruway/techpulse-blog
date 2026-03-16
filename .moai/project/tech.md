# TechPulse Blog — Technology Document

> 生成日: 2026-03-16
> ソース基盤: コードベース自動分析

---

## 1. 技術スタック概要

| 階層 | 技術 | バージョン |
|------|------|-----------|
| **言語** | Go | 1.26.1 |
| **Webフレームワーク** | Echo | v4.15.1 |
| **テンプレート** | templ | （次Phase導入） |
| **フロントエンド** | HTMX + Tailwind CSS | （次Phase導入） |
| **DB** | SQLite (WALモード) | — |
| **検索** | Bleve | — |
| **マークダウン** | goldmark | — |
| **インフラ** | Oracle Cloud Always Free | — |
| **CI/CD** | GitHub Actions | ubuntu-latest |

---

## 2. フレームワーク選択根拠

### Echo (Webフレームワーク)

- **根拠**: 軽量・高速なGoのWebフレームワーク
- **利点**: ミドルウェアエコシステム、シンプルなAPI
- **用途**: HTTPハンドラー、静的ファイル配信、ミドルウェア管理

### templ (テンプレート)

- **根拠**: Goの型安全テンプレートエンジン
- **利点**: コンパイル時検証、Goとの自然な統合
- **用途**: HTMLレンダリング、コンポーネントベースUI

### HTMX (フロントエンド)

- **根拠**: JSフレームワーク不要のインタラクティブUI
- **利点**: サーバーサイドレンダリングとの相性、軽量
- **用途**: ページネーション、動的コンテンツロード

### SQLite (データベース)

- **根拠**: 組み込みDB、サーバー不要
- **利点**: Always Free環境に最適、WALモードで読み取り性能向上
- **用途**: ポスト、タグ、カテゴリの永続化

---

## 3. アーキテクチャ決定

### レイヤー分離（4層）

```
決定: Handler → Service → Repository → Model の4層分離
根拠: テスト容易性、責務の明確化
トレードオフ: 小規模プロジェクトにはオーバーヘッド
緩和: 初期はシンプルに保ち、成長に応じて拡張
```

### インターフェース基盤DI

```
決定: ServiceはRepositoryインターフェースに依存
根拠: テスト時にモック差し替え可能
トレードオフ: インターフェース定義のボイラープレート
緩和: 必要最小限のメソッドのみ定義
```

### マークダウンベースコンテンツ

```
決定: ファイルシステム上のマークダウンをDBに同期
根拠: Obsidianとの互換性、Git管理可能
トレードオフ: DB同期のタイミング管理
緩和: ファイル監視 + 起動時フルスキャン
```

---

## 4. 主要依存関係

### 現在

| パッケージ | 用途 |
|-----------|------|
| github.com/labstack/echo/v4 | Webフレームワーク |

### 予定（Phase 1完了まで）

| パッケージ | 用途 |
|-----------|------|
| github.com/a-h/templ | テンプレートエンジン |
| github.com/yuin/goldmark | マークダウンパーサー |
| github.com/blevesearch/bleve | 全文検索 |
| modernc.org/sqlite | SQLiteドライバ（CGO不要） |

---

## 5. 開発環境

### 必須ツール

| ツール | バージョン | 用途 |
|--------|-----------|------|
| Go | 1.26+ | 言語 |
| templ | latest | テンプレート生成 |
| air | latest | ホットリロード |
| golangci-lint | latest | 静的解析 |
| Git | 最新 | バージョン管理 |

### ビルドコマンド

```bash
# ビルド
make build

# 開発モード（ホットリロード）
make dev

# テスト
make test

# カバレッジ
make coverage

# リント
make lint

# フォーマット
make fmt

# templ生成
make templ
```

---

## 6. CI/CDパイプライン

### GitHub Actions (ci.yml)

| Job | 内容 |
|-----|------|
| コード品質 | `gofmt` + `golangci-lint` + `templ generate` + `go build` |
| テスト | `go test -v -race -coverprofile` + Codecovアップロード |

**トリガー**: main push + PR (*.md, content/** 除外)

---

## 7. 開発方法論

| 設定 | 値 | 説明 |
|------|-----|------|
| 方法論 | Hybrid | 新規コード: TDD、レガシー: DDD |
| カバレッジ目標 | 80% | 新規80%、レガシー70% |
| LSP品質ゲート | 有効 | エラー0、型エラー0、リントエラー0 |
| コミット規約 | Conventional | feat, fix, docs, style, refactor, test, chore |
| コードコメント | 日本語 | 全てのコメント日本語 |
| ドキュメント言語 | 日本語 | PR, Issue, アーキテクチャ文書 |
