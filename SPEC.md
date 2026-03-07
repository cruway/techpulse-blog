# TechPulse Blog - Development Specification

## 1. Project Overview

**프로젝트명**: TechPulse Blog
**목적**: IT 기술 트렌드를 수집하고, 대화형으로 정리하여 블로그로 게시하는 시스템
**핵심 가치**: 비용 최소화 ($0 운영) + 커리어 고찰 + 기술 탐색

---

## 2. System Architecture

```
n8n (RSS 수집, $0)
    |  Markdown 저장
Obsidian Vault (= Git repo)
    |-- inbox/        <- n8n이 매일 저장하는 피드
    |-- drafts/       <- Claude Code로 정리한 초안
    |-- posts/        <- 게시 확정된 글
    |-- knowledge/    <- 커리어/기술 메모
    |-- config/
    |   +-- feeds.yaml  <- 관심 기술 키워드 설정
    |
    |  git push (수동)
Go Blog Server (Oracle Cloud, $0)
    |-- 블로그 렌더링 (templ + HTMX)
    |-- 기사 검색 (Bleve)
    +-- 기술 탐색 UI (수집 데이터 기반)
```

---

## 3. Tech Stack

| Layer | Technology | Rationale |
|-------|-----------|-----------|
| Backend | Go (Echo) | 경량, 고성능, 단일 바이너리 배포 |
| Template | templ | 타입 세이프 Go 템플릿, 컴파일 타임 체크 |
| Frontend | HTMX + Tailwind CSS | JS 최소화, Go SSR과 궁합 |
| Markdown Parser | goldmark | Go 네이티브, 확장 가능 (Mermaid, frontmatter) |
| Full-text Search | Bleve | Go 네이티브 전문검색 엔진 |
| Database | SQLite | 경량, 서버리스, 백업 용이 |
| Automation | n8n (self-hosted) | RSS 수집 자동화, Docker 기반 |
| Editor | Obsidian | 마크다운 편집, 지식 그래프, Git 연동 |
| AI Processing | Claude Code (기존 구독) | 수동 대화형 콘텐츠 정리 |
| Infra | Oracle Cloud (Always Free) | ARM VM 4코어/24GB 무료 |
| Reverse Proxy | Caddy | 자동 HTTPS, 설정 간단 |
| Domain | Cloudflare + .dev | CDN + DNS + 강제 HTTPS |

---

## 4. Core Features

### 4.1 Blog Engine (Phase 1)

- [ ] 마크다운 파일 기반 블로그 렌더링
- [ ] frontmatter 파싱 (title, date, tags, category, status)
- [ ] Mermaid 다이어그램 클라이언트 렌더링
- [ ] 코드 하이라이팅 (highlight.js 또는 Prism)
- [ ] 반응형 레이아웃 (Tailwind CSS)
- [ ] 태그/카테고리별 목록
- [ ] 페이지네이션
- [ ] RSS 피드 생성

### 4.2 Obsidian Integration (Phase 2)

- [ ] Obsidian vault 폴더 구조 설계
- [ ] frontmatter 규격 정의
- [ ] posts/ 폴더 → 블로그 자동 반영 (git push 트리거)
- [ ] Obsidian Git 플러그인 설정 가이드

### 4.3 RSS Collection Pipeline (Phase 3)

- [ ] n8n 워크플로우: 스케줄 → RSS 수집 → 키워드 필터링
- [ ] feeds.yaml 기반 소스/키워드 관리
- [ ] inbox/ 폴더에 일별 피드 마크다운 생성
- [ ] 수집 소스: Hacker News, dev.to, GitHub Trending, Zenn 등

### 4.4 Tech Explorer (Phase 4)

- [ ] Bleve 기반 전문검색 인덱싱
- [ ] 기술 탐색 UI (검색 + 필터)
- [ ] 수집된 기사에서 관련 링크 제시
- [ ] 내 블로그 관련 포스트 연결
- [ ] 태그 기반 기술 카테고리 탐색

### 4.5 Deployment (Phase 5)

- [ ] Dockerfile (Go 서버 + n8n)
- [ ] docker-compose.yml
- [ ] Caddy 설정 (리버스 프록시 + HTTPS)
- [ ] Oracle Cloud ARM VM 배포
- [ ] CI/CD (GitHub Actions → Oracle Cloud)

### 4.6 Domain & DNS (Phase 6)

