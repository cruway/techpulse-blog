# 코드 스타일 규칙

> **관련 문서**: [CODING_RULES.md](./CODING_RULES.md) | [PR_CHECKLIST.md](./PR_CHECKLIST.md)

---

## 1. 주석 언어

- **모든 코드 주석은 일본어로 작성**
- 함수/구조체/인터페이스 문서 주석 포함
- 인라인 주석 포함
- 예외: 외부 라이브러리 인터페이스 구현 시 영문 유지 가능

```go
// PostServiceはブログポストのCRUD操作を処理します。
//
// マークダウンパース、frontmatter抽出、検索インデキシングを含みます。
type PostService struct {
    // ポスト保存所
    repo PostRepository
    // 全文検索エンジン
    search *bleve.Index
}

// GetBySlugはスラッグでポストを照会します。
// ポストがない場合はErrPostNotFoundを返します。
func (s *PostService) GetBySlug(slug string) (*Post, error) {
    // キャッシュから先に照会
    if cached, ok := s.cache.Get(slug); ok {
        return cached, nil
    }
    return s.repo.FindBySlug(slug)
}
```

---

## 2. 네이밍 규칙

| 대상 | 규칙 | 예시 |
|------|------|------|
| 패키지 | lowercase | `handler`, `service`, `model` |
| 구조체/인터페이스 | PascalCase | `PostService`, `SearchEngine` |
| 공개 함수/메서드 | PascalCase | `GetBySlug`, `ParseMarkdown` |
| 비공개 함수/메서드 | camelCase | `parseContent`, `buildIndex` |
| 상수 | PascalCase 또는 camelCase | `MaxPageSize`, `defaultTimeout` |
| 파일 | snake_case | `post_handler.go`, `feed_service.go` |
| 테스트 파일 | snake_case + _test | `post_handler_test.go` |

---

## 3. Go 코드 포맷팅

- `gofmt` / `goimports` 필수 적용
- `golangci-lint` 정적 분석 통과

---

## 4. 품질 규칙

### 하드코딩 금지

```go
// 위반
if pageSize > 20 { ... }

// 준수
const MaxPageSize = 20
if pageSize > MaxPageSize { ... }
```

### 에러 처리 필수

```go
// 위반
result, _ := doSomething()

// 준수
result, err := doSomething()
if err != nil {
    return fmt.Errorf("작업 실패: %w", err)
}
```

### 공개 API 문서화

```go
// 위반 - 문서 없음
func ParseMarkdown(content []byte) (*Post, error) { ... }

// 준수
// ParseMarkdown은 마크다운 콘텐츠를 파싱하여 Post 구조체로 변환합니다.
//
// frontmatter(YAML)와 본문을 분리하여 처리합니다.
// Mermaid 코드 블록은 클라이언트 렌더링용으로 보존됩니다.
func ParseMarkdown(content []byte) (*Post, error) { ... }
```
