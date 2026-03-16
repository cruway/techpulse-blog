# TechPulse Blog - Development Specification

## 1. Project Overview

**プロジェクト名**: TechPulse Blog
**目的**: IT技術トレンドを収集し、対話形式で整理してブログとして公開するシステム
**核心価値**: コスト最小化（$0運営）+ キャリア考察 + 技術探索

---

## 2. System Architecture

```
n8n (RSS収集, $0)
    |  Markdown保存
Obsidian Vault (= Git repo)
    |-- inbox/        <- n8nが毎日保存するフィード
    |-- drafts/       <- Claude Codeで整理した下書き
    |-- posts/        <- 公開確定した記事
    |-- knowledge/    <- キャリア/技術メモ
    |-- config/
    |   +-- feeds.yaml  <- 関心技術キーワード設定
    |
    |  git push (手動)
Go Blog Server (Oracle Cloud, $0)
    |-- ブログレンダリング (templ + HTMX)
    |-- 記事検索 (Bleve)
    +-- 技術探索UI (収集データ基盤)
```

---

## 3. Tech Stack

| Layer | Technology | Rationale |
|-------|-----------|-----------|
| Backend | Go (Echo) | 軽量、高性能、単一バイナリデプロイ |
| Template | templ | タイプセーフGoテンプレート、コンパイルタイムチェック |
| Frontend | HTMX + Tailwind CSS | JS最小化、Go SSRと相性良好 |
| Markdown Parser | goldmark | Goネイティブ、拡張可能（Mermaid、frontmatter） |
| Full-text Search | Bleve | Goネイティブ全文検索エンジン |
| Database | SQLite | 軽量、サーバーレス、バックアップ容易 |
| Automation | n8n (self-hosted) | RSS収集自動化、Docker基盤 |
| Editor | Obsidian | マークダウン編集、知識グラフ、Git連携 |
| AI Processing | Claude Code (既存サブスクリプション) | 手動対話形式コンテンツ整理 |
| Infra | Oracle Cloud (Always Free) | ARM VM 4コア/24GB無料 |
| Reverse Proxy | Caddy | 自動HTTPS、設定簡単 |
| Domain | Cloudflare + .dev | CDN + DNS + 強制HTTPS |

---

## 4. Core Features

### 4.1 Blog Engine (Phase 1)

- [ ] マークダウンファイル基盤ブログレンダリング
- [ ] frontmatterパース (title, date, tags, category, status)
- [ ] Mermaidダイアグラムクライアントレンダリング
- [ ] コードハイライト (highlight.js または Prism)
- [ ] レスポンシブレイアウト (Tailwind CSS)
- [ ] タグ/カテゴリ別リスト
- [ ] ページネーション
- [ ] RSSフィード生成

### 4.2 Obsidian Integration (Phase 2)

- [ ] Obsidian vaultフォルダ構造設計
- [ ] frontmatter規格定義
- [ ] posts/ フォルダ → ブログ自動反映 (git pushトリガー)
- [ ] Obsidian Gitプラグイン設定ガイド

### 4.3 RSS Collection Pipeline (Phase 3)

- [ ] n8nワークフロー: スケジュール → RSS収集 → キーワードフィルタリング
- [ ] feeds.yaml基盤ソース/キーワード管理
- [ ] inbox/ フォルダに日別フィードマークダウン生成
- [ ] 収集ソース: Hacker News, dev.to, GitHub Trending, Zenn等

### 4.4 Tech Explorer (Phase 4)

- [ ] Bleve基盤全文検索インデキシング
- [ ] 技術探索UI（検索 + フィルター）
- [ ] 収集した記事から関連リンク提示
- [ ] 自分のブログ関連ポスト連結
- [ ] タグ基盤技術カテゴリ探索

### 4.5 Deployment (Phase 5)

- [ ] Dockerfile (Goサーバー + n8n)
- [ ] docker-compose.yml
- [ ] Caddy設定（リバースプロキシ + HTTPS）
- [ ] Oracle Cloud ARM VMデプロイ
- [ ] CI/CD (GitHub Actions → Oracle Cloud)

### 4.6 Domain & DNS (Phase 6)

- [ ] .devドメイン購入 (Cloudflare)
- [ ] DNS設定
- [ ] Cloudflare CDN有効化

---

## 5. Data Models

### 5.1 Post（ブログ記事）

```yaml
# frontmatter
title: "Go 1.24の新機能まとめ"
date: 2026-03-08
tags: [go, release, language]
category: backend
status: published  # draft | review | published
source_url: "https://example.com/original"
mermaid: true
---
# 本文 (Markdown)
```

### 5.2 Feed Item（収集記事）

