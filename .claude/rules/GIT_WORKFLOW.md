# Git 워크플로우

> **관련 문서**: [PR_CHECKLIST.md](./PR_CHECKLIST.md)

모든 작업은 다음 규칙을 따릅니다:

## 0. Status 워크플로우 (Backlog → Todo → In progress → Done)

```
Backlog ──→ Todo ──→ In progress ──→ Done
(비활성 Epic)  (작업 대기)  (진행중)       (완료)
```

1. **Todo 항목 우선 소진**: Todo에 있는 모든 Feat를 먼저 완료한다
2. **In progress 완료**: 현재 진행중인 모든 작업을 Done으로 처리한다
3. **Backlog에서 승격**: Todo와 In progress가 모두 비어야 Backlog에서 다음 Epic을 가져온다
4. **Epic 단위 이동**: Backlog → In progress로 Epic 이동 시, 하위 Feat들을 Todo에 배치한다

## 1. 이슈 생성

> **GitHub MCP 도구 우선**: `mcp__github__issue_write`, `mcp__github__issue_read` 사용

- **이슈 제목**: `[타입] 제목` (이슈 번호 붙이지 않음)
- **PR 제목과 구분**: PR 제목만 `(#이슈번호)` 포함

| 타입 | 용도 |
|------|------|
| [Epic] | Phase 단위 |
| [Feat] | 기능 구현 |
| [Fix] | 버그 수정 |
| [Chore] | CI, 설정 |
| [Docs] | 문서 |
| [Refactor] | 리팩토링 |
| [Test] | 테스트 |

## 2. 브랜치 관리

- 명명: `feature/{이슈번호}-{간단한설명}` (예: `feature/3-blog-engine-mvp`)
- main 직접 커밋 금지
- 머지 후 자동 삭제

## 3. 커밋

```
type: 簡単な説明

Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
```

type: feat, fix, docs, style, refactor, test, chore
- 커밋 메시지는 일본어로 작성 (type은 영어)

## 4. PR 생성

> **MCP 도구 우선**: `mcp__github__create_pull_request`, `mcp__github__merge_pull_request`

- PR 제목: `[타입] 간단한 설명 (#이슈번호)`
- 이슈 연결: `Closes #이슈번호`
- 필수: Assignees (cruway), Labels
- Mermaid 다이어그램 필수

## 5. 머지

- CI 통과 필수
- **사용자 승인 필수** — Claude는 절대 임의로 머지하지 않음
- Squash merge 권장

## 6. 완료 처리

### 6-1. Feat 완료 (PR 머지 직후 — Claude 자동 실행)

1. Feat 이슈 체크박스 완료
2. 상위 Epic 체크박스 업데이트
3. 프로젝트 필드 갱신: Status → Done

### 6-2. Epic 완료 (모든 Feat 완료 시 — 사용자 승인 후)

1. Epic 완료 조건 체크박스 업데이트
2. Epic 이슈 닫기 (사용자 승인 필수)

## CI 체크 항목

- 코드 포맷팅 (`gofmt`)
- 정적 분석 (`golangci-lint`)
- 테스트 (`go test ./...`)
- templ 생성 (`templ generate`)
- 빌드 (`go build`)
