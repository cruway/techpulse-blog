package model

import (
	"testing"
	"time"
)

// TestPostStatusIsValid はPostStatusの有効性検証をテストします。
func TestPostStatusIsValid(t *testing.T) {
	tests := []struct {
		name   string
		status PostStatus
		want   bool
	}{
		{name: "draft は有効", status: PostStatusDraft, want: true},
		{name: "review は有効", status: PostStatusReview, want: true},
		{name: "published は有効", status: PostStatusPublished, want: true},
		{name: "空文字は無効", status: PostStatus(""), want: false},
		{name: "不明な値は無効", status: PostStatus("archived"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.IsValid(); got != tt.want {
				t.Errorf("PostStatus(%q).IsValid() = %v, 期待値 %v", tt.status, got, tt.want)
			}
		})
	}
}

// TestPostValidate はPost構造体のバリデーションをテストします。
func TestPostValidate(t *testing.T) {
	validPost := Post{
		Title:   "テスト記事",
		Slug:    "test-article",
		Content: "# テスト",
		Date:    time.Now(),
		Status:  PostStatusDraft,
	}

	tests := []struct {
		wantErr error
		modify  func(p *Post)
		name    string
	}{
		{
			name:    "正常なポスト",
			modify:  func(p *Post) {},
			wantErr: nil,
		},
		{
			name:    "タイトル空",
			modify:  func(p *Post) { p.Title = "" },
			wantErr: ErrInvalidTitle,
		},
		{
			name:    "スラッグ空",
			modify:  func(p *Post) { p.Slug = "" },
			wantErr: ErrInvalidSlug,
		},
		{
			name:    "無効なステータス",
			modify:  func(p *Post) { p.Status = PostStatus("unknown") },
			wantErr: ErrInvalidStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := validPost
			tt.modify(&p)
			err := p.Validate()

			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("Validate() = %v, 期待値 nil", err)
				}
				return
			}

			if err == nil {
				t.Errorf("Validate() = nil, 期待値 %v", tt.wantErr)
				return
			}

			if err != tt.wantErr {
				t.Errorf("Validate() = %v, 期待値 %v", err, tt.wantErr)
			}
		})
	}
}