```yaml
# inbox/2026-03-08-feed.md frontmatter
date: 2026-03-08
keywords: [go, ai, cloud]
total_items: 12
filtered_items: 5
---
# 記事リスト
```

### 5.3 feeds.yaml（収集設定）

```yaml
sources:
  - name: "Hacker News"
    type: api
    url: "https://hacker-news.firebaseio.com/v0"
    min_score: 100

  - name: "dev.to"
    type: rss
    url: "https://dev.to/feed"

  - name: "GitHub Trending"
    type: scrape
    languages: [go, typescript, rust]

keywords:
  include:
    - go
    - golang
    - ai
    - llm
    - cloud
    - webassembly
    - kubernetes
  exclude:
    - crypto
    - nft

schedule: "0 9 * * *"  # 毎日午前9時
```

---

## 6. API Endpoints

```
GET  /                      # メイン（最新記事リスト）
GET  /posts                 # 記事リスト（ページネーション）
GET  /posts/:slug           # 記事詳細
GET  /tags                  # タグリスト
GET  /tags/:tag             # タグ別記事リスト
GET  /categories/:category  # カテゴリ別記事リスト
GET  /search                # 検索 (Bleve)
GET  /explore               # 技術探索UI
GET  /feed.xml              # RSSフィード
GET  /api/search            # 検索API (HTMX)
GET  /api/explore           # 探索API (HTMX)
```

---

## 7. Directory Structure (Go Project)

```
techpulse-blog/
|-- cmd/
|   +-- server/
|       +-- main.go
|-- internal/
|   |-- handler/        # HTTPハンドラー
|   |-- service/        # ビジネスロジック
|   |-- repository/     # データアクセス
|   |-- model/          # データモデル
|   |-- markdown/       # マークダウンパース
|   +-- search/         # Bleve検索
|-- templates/           # templファイル
|-- static/              # CSS, JS, 画像
|-- content/             # Obsidian vault (git submodule)
|   |-- inbox/
|   |-- drafts/
|   |-- posts/
|   |-- knowledge/
|   +-- config/
|-- deploy/
|   |-- Dockerfile
|   |-- docker-compose.yml
|   +-- Caddyfile
|-- go.mod
|-- go.sum
+-- README.md
```

---

## 8. Development Phases & Milestones

### Phase 1: Blog Engine (MVP)
**Goal**: マークダウンファイルを読み込んでブログとしてレンダリング

- Goプロジェクト初期化
- マークダウンパース + frontmatter処理
- templ基盤レイアウト（メイン、記事リスト、記事詳細）
- HTMX基盤ページネーション
- Tailwind CSSスタイリング
- Mermaid.jsクライアントレンダリング
- コードハイライト
- ローカル開発環境（hot reload）

### Phase 2: Obsidian Integration
**Goal**: Obsidian vaultとブログソース連携

- vaultフォルダ構造確定
- frontmatter規格確定
- Git submoduleまたは単一repo決定
- Obsidianプラグイン推薦/設定

### Phase 3: RSS Collection
**Goal**: n8nで技術記事自動収集

- n8n Docker設定
- RSS収集ワークフロー
- feeds.yamlパーサー
- キーワードフィルタリングロジック
- inbox/ マークダウン生成

### Phase 4: Tech Explorer
**Goal**: 収集した記事の検索/探索

- Bleveインデキシング (posts + inbox)
- 検索UI (HTMX)
- タグ/カテゴリフィルター
- 関連ポスト推薦

### Phase 5: Deployment
**Goal**: Oracle Cloudデプロイ

- Dockerイメージビルド
- docker-compose (Go + n8n + Caddy)
- Oracle Cloud VMプロビジョニング
- GitHub Actions CI/CD

### Phase 6: Domain
**Goal**: カスタムドメイン適用

- .devドメイン購入
- Cloudflare設定
- HTTPS確認

---

## 9. Non-functional Requirements

- **Performance**: ページロード < 200ms (SSR)
- **Cost**: 月$0運営（ドメイン除く）
- **Security**: HTTPS必須、XSS防止、CSPヘッダー
- **SEO**: メタタグ、OGP、sitemap.xml
- **Accessibility**: セマンティックHTML、キーボードナビゲーション
- **Backup**: Git基盤コンテンツバックアップ、SQLite定期バックアップ

---

## 10. Open Questions

- [ ] Obsidian vaultを別repoで管理するか、ブログrepo内submoduleにするか
- [ ] ダークモード対応の可否
- [ ] コメント機能の必要性（giscus等）
- [ ] 多言語対応（日/英）の必要性
- [ ] Analytics（Umami等self-hosted）導入の可否
