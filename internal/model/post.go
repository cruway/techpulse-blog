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
	// ID はデータベース上の一意な識別子です。
	ID int64

	// Title はポストのタイトルです。
	Title string

	// Slug はURLに使用されるスラッグです。
	Slug string

	// Content はマークダウン原文です。
	Content string

	// HTML はレンダリング済みHTMLです。
	HTML string

	// Excerpt は一覧表示用の抜粋です。
	Excerpt string

	// Date は記事の公開日です。
	Date time.Time

	// Tags はポストに紐づくタグ一覧です。
	Tags []string

	// Category はポストのカテゴリです。
	Category string

	// Status はポストの公開状態です。
	Status PostStatus

	// SourceURL はObsidian元ファイルパスです（Phase 2用）。
	SourceURL string

	// Mermaid はMermaidコードブロックの含有フラグです。
	Mermaid bool

	// CreatedAt はレコード作成日時です。
	CreatedAt time.Time

	// UpdatedAt はレコード更新日時です。
	UpdatedAt time.Time
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
