# TechPulse Blog 프로젝트 규칙

## 핵심 원칙

1. **이슈 기반 작업**: 모든 작업은 이슈 생성 후 시작
2. **PR 머지 승인 필수**: Claude는 사용자의 명시적 승인 없이 PR을 머지하지 않음
3. **문서화 필수**: 모든 PR에 Mermaid 다이어그램 및 변경 이력 문서 포함
4. **일본어 주석**: 모든 코드 주석, 커밋 메시지, 이슈/PR 내용은 일본어로 작성

## Claude 자동 작업 제한 (절대 규칙)

| 작업 | 승인 필요 |
|------|----------|
| PR 머지 | 필수 |
| 브랜치 삭제 | 필수 |
| main 브랜치 push | 필수 |
| 이슈 닫기 | 필수 |
| force push | 필수 |

## 프로젝트 정보

- **저장소**: cruway/techpulse-blog
- **언어**: Go (Echo + templ + HTMX)
- **프론트엔드**: HTMX + Tailwind CSS
- **DB**: SQLite
- **검색**: Bleve
- **인프라**: Oracle Cloud (Always Free)

## Status 워크플로우

```
Backlog → Todo → In progress → Done
```
- Todo 우선 소진 → In progress 완료 → Backlog에서 다음 Epic 승격
- Epic: Backlog(대기) → In progress(활성) → Done
- Feat: Todo(대기) → In progress(구현중) → Done

## 이슈 타입 / 크기 기준

| 타입 | 용도 | | 크기 | 작업량 |
|------|------|-|------|--------|
| [Epic] | Phase 단위 | | XS | < 2시간 |
| [Feat] | 기능 구현 | | S | 반나절 |
| [Fix] | 버그 수정 | | M | 1일 |
| [Chore] | CI, 설정 | | L | 1.5~2일 |
| [Docs] | 문서 | | XL | > 3일 (분할) |
| [Refactor] | 리팩토링 | | | |
| [Test] | 테스트 | | | |

## GitHub MCP 도구 우선 사용

GitHub MCP 서버가 연결되어 있으므로, 이슈/PR 관련 작업은 **MCP 도구를 우선 사용**합니다.

| 작업 | MCP 도구 | 비고 |
|------|----------|------|
| 이슈 조회 | `mcp__github__issue_read` | `gh issue view` 대체 |
| 이슈 생성/수정 | `mcp__github__issue_write` | `gh issue create/edit` 대체 |
| 이슈 검색 | `mcp__github__search_issues` | `gh search issues` 대체 |
| Sub-issue 연결 | `mcp__github__sub_issue_write` | `gh api graphql` 뮤테이션 대체 |
| PR 생성 | `mcp__github__create_pull_request` | `gh pr create` 대체 |
| PR 조회 | `mcp__github__pull_request_read` | `gh pr view` 대체 |
| PR 머지 | `mcp__github__merge_pull_request` | `gh pr merge` 대체 |
| PR 수정 | `mcp__github__update_pull_request` | `gh pr edit` 대체 |
| 라벨 조회 | `mcp__github__get_label` | `gh label list` 대체 |
| 코드 검색 | `mcp__github__search_code` | `gh search code` 대체 |

**`gh` CLI를 계속 사용해야 하는 영역** (MCP 미지원):
- GitHub Projects v2 필드 편집 (`gh project item-edit`)
- 프로젝트에 이슈 추가 (`gh project item-add`)
- 마일스톤 생성 (`gh api repos/.../milestones`)
- PR checks 확인 (`gh pr checks`)

## 참조 문서 (필요시 Read로 로딩)

| 분류 | 파일 | 자동 로딩 |
|------|------|----------|
| **코드 스타일** | .claude/rules/CODE_STYLE.md | 항상 |
| **코딩 규칙** | .claude/rules/CODING_RULES.md | Go 파일 작업시 |
| **Git 워크플로우** | .claude/rules/GIT_WORKFLOW.md | 항상 |
| **PR 체크리스트** | .claude/rules/PR_CHECKLIST.md | 이슈/PR 작업시 |
| **프로젝트 필드** | .claude/rules/PROJECT_FIELDS.md | 이슈/PR 작업시 |
| **작업 분할** | .claude/rules/EPIC_STORY.md | 이슈/PR 작업시 |
| **CLI 명령어** | .claude/rules/CLI_COMMANDS.md | 이슈/PR 작업시 |
| **라벨** | .claude/rules/LABELS.md | 이슈/PR 작업시 |
