# 코딩 규칙

> **관련 문서**: [CODE_STYLE.md](./CODE_STYLE.md) | [PR_CHECKLIST.md](./PR_CHECKLIST.md)

TechPulse Blog 프로젝트의 아키텍처, 패턴, 성능, 에러 처리 규칙을 정의합니다.

---

## 1. 아키텍처 규칙

### 1.1 레이어 분리 원칙

```
cmd/server/main.go          ← 엔트리포인트
    |
internal/handler/            ← HTTP 핸들러 (요청/응답)
    |
internal/service/            ← 비즈니스 로직
    |
internal/repository/         ← 데이터 접근 (SQLite, 파일시스템)
    |
internal/model/              ← 데이터 모델 (구조체)
```

| 계층 | 패키지 | 허용 의존성 |
|------|--------|-------------|
| Model | `internal/model` | 표준 라이브러리만 |
| Repository | `internal/repository` | model |
| Service | `internal/service` | model, repository |
| Handler | `internal/handler` | model, service |
| Markdown | `internal/markdown` | model, goldmark |
| Search | `internal/search` | model, bleve |

### 1.2 의존성 규칙

- **하위 → 상위 의존 금지**: model이 handler를 import하면 안 됨
- **인터페이스 기반 의존성 주입**: service는 repository 인터페이스에 의존
- **순환 의존 금지**: 패키지 간 순환 import 불가 (Go 컴파일러가 차단)

```go
// 준수 - 인터페이스 기반
type PostRepository interface {
    FindBySlug(slug string) (*model.Post, error)
    FindAll(opts ListOptions) ([]*model.Post, error)
}

type PostService struct {
    repo PostRepository  // 인터페이스 의존
}
```

---

## 2. 에러 처리 규칙

### 2.1 에러 래핑 필수

```go
// 위반
return err

// 준수
return fmt.Errorf("포스트 조회 실패 (slug=%s): %w", slug, err)
```

### 2.2 커스텀 에러 타입

```go
// 도메인 에러 정의
var (
    ErrPostNotFound = errors.New("포스트를 찾을 수 없습니다")
    ErrInvalidSlug  = errors.New("유효하지 않은 슬러그입니다")
)
```

### 2.3 HTTP 에러 응답

```go
// handler 레이어에서만 HTTP 상태 코드 결정
func (h *PostHandler) GetPost(c echo.Context) error {
    post, err := h.service.GetBySlug(slug)
    if err != nil {
        if errors.Is(err, service.ErrPostNotFound) {
            return c.Render(http.StatusNotFound, "404", nil)
        }
        return c.Render(http.StatusInternalServerError, "500", nil)
    }
    return c.Render(http.StatusOK, "post", post)
}
```

---

## 3. 성능 규칙

### 3.1 마크다운 파싱 캐시

- 파싱된 HTML은 메모리 캐시에 보관
- 파일 변경 감지 시 캐시 무효화
- 최대 캐시 크기 설정

### 3.2 데이터베이스

- SQLite WAL 모드 사용 (읽기 성능 향상)
- 인덱스: slug, date, tags, category
- 쿼리에 LIMIT/OFFSET 필수

### 3.3 HTTP 성능

- 정적 파일 캐시 헤더 설정
- gzip 압축 미들웨어 적용
- templ 컴포넌트 렌더링 (스트리밍)

---

## 4. 테스트 규칙

### 4.1 테스트 구조

```
internal/
    service/
        post_service.go
        post_service_test.go    ← 같은 패키지에 테스트
    handler/
        post_handler.go
        post_handler_test.go
```

### 4.2 테이블 드리븐 테스트

```go
func TestParseMarkdown(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    *Post
        wantErr bool
    }{
        {
            name:  "정상적인 포스트",
            input: "---\ntitle: test\n---\n# Hello",
            want:  &Post{Title: "test"},
        },
        {
            name:    "frontmatter 없음",
            input:   "# No frontmatter",
            wantErr: true,
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseMarkdown([]byte(tt.input))
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseMarkdown() error = %v, wantErr %v", err, tt.wantErr)
            }
            if !tt.wantErr && got.Title != tt.want.Title {
                t.Errorf("ParseMarkdown() title = %v, want %v", got.Title, tt.want.Title)
            }
        })
    }
}
```

---

## 5. Quick Reference

| 항목 | 준수 | 위반 |
|------|------|------|
| 의존성 방향 | handler → service → repo → model | model → handler |
| 에러 처리 | `fmt.Errorf("...: %w", err)` | `return err` (래핑 없이) |
| 에러 무시 | 명시적 `_ =` + 주석 | 암묵적 무시 |
| HTTP 상태 | handler에서만 결정 | service에서 HTTP 코드 반환 |
| SQL 쿼리 | Prepared statement | 문자열 연결 |
| 테스트 | 테이블 드리븐 | 단일 케이스 |

### 유용한 명령어

```bash
# 코드 포맷팅
gofmt -w .
goimports -w .

# 정적 분석
golangci-lint run

# 테스트 실행
go test ./...

# 테스트 커버리지
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# templ 코드 생성
templ generate

# 빌드
go build -o bin/server ./cmd/server
```
