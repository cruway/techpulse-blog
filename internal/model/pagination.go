package model

const (
	// DefaultPageSize はデフォルトのページサイズです。
	DefaultPageSize = 10

	// MaxPageSize は最大ページサイズです。
	MaxPageSize = 100
)

// ListOptions はポスト一覧取得時のフィルタ・ページネーション条件です。
type ListOptions struct {
	// Tag はタグフィルタです。空の場合はフィルタなし。
	Tag string

	// Category はカテゴリフィルタです。空の場合はフィルタなし。
	Category string

	// Status はステータスフィルタです。空の場合はフィルタなし。
	Status PostStatus

	// SortBy はソートキーです（"date", "title", "updated"）。
	SortBy string

	// Page はページ番号（1始まり）です。
	Page int

	// PageSize は1ページあたりの件数です。
	PageSize int

	// SortDesc がtrueの場合、降順ソートします。
	SortDesc bool
}

// Validate はListOptionsの値を検証します。
func (o *ListOptions) Validate() error {
	if o.Page < 1 {
		return ErrInvalidPage
	}
	if o.PageSize < 1 || o.PageSize > MaxPageSize {
		return ErrInvalidPageSize
	}
	return nil
}

// PageResult はページネーション付きの結果セットです。
type PageResult[T any] struct {
	// Items は現在のページのアイテム一覧です。
	Items []T

	// TotalItems はフィルタ条件に合致する全アイテム数です。
	TotalItems int

	// TotalPages は全ページ数です。
	TotalPages int

	// Page は現在のページ番号です。
	Page int

	// PageSize は1ページあたりの件数です。
	PageSize int
}

// NewPageResult はPageResultを生成します。
// TotalPagesはtotalItemsとpageSizeから自動計算されます。
func NewPageResult[T any](items []T, totalItems, page, pageSize int) PageResult[T] {
	totalPages := 0
	if totalItems > 0 && pageSize > 0 {
		totalPages = (totalItems + pageSize - 1) / pageSize
	}

	return PageResult[T]{
		Items:      items,
		TotalItems: totalItems,
		TotalPages: totalPages,
		Page:       page,
		PageSize:   pageSize,
	}
}
