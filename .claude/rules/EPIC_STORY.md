# Epic-Feat 作業分割戦略

> **関連文書**: [PROJECT_FIELDS.md](./PROJECT_FIELDS.md) | [CLI_COMMANDS.md](./CLI_COMMANDS.md) | [LABELS.md](./LABELS.md)

## Issue階層構造

全ての大規模作業（Epic）はFeat Issueに分割して進行します：

```
Epic Issue（Phase全体）
|-- Feat Issue #1 → PR #1（Closes #1）
|-- Feat Issue #2 → PR #2（Closes #2）
|-- Feat Issue #3 → PR #3（Closes #3）
+-- ...
```

## Feat Issue作成規則

### 1. 分割基準

- **サイズ**: 半日〜2日の作業量（小さすぎず大きすぎず）
- **単一責任**: 一つの機能・目的のみ担当
- **独立テスト**: Feat単位でテスト可能であること
- **1:1マッピング**: 1 Feat Issue = 1 PR

### 2. Wave構造

依存関係に従いWaveでグループ化：

```
Wave 1: 基盤（依存なし）- 並列作業可能
Wave 2: 拡張（Wave 1依存）
Wave 3: 核心（Wave 1, 2依存）
...
```

### 3. Feat Issue必須セクション

```markdown
## 概要
（何を実装するか簡単に）

## 目標
（達成すべきこと）

## 実装内容
（具体的な実装項目 - チェックボックス）

## 完了条件
（PRマージ前に満たすべき条件）

## 関連
- 上位Epic: #番号
- 依存: #番号（ある場合）

## 予想作業量
- サイズ: S/M/L
- 予想: 半日/1日/1.5日
```

## Epic Issue管理

### Epic Issue更新

Feat Issue作成後、Epic Issue本文に追加：
- 全Feat Issue一覧（チェックボックス）
- 依存関係グラフ（Mermaid）
- Wave別グループ化

### Epicストーリーポイント

**EpicのSP = 全Sub-issue（Feat）のSP合計**

GitHub Projects v2は自動合計をサポートしないため手動管理：

1. **Feat作成時**: 各FeatにSP設定
2. **Epic更新**: Feat全体合計をEpicのSPに入力
3. **変更時更新**: FeatのSP変更時Epicも更新

| Epic | Feat 1 | Feat 2 | Feat 3 | 合計 |
|------|--------|--------|--------|------|
| #1 | 3点 | 5点 | 2点 | **10点** |

### Epicスクラムフィールド

| フィールド | Epic設定基準 | 例 |
|------------|-------------|-----|
| ストーリーポイント | 全Feat合計 | 10 |
| スプリント | 最初のFeatが属するスプリント | 最初のスプリント選択 |
| スプリント目標 | Epic全体目標 | "Phase 1 ブログエンジン完成" |
| 完了定義 | Epic完了条件 | "全Feat完了 + 統合テスト" |
| ブロッカー | Epicレベルのブロッカー | "外部API連携検討必要" |

## Feat Issue設定

### Sub-issue連結（MCPツール優先）

GitHub MCPサーバー接続時、Sub-issue連結はMCPツールを優先使用します：
- `mcp__github__sub_issue_write(method: "add")` — Sub-issue追加
- `mcp__github__sub_issue_write(method: "remove")` — Sub-issue削除
- `mcp__github__sub_issue_write(method: "reprioritize")` — 順序変更

### 必須設定項目

| 設定 | 内容 |
|------|------|
| Assignee | GitHubログインユーザー |
| Milestone | Phase別マイルストーン |
| Labels | phase-X-xxx, size-s/m/l, type-feature |
| Parent Issue | 上位Epic（Sub-issueとして連結） |
| プロジェクトフィールド | [PROJECT_FIELDS.md](./PROJECT_FIELDS.md) 参照 |

### マイルストーンマッピング

| Phase | マイルストーン |
|-------|---------------|
| Phase 1 | v0.1.0 - ブログエンジン |
| Phase 2 | v0.2.0 - Obsidian連携 |
| Phase 3 | v0.3.0 - RSS収集 |
| Phase 4 | v0.4.0 - 技術エクスプローラー |
| Phase 5 | v0.5.0 - デプロイ |
| Phase 6 | v1.0.0 - ドメイン・正式公開 |

### 作業順序

1. Epic Issue分析
2. Feat Issue作成（Wave順序で）
3. 各FeatにAssignee設定
4. 各Featにマイルストーン設定
5. 各FeatをEpicのSub-issueとして連結
6. 各Featにプロジェクトフィールド設定
7. Epic Issue本文にFeat一覧追加
8. Wave 1から順次実装
9. 各Featは独立ブランチ + PRで進行

## サイズ基準

| サイズ | 作業量 | 例 |
|--------|--------|-----|
| XS | < 2時間 | 設定変更、バグ修正 |
| S | 半日 | 単一モデル/ユーティリティ実装 |
| M | 1日 | サービスクラス実装 |
| L | 1.5〜2日 | 統合システム実装 |
| XL | > 3日 | Epic（分割必要） |

**XL以上は必ずFeatに分割**
