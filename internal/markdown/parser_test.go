package markdown

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cruway/techpulse-blog/internal/model"
)

// TestParseContent_正常系 は基本的なマークダウン変換をテストします。
func TestParseContent_正常系(t *testing.T) {
	input := []byte(`---
title: "テスト記事"
date: 2026-03-16
tags: [go, htmx]
category: "技術"
status: draft
---
# 見出し

本文テキストです。
`)

	p := NewParser()
	post, err := p.ParseContent(input)
	if err != nil {
		t.Fatalf("ParseContent() error = %v", err)
	}

	if post.Title != "テスト記事" {
		t.Errorf("Title = %q, 期待値 %q", post.Title, "テスト記事")
	}
	if post.Category != "技術" {
		t.Errorf("Category = %q, 期待値 %q", post.Category, "技術")
	}
	if post.Status != model.PostStatusDraft {
		t.Errorf("Status = %q, 期待値 %q", post.Status, model.PostStatusDraft)
	}
	if !strings.Contains(post.HTML, "<h1") {
		t.Errorf("HTML に <h1> が含まれていません: %s", post.HTML)
	}
	if post.Content == "" {
		t.Error("Content が空です")
	}
}

// TestParseContent_全フィールド はfrontmatter全フィールドのマッピングをテストします。
func TestParseContent_全フィールド(t *testing.T) {
	input := []byte(`---
title: "全フィールドテスト"
date: 2026-01-15
tags: [rust, wasm, frontend]
category: "チュートリアル"
status: published
source_url: "/path/to/obsidian/note.md"
mermaid: true
---
本文
`)

	p := NewParser()
	post, err := p.ParseContent(input)
	if err != nil {
		t.Fatalf("ParseContent() error = %v", err)
	}

	tests := []struct {
		name string
		got  string
		want string
	}{
		{"Title", post.Title, "全フィールドテスト"},
		{"Category", post.Category, "チュートリアル"},
		{"Status", string(post.Status), "published"},
		{"SourceURL", post.SourceURL, "/path/to/obsidian/note.md"},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %q, 期待値 %q", tt.name, tt.got, tt.want)
		}
	}

	if !post.Mermaid {
		t.Error("Mermaid = false, 期待値 true")
	}
	if post.Date.Year() != 2026 || post.Date.Month() != 1 || post.Date.Day() != 15 {
		t.Errorf("Date = %v, 期待値 2026-01-15", post.Date)
	}
}

// TestParseContent_タグ変換 はタグの型変換をテストします。
func TestParseContent_タグ変換(t *testing.T) {
	input := []byte(`---
title: "タグテスト"
date: 2026-03-16
tags: [go, echo, htmx]
category: "技術"
status: draft
---
本文
`)

	p := NewParser()
	post, err := p.ParseContent(input)
	if err != nil {
		t.Fatalf("ParseContent() error = %v", err)
	}

	wantTags := []string{"go", "echo", "htmx"}
	if len(post.Tags) != len(wantTags) {
		t.Fatalf("Tags len = %d, 期待値 %d", len(post.Tags), len(wantTags))
	}
	for i, tag := range wantTags {
		if post.Tags[i] != tag {
			t.Errorf("Tags[%d] = %q, 期待値 %q", i, post.Tags[i], tag)
		}
	}
}

// TestParseContent_Mermaid検知 はMermaidコードブロックの検知と保持をテストします。
func TestParseContent_Mermaid検知(t *testing.T) {
	input := []byte("---\ntitle: \"Mermaidテスト\"\ndate: 2026-03-16\ntags: []\ncategory: \"技術\"\nstatus: draft\n---\n\n```mermaid\ngraph TD\n    A --> B\n```\n")

	p := NewParser()
	post, err := p.ParseContent(input)
	if err != nil {
		t.Fatalf("ParseContent() error = %v", err)
	}

	if !post.Mermaid {
		t.Error("Mermaid = false, 期待値 true（コードブロック検知）")
	}
	if !strings.Contains(post.HTML, `class="mermaid"`) {
		t.Errorf("HTML に class=\"mermaid\" が含まれていません: %s", post.HTML)
	}
}

// TestParseContent_GFM はGFM拡張（テーブル）の変換をテストします。
func TestParseContent_GFM(t *testing.T) {
	input := []byte(`---
title: "GFMテスト"
date: 2026-03-16
tags: []
category: "技術"
status: draft
---

| 列1 | 列2 |
|-----|-----|
| A   | B   |
`)

	p := NewParser()
	post, err := p.ParseContent(input)
	if err != nil {
		t.Fatalf("ParseContent() error = %v", err)
	}

	if !strings.Contains(post.HTML, "<table") {
		t.Errorf("HTML に <table> が含まれていません: %s", post.HTML)
	}
}

// TestParseContent_frontmatterなし はfrontmatterがない場合のエラーをテストします。
func TestParseContent_frontmatterなし(t *testing.T) {
	input := []byte("# タイトルのみ\n\n本文テキスト")

	p := NewParser()
	_, err := p.ParseContent(input)
	if err == nil {
		t.Error("ParseContent() error = nil, エラーが期待されます")
	}
}

// TestParseContent_空コンテンツ は空のコンテンツのエラーをテストします。
func TestParseContent_空コンテンツ(t *testing.T) {
	p := NewParser()
	_, err := p.ParseContent([]byte{})
	if err != ErrEmptyContent {
		t.Errorf("ParseContent() error = %v, 期待値 %v", err, ErrEmptyContent)
	}
}

// TestParseContent_Excerpt生成 は抜粋の生成をテストします。
func TestParseContent_Excerpt生成(t *testing.T) {
	// 200文字以上の本文
	longText := strings.Repeat("あ", 300)
	input := []byte("---\ntitle: \"Excerptテスト\"\ndate: 2026-03-16\ntags: []\ncategory: \"技術\"\nstatus: draft\n---\n\n" + longText + "\n")

	p := NewParser()
	post, err := p.ParseContent(input)
	if err != nil {
		t.Fatalf("ParseContent() error = %v", err)
	}

	if post.Excerpt == "" {
		t.Error("Excerpt が空です")
	}
	// Excerptは200文字以内（runeベース）
	if len([]rune(post.Excerpt)) > 200 {
		t.Errorf("Excerpt長 = %d runes, 期待値 <= 200", len([]rune(post.Excerpt)))
	}
}

// TestParseFile_正常系 はファイルからの読み込みをテストします。
func TestParseFile_正常系(t *testing.T) {
	// 一時ファイル作成
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")
	content := []byte(`---
title: "ファイルテスト"
date: 2026-03-16
tags: [test]
category: "テスト"
status: draft
---
# ファイルからの読み込み

テスト本文です。
`)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("ファイル作成失敗: %v", err)
	}

	p := NewParser()
	post, err := p.ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}

	if post.Title != "ファイルテスト" {
		t.Errorf("Title = %q, 期待値 %q", post.Title, "ファイルテスト")
	}
}

// TestParseFile_存在しないファイル はファイル不在のエラーをテストします。
func TestParseFile_存在しないファイル(t *testing.T) {
	p := NewParser()
	_, err := p.ParseFile("/nonexistent/path/file.md")
	if err == nil {
		t.Error("ParseFile() error = nil, エラーが期待されます")
	}
}
