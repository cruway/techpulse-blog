package model

import (
	"testing"
)

// TestListOptionsValidate はListOptionsのバリデーションをテストします。
func TestListOptionsValidate(t *testing.T) {
	tests := []struct {
		name    string
		opts    ListOptions
		wantErr error
	}{
		{
			name:    "デフォルト値で正常",
			opts:    ListOptions{Page: 1, PageSize: 10},
			wantErr: nil,
		},
		{
			name:    "フィルタ付きで正常",
			opts:    ListOptions{Page: 1, PageSize: 20, Tag: "go", Category: "技術"},
			wantErr: nil,
		},
		{
			name:    "Page が0",
			opts:    ListOptions{Page: 0, PageSize: 10},
			wantErr: ErrInvalidPage,
		},
		{
			name:    "Page が負数",
			opts:    ListOptions{Page: -1, PageSize: 10},
			wantErr: ErrInvalidPage,
		},
		{
			name:    "PageSize が0",
			opts:    ListOptions{Page: 1, PageSize: 0},
			wantErr: ErrInvalidPageSize,
		},
		{
			name:    "PageSize が上限超過",
			opts:    ListOptions{Page: 1, PageSize: MaxPageSize + 1},
			wantErr: ErrInvalidPageSize,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()

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

// TestNewPageResult はPageResult生成のテストです。
func TestNewPageResult(t *testing.T) {
	tests := []struct {
		name       string
		items      []string
		totalItems int
		page       int
		pageSize   int
		wantPages  int
	}{
		{
			name:       "1ページに収まる",
			items:      []string{"a", "b"},
			totalItems: 2,
			page:       1,
			pageSize:   10,
			wantPages:  1,
		},
		{
			name:       "複数ページ",
			items:      []string{"a", "b", "c"},
			totalItems: 25,
			page:       1,
			pageSize:   10,
			wantPages:  3,
		},
		{
			name:       "ちょうど割り切れる",
			items:      []string{"a"},
			totalItems: 20,
			page:       2,
			pageSize:   10,
			wantPages:  2,
		},
		{
			name:       "アイテム0件",
			items:      []string{},
			totalItems: 0,
			page:       1,
			pageSize:   10,
			wantPages:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewPageResult(tt.items, tt.totalItems, tt.page, tt.pageSize)

			if result.TotalPages != tt.wantPages {
				t.Errorf("TotalPages = %d, 期待値 %d", result.TotalPages, tt.wantPages)
			}
			if result.Page != tt.page {
				t.Errorf("Page = %d, 期待値 %d", result.Page, tt.page)
			}
			if result.PageSize != tt.pageSize {
				t.Errorf("PageSize = %d, 期待値 %d", result.PageSize, tt.pageSize)
			}
			if result.TotalItems != tt.totalItems {
				t.Errorf("TotalItems = %d, 期待値 %d", result.TotalItems, tt.totalItems)
			}
			if len(result.Items) != len(tt.items) {
				t.Errorf("Items len = %d, 期待値 %d", len(result.Items), len(tt.items))
			}
		})
	}
}
