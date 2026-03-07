# GitHub CLI コマンド集

> **関連文書**: [GIT_WORKFLOW.md](./GIT_WORKFLOW.md) | [PROJECT_FIELDS.md](./PROJECT_FIELDS.md)

## プロジェクトID情報

```bash
PROJECT_ID="PVT_kwHOAMfvP84BRFVI"
```

## GitHub MCP ツール代替

GitHub MCPサーバー接続時、以下の作業は `gh` CLI の代わりにMCPツールを優先使用します。

| 既存CLI | MCPツール | 備考 |
|---------|----------|------|
| `gh issue create` | `mcp__github__issue_write(method: "create")` | Issue作成 |
| `gh issue edit` | `mcp__github__issue_write(method: "update")` | Issue修正 |
| `gh issue view` | `mcp__github__issue_read(method: "get")` | Issue照会 |
| `gh pr create` | `mcp__github__create_pull_request` | PR作成 |
| `gh pr view` | `mcp__github__pull_request_read` | PR照会 |
| `gh pr merge` | `mcp__github__merge_pull_request` | PRマージ |
| `gh api graphql` (Sub-issue) | `mcp__github__sub_issue_write` | Sub-issue連結 |

**MCP未対応 — `gh` CLI継続使用：**
- `gh project item-add`（プロジェクトにIssue追加）
- `gh project item-edit`（プロジェクトフィールド編集）
- `gh api repos/.../milestones`（マイルストーン作成）
- `gh pr checks`（PR CI状態確認）

---

## Issueプロジェクト追加・フィールド設定

```bash
# 1. プロジェクトにIssue追加
gh project item-add 13 --owner cruway --url https://github.com/cruway/techpulse-blog/issues/{Issue番号}

# 2. プロジェクトアイテムID照会
ITEM_ID=$(gh api graphql -f query='
query {
  repository(owner: "cruway", name: "techpulse-blog") {
    issue(number: {Issue番号}) {
      projectItems(first: 1) {
        nodes { id }
      }
    }
  }
}' --jq '.data.repository.issue.projectItems.nodes[0].id')
```

## フィールドID・オプションID参照

### Status
```bash
# Backlog: f5b86059, Todo: ada40d47, In progress: 78c7b941, Done: f35112f5
gh project item-edit --project-id "$PROJECT_ID" --id "$ITEM_ID" \
  --field-id "PVTSSF_lAHOAMfvP84BRFVIzg_A_54" --single-select-option-id "78c7b941"
```

### 優先度
```bash
# P0: 79628723, P1: 0a877460, P2: da944a9c
gh project item-edit --project-id "$PROJECT_ID" --id "$ITEM_ID" \
  --field-id "PVTSSF_lAHOAMfvP84BRFVIzg_A_7A" --single-select-option-id "0a877460"
```

### サイズ
```bash
# XS: 911790be, S: b277fb01, M: 86db8eb3, L: 853c8207, XL: 2d0801e2
gh project item-edit --project-id "$PROJECT_ID" --id "$ITEM_ID" \
  --field-id "PVTSSF_lAHOAMfvP84BRFVIzg_A_7E" --single-select-option-id "86db8eb3"
```

### 作業種別
```bash
# 機能開発: 6bae7871, バグ修正: 12beaeae, リファクタリング: c8ce292a
# ドキュメント: 6ffecd88, テスト: 518e1f34, 設計: 5ce0abba
gh project item-edit --project-id "$PROJECT_ID" --id "$ITEM_ID" \
  --field-id "PVTSSF_lAHOAMfvP84BRFVIzg_BBC0" --single-select-option-id "6bae7871"
```

### 難易度
```bash
# 簡単: a197513f, 普通: 930890ea, 難しい: b3b43329, 要検討: 6ea2645d
gh project item-edit --project-id "$PROJECT_ID" --id "$ITEM_ID" \
  --field-id "PVTSSF_lAHOAMfvP84BRFVIzg_BBC4" --single-select-option-id "930890ea"
```

### 担当
```bash
# メイン: 1dd06f8b, サブ: 0b0aef56, 共同: 784d5685
gh project item-edit --project-id "$PROJECT_ID" --id "$ITEM_ID" \
  --field-id "PVTSSF_lAHOAMfvP84BRFVIzg_BBC8" --single-select-option-id "1dd06f8b"
```

### レビュー状態
```bash
# 未レビュー: a8753eb4, レビュー中: 1b91480e, 承認済み: 24907a2b, 修正必要: 38060299
gh project item-edit --project-id "$PROJECT_ID" --id "$ITEM_ID" \
  --field-id "PVTSSF_lAHOAMfvP84BRFVIzg_BBDA" --single-select-option-id "a8753eb4"
```

### ストーリーポイント
```bash
gh project item-edit --project-id "$PROJECT_ID" --id "$ITEM_ID" \
  --field-id "PVTF_lAHOAMfvP84BRFVIzg_A_7I" --number 3
```

### スプリント
```bash
# Iteration 1: 381c7c80, Iteration 2: 54cf5c95, Iteration 3: d2c335bc
# Iteration 4: b6a8f1bb, Iteration 5: 955c1297
gh project item-edit --project-id "$PROJECT_ID" --id "$ITEM_ID" \
  --field-id "PVTIF_lAHOAMfvP84BRFVIzg_A_7M" --iteration-id "381c7c80"
```

