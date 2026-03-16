---
title: "Hello, TechPulse Blog!"
date: 2026-03-16
tags: [go, echo, htmx]
category: "お知らせ"
status: published
source_url: ""
mermaid: true
---

# TechPulse Blogへようこそ

これはTechPulse Blogの最初の記事です。Go + Echo + HTMX + Tailwind CSSで構築された技術ブログです。

## 技術スタック

| 技術 | 用途 |
|------|------|
| Go | バックエンド |
| Echo | Webフレームワーク |
| templ | テンプレート |
| HTMX | フロントエンド |
| Tailwind CSS | スタイリング |
| SQLite | データベース |

## アーキテクチャ

```mermaid
graph TD
    A[マークダウン] --> B[goldmark パーサー]
    B --> C[HTML レンダリング]
    C --> D[templ テンプレート]
    D --> E[ブラウザ]
```

## コード例

```go
func main() {
    fmt.Println("Hello, TechPulse!")
}
```

今後もIT技術トレンドの情報を発信していきます。
