# TechPulse Blog プロジェクトルール

## 核心原則

1. **Issue基盤作業**: 全ての作業はIssue作成後に開始
2. **PRマージ承認必須**: Claudeはユーザーの明示的承認なしにPRをマージしない
3. **ドキュメント化必須**: 全てのPRにMermaidダイアグラム及び変更履歴ドキュメントを含む
4. **日本語コメント**: 全てのコードコメント、コミットメッセージ、Issue/PR内容は日本語で作成

## ローカルCI必須（絶対ルール）

コミット・プッシュ前に必ず `make ci` を実行し、全チェック通過を確認すること。

| タイミング | コマンド | 内容 |
|-----------|---------|------|
| コミット前 | `make ci-quick` | gofmt + go vet + テスト |
| プッシュ前 | `make ci` | fmt + lint + build + test + カバレッジ閾値 |

**CI失敗時はコミット・プッシュを行わない。** 問題を修正してから再実行する。

## Claude自動作業制限（絶対ルール）

| 作業 | 承認必要 |
|------|----------|
| PRマージ | 必須 |
| ブランチ削除 | 必須 |
| mainブランチpush | 必須 |
| Issueクローズ | 必須 |
| force push | 必須 |

## プロジェクト情報

- **リポジトリ**: cruway/techpulse-blog
- **言語**: Go (Echo + templ + HTMX)
- **フロントエンド**: HTMX + Tailwind CSS
- **DB**: SQLite
- **検索**: Bleve
- **インフラ**: Oracle Cloud (Always Free)

## Statusワークフロー

```
Backlog → Todo → In progress → Done
```
- Todo優先消化 → In progress完了 → Backlogから次のEpic昇格
- Epic: Backlog(待機) → In progress(活性) → Done
- Feat: Todo(待機) → In progress(実装中) → Done

## Issueタイプ / サイズ基準

| タイプ | 用途 | | サイズ | 作業量 |
|------|------|-|------|--------|
| [Epic] | Phase単位 | | XS | < 2時間 |
| [Feat] | 機能実装 | | S | 半日 |
| [Fix] | バグ修正 | | M | 1日 |
| [Chore] | CI、設定 | | L | 1.5〜2日 |
| [Docs] | ドキュメント | | XL | > 3日（分割） |
| [Refactor] | リファクタリング | | | |
| [Test] | テスト | | | |

## GitHub MCPツール優先使用

GitHub MCPサーバーが接続されているため、Issue/PR関連作業は**MCPツールを優先使用**します。

| 作業 | MCPツール | 備考 |
|------|----------|------|
| Issue照会 | `mcp__github__issue_read` | `gh issue view` 代替 |
| Issue作成/修正 | `mcp__github__issue_write` | `gh issue create/edit` 代替 |
| Issue検索 | `mcp__github__search_issues` | `gh search issues` 代替 |
| Sub-issue連結 | `mcp__github__sub_issue_write` | `gh api graphql` ミューテーション代替 |
| PR作成 | `mcp__github__create_pull_request` | `gh pr create` 代替 |
| PR照会 | `mcp__github__pull_request_read` | `gh pr view` 代替 |
| PRマージ | `mcp__github__merge_pull_request` | `gh pr merge` 代替 |
| PR修正 | `mcp__github__update_pull_request` | `gh pr edit` 代替 |
| ラベル照会 | `mcp__github__get_label` | `gh label list` 代替 |
| コード検索 | `mcp__github__search_code` | `gh search code` 代替 |

**`gh` CLIを継続使用する領域**（MCP未対応）:
- GitHub Projects v2フィールド編集 (`gh project item-edit`)
- プロジェクトにIssue追加 (`gh project item-add`)
- マイルストーン作成 (`gh api repos/.../milestones`)
- PR checks確認 (`gh pr checks`)

## 参照ドキュメント（必要時Readでロード）

| 分類 | ファイル | 自動ロード |
|------|------|----------|
| **コードスタイル** | .claude/rules/CODE_STYLE.md | 常時 |
| **コーディングルール** | .claude/rules/CODING_RULES.md | Goファイル作業時 |
| **Gitワークフロー** | .claude/rules/GIT_WORKFLOW.md | 常時 |
| **PRチェックリスト** | .claude/rules/PR_CHECKLIST.md | Issue/PR作業時 |
| **プロジェクトフィールド** | .claude/rules/PROJECT_FIELDS.md | Issue/PR作業時 |
| **作業分割** | .claude/rules/EPIC_STORY.md | Issue/PR作業時 |
| **CLIコマンド** | .claude/rules/CLI_COMMANDS.md | Issue/PR作業時 |
| **ラベル** | .claude/rules/LABELS.md | Issue/PR作業時 |