### 日付フィールド
```bash
# 開始日
gh project item-edit --project-id "$PROJECT_ID" --id "$ITEM_ID" \
  --field-id "PVTF_lAHOAMfvP84BRFVIzg_BBFk" --date "2026-03-08"

# 目標日
gh project item-edit --project-id "$PROJECT_ID" --id "$ITEM_ID" \
  --field-id "PVTF_lAHOAMfvP84BRFVIzg_BBFo" --date "2026-03-10"
```

### テキストフィールド
```bash
# スプリント目標
gh project item-edit --project-id "$PROJECT_ID" --id "$ITEM_ID" \
  --field-id "PVTF_lAHOAMfvP84BRFVIzg_BBGU" --text "目標内容"

# 完了定義
gh project item-edit --project-id "$PROJECT_ID" --id "$ITEM_ID" \
  --field-id "PVTF_lAHOAMfvP84BRFVIzg_BBGY" --text "テスト80%, CI通過"

# ブロッカー
gh project item-edit --project-id "$PROJECT_ID" --id "$ITEM_ID" \
  --field-id "PVTF_lAHOAMfvP84BRFVIzg_BBGc" --text "なし"
```

## Epic関連コマンド

### Epicストーリーポイント合計確認
```bash
EPIC_ID="I_xxxxx"  # EpicのNode ID
TOTAL_POINTS=$(gh api graphql -f query='
query($id: ID!) {
  node(id: $id) {
    ... on Issue {
      subIssues(first: 50) {
        nodes {
          projectItems(first: 1) {
            nodes {
              fieldValueByName(name: "ストーリーポイント") {
                ... on ProjectV2ItemFieldNumberValue {
                  number
                }
              }
            }
          }
        }
      }
    }
  }
}' -f id="$EPIC_ID" --jq '[.data.node.subIssues.nodes[].projectItems.nodes[].fieldValueByName.number // 0] | add')
echo "合計ストーリーポイント: $TOTAL_POINTS"
```

### Sub-issue関係設定
```bash
gh api graphql -f query='
mutation {
  addSubIssue(input: {
    issueId: "'$(gh api repos/cruway/techpulse-blog/issues/{Epic番号} --jq '.node_id')'",
    subIssueId: "'$(gh api repos/cruway/techpulse-blog/issues/{Feat番号} --jq '.node_id')'"
  }) {
    issue { number }
    subIssue { number }
  }
}'
```

## GitHub Project ワークフロー（自動化）

GitHub Project #13に設定されたビルトインワークフローです。
このワークフローはGitHub UIでのみ修正可能です（Settings → Workflows）。

| # | ワークフロー | 状態 | 動作 |
|---|-------------|------|------|
| 1 | Auto-add to project | 要設定 | Issue/PR作成時にプロジェクトに自動追加 |
| 2 | Auto-add sub-issues | 要設定 | Sub-issue作成時にプロジェクトに自動追加 |
| 3 | Item added to project | 要設定 | プロジェクト追加時 Status → **Todo** |
| 4 | Pull request linked | 要設定 | PRがIssueに連結時 Status → **In progress** |
| 5 | Item closed | 要設定 | Issue閉じ時 Status → **Done** |
| 6 | Pull request merged | 要設定 | PRマージ時 Status → **Done** |
| 7 | Auto-close issue | 要設定 | 連結PR マージ時Issue自動閉じ |

### Backlogワークフロー運用規則

- **Feat Issue**: 自動でTodo配置 → ワークフローデフォルト動作と一致
- **Epic Issue**: 自動Todo配置後 **手動でBacklog移動必要**
- EpicをBacklogに移動するCLI:
```bash
# EpicをBacklogに移動
gh project item-edit --project-id "$PROJECT_ID" --id "$ITEM_ID" \
  --field-id "PVTSSF_lAHOAMfvP84BRFVIzg_A_54" --single-select-option-id "f5b86059"
```

## その他の便利なコマンド

### Assignee設定
```bash
gh issue edit {Issue番号} --repo cruway/techpulse-blog --add-assignee $(gh api user --jq '.login')
```

### マイルストーン設定
```bash
gh issue edit {Issue番号} --repo cruway/techpulse-blog --milestone "v0.1.0 - ブログエンジン"
```

### PR作成（ラベル含む）
```bash
gh pr create --label "phase-1-blog-engine,type-feature,size-s" --title "[Feat] タイトル (#Issue番号)" --body "..."
```

### Project Status Update（Epic基盤）

Epic追加/完了時にプロジェクト状態を更新します。

```bash
# Status Update作成 (ON_TRACK, AT_RISK, OFF_TRACK, COMPLETE, INACTIVE)
gh api graphql -f query='mutation {
  createProjectV2StatusUpdate(input: {
    projectId: "PVT_kwHOAMfvP84BRFVI"
    body: "## Phase X: タイトル\n\n### 完了\n- ...\n\n### 進行中\n- ..."
    status: ON_TRACK
  }) { statusUpdate { id status } }
}'
```

```bash
# Project README更新
gh api graphql -f query='mutation {
  updateProjectV2(input: {
    projectId: "PVT_kwHOAMfvP84BRFVI"
    readme: "README内容"
    shortDescription: "短い説明"
  }) { projectV2 { id } }
}'
```
