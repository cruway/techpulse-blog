# Gitワークフロー

> **関連文書**: [PR_CHECKLIST.md](./PR_CHECKLIST.md)

全ての作業は以下のルールに従います:

## 0. Statusワークフロー (Backlog → Todo → In progress → Done)

```
Backlog ──→ Todo ──→ In progress ──→ Done
(非活性Epic)  (作業待機)  (進行中)       (完了)
```

1. **Todo項目優先消化**: Todoにある全てのFeatを先に完了する
2. **In progress完了**: 現在進行中の全ての作業をDoneに処理する
3. **Backlogから昇格**: TodoとIn progressが全て空になってからBacklogから次のEpicを取得する
4. **Epic単位移動**: Backlog → In progressにEpic移動時、下位FeatをTodoに配置する

## 1. Issue作成

> **GitHub MCPツール優先**: `mcp__github__issue_write`, `mcp__github__issue_read` 使用

- **Issueタイトル**: `[タイプ] タイトル`（Issue番号は付けない）
- **PRタイトルと区別**: PRタイトルのみ `(#Issue番号)` を含む

| タイプ | 用途 |
|------|------|
| [Epic] | Phase単位 |
| [Feat] | 機能実装 |
| [Fix] | バグ修正 |
| [Chore] | CI、設定 |
| [Docs] | ドキュメント |
| [Refactor] | リファクタリング |
| [Test] | テスト |

## 2. ブランチ管理

- 命名: `feature/{Issue番号}-{簡単な説明}` (例: `feature/3-blog-engine-mvp`)
- main直接コミット禁止
- マージ後自動削除

## 3. コミット

```
type: 簡単な説明

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
```

type: feat, fix, docs, style, refactor, test, chore
- コミットメッセージは日本語で作成（typeは英語）

## 4. PR作成

> **MCPツール優先**: `mcp__github__create_pull_request`, `mcp__github__merge_pull_request`

- PRタイトル: `[タイプ] 簡単な説明 (#Issue番号)`
- Issue連結: `Closes #Issue番号`
- 必須: Assignees (cruway), Labels
- Mermaidダイアグラム必須

## 5. マージ

- CI通過必須
- **ユーザー承認必須** — Claudeは絶対に任意でマージしない
- Squash merge推奨

## 6. 完了処理

### 6-1. Feat完了（PRマージ直後 — Claude自動実行）

1. Feat Issueチェックボックス完了
2. 上位Epicチェックボックス更新
3. プロジェクトフィールド更新: Status → Done

### 6-2. Epic完了（全Feat完了時 — ユーザー承認後）

1. Epic完了条件チェックボックス更新
2. Epic Issueクローズ（ユーザー承認必須）

## CIチェック項目

- コードフォーマット (`gofmt`)
- 静的解析 (`golangci-lint`)
- テスト (`go test ./...`)
- templ生成 (`templ generate`)
- ビルド (`go build`)
