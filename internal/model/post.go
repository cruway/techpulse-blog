package model

import "time"

// @MX:ANCHOR: [AUTO] 後続全Feat（#9-#15）が依存する基盤モデル
// @MX:REASON: fan_in >= 6（Repository, Service, Handler, Markdown, Search, Template）

// PostStatus はポストの公開状態を表します。
type PostStatus string

const (
	// PostStatusDraft は下書き状態です。
	PostStatusDraft PostStatus = "draft"

	// PostStatusReview はレビュー待ち状態です。
	PostStatusReview PostStatus = "review"

	// PostStatusPublished は公開状態です。
	PostStatusPublished PostStatus = "published"
)

// IsValid はPostStatusが有効な値かどうかを返します。
func (s PostStatus) IsValid() bool {
	switch s {
	case PostStatusDraft, PostStatusReview, PostStatusPublished:
		return true
	default:
		return false
	}
}

// Post はブログポストを表す構造体です。
//
// マークダウン原文（Content）とレンダリング済みHTML（HTML）の両方を保持し、
// frontmatterから抽出されたメタデータ（Title, Tags, Category等）を含みます。
type Post struct {
	Date      time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	Category  string
	Title     string
	Slug      string
	Content   string
	HTML      string
	Excerpt   string
	Status    PostStatus
	SourceURL string
	Tags      []string
	ID        int64
	Mermaid   bool
}

// Validate はPostの必須フィールドを検証します。
func (p *Post) Validate() error {
	if p.Title == "" {
		return ErrInvalidTitle
	}
	if p.Slug == "" {
		return ErrInvalidSlug
	}
	if !p.Status.IsValid() {
		return ErrInvalidStatus
	}
	return nil
}

// Tag はタグ情報を表す構造体です。
type Tag struct {
	// Name はタグ名です。
	Name string

	// Count はこのタグが付けられたポスト数です。
	Count int
}

// Category はカテゴリ情報を表す構造体です。
type Category struct {
	// Name はカテゴリ名です。
	Name string

	// Count はこのカテゴリに属するポスト数です。
	Count int
}
