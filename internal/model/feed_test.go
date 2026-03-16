package model

import (
	"testing"
	"time"
)

// TestFeedItemValidate はFeedItemのバリデーションをテストします。
func TestFeedItemValidate(t *testing.T) {
	validItem := FeedItem{
		FeedURL:     "https://example.com/feed",
		Title:       "テスト記事",
		Link:        "https://example.com/article",
		PublishedAt: time.Now(),
	}

	tests := []struct {
		name    string
		modify  func(f *FeedItem)
		wantErr bool
	}{
		{
			name:    "正常なフィードアイテム",
			modify:  func(f *FeedItem) {},
			wantErr: false,
		},
		{
			name:    "FeedURL空",
			modify:  func(f *FeedItem) { f.FeedURL = "" },
			wantErr: true,
		},
		{
			name:    "Title空",
			modify:  func(f *FeedItem) { f.Title = "" },
			wantErr: true,
		},
		{
			name:    "Link空",
			modify:  func(f *FeedItem) { f.Link = "" },
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := validItem
			tt.modify(&f)
			err := f.Validate()

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
