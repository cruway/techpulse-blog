package model

import "time"

// FeedItem はRSSフィードから取得した記事を表す構造体です（Phase 3用）。
type FeedItem struct {
	// PublishedAt は記事の公開日時です。
	PublishedAt time.Time

	// FetchedAt はフィード取得日時です。
	FetchedAt time.Time

	// FeedURL はフィードの取得元URLです。
	FeedURL string

	// Title は記事のタイトルです。
	Title string

	// Link は記事のURLです。
	Link string

	// Description は記事の要約です。
	Description string

	// ID はデータベース上の一意な識別子です。
	ID int64
}

// Validate はFeedItemの必須フィールドを検証します。
func (f *FeedItem) Validate() error {
	if f.FeedURL == "" || f.Title == "" || f.Link == "" {
		return ErrInvalidFeedItem
	}
	return nil
}