- [ ] .dev 도메인 구매 (Cloudflare)
- [ ] DNS 설정
- [ ] Cloudflare CDN 활성화

---

## 5. Data Models

### 5.1 Post (블로그 글)

```yaml
# frontmatter
title: "Go 1.24의 새로운 기능 정리"
date: 2026-03-08
tags: [go, release, language]
category: backend
status: published  # draft | review | published
source_url: "https://example.com/original"
mermaid: true
---
# 본문 (Markdown)
```

### 5.2 Feed Item (수집 기사)

```yaml
# inbox/2026-03-08-feed.md frontmatter
date: 2026-03-08
keywords: [go, ai, cloud]
total_items: 12
filtered_items: 5
---
# 기사 목록
```

### 5.3 feeds.yaml (수집 설정)

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

schedule: "0 9 * * *"  # 매일 오전 9시
```

---

## 6. API Endpoints

```
GET  /                      # 메인 (최신 글 목록)
GET  /posts                 # 글 목록 (페이지네이션)
GET  /posts/:slug           # 글 상세
GET  /tags                  # 태그 목록
GET  /tags/:tag             # 태그별 글 목록
GET  /categories/:category  # 카테고리별 글 목록
GET  /search                # 검색 (Bleve)
GET  /explore               # 기술 탐색 UI
GET  /feed.xml              # RSS 피드
GET  /api/search            # 검색 API (HTMX)
GET  /api/explore           # 탐색 API (HTMX)
```

---

## 7. Directory Structure (Go Project)

```
techpulse-blog/
|-- cmd/
|   +-- server/
|       +-- main.go
|-- internal/
|   |-- handler/        # HTTP 핸들러
|   |-- service/        # 비즈니스 로직
|   |-- repository/     # 데이터 접근
|   |-- model/          # 데이터 모델
|   |-- markdown/       # 마크다운 파싱
|   +-- search/         # Bleve 검색
|-- templates/           # templ 파일
|-- static/              # CSS, JS, 이미지
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
**Goal**: 마크다운 파일을 읽어서 블로그로 렌더링

- Go 프로젝트 초기화
- 마크다운 파싱 + frontmatter 처리
- templ 기반 레이아웃 (메인, 글 목록, 글 상세)
- HTMX 기반 페이지네이션
- Tailwind CSS 스타일링
- Mermaid.js 클라이언트 렌더링
- 코드 하이라이팅
- 로컬 개발 환경 (hot reload)

### Phase 2: Obsidian Integration
**Goal**: Obsidian vault와 블로그 소스 연동

- vault 폴더 구조 확정
- frontmatter 규격 확정
- Git submodule 또는 단일 repo 결정
- Obsidian 플러그인 추천/설정

### Phase 3: RSS Collection
**Goal**: n8n으로 기술 기사 자동 수집

- n8n Docker 설정
- RSS 수집 워크플로우
- feeds.yaml 파서
- 키워드 필터링 로직
- inbox/ 마크다운 생성

### Phase 4: Tech Explorer
**Goal**: 수집된 기사 검색/탐색

- Bleve 인덱싱 (posts + inbox)
- 검색 UI (HTMX)
- 태그/카테고리 필터
- 관련 포스트 추천

### Phase 5: Deployment
**Goal**: Oracle Cloud 배포

- Docker 이미지 빌드
- docker-compose (Go + n8n + Caddy)
- Oracle Cloud VM 프로비저닝
- GitHub Actions CI/CD

### Phase 6: Domain
**Goal**: 커스텀 도메인 적용

- .dev 도메인 구매
- Cloudflare 설정
- HTTPS 확인

---

## 9. Non-functional Requirements

- **Performance**: 페이지 로드 < 200ms (SSR)
- **Cost**: 월 $0 운영 (도메인 제외)
- **Security**: HTTPS 필수, XSS 방지, CSP 헤더
- **SEO**: 메타 태그, OGP, sitemap.xml
- **Accessibility**: 시맨틱 HTML, 키보드 내비게이션
- **Backup**: Git 기반 콘텐츠 백업, SQLite 정기 백업

---

## 10. Open Questions

- [ ] Obsidian vault를 별도 repo로 관리할지, 블로그 repo 내 submodule로 할지
- [ ] 다크모드 지원 여부
- [ ] 댓글 기능 필요 여부 (giscus 등)
- [ ] 다국어 지원 (한/영) 필요 여부
- [ ] Analytics (Umami 등 self-hosted) 도입 여부
