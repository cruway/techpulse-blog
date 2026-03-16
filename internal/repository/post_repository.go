// Package repository はデータアクセス層を提供します。
package repository

import (
	"github.com/cruway/techpulse-blog/internal/model"
)

// PostRepository はポストのデータアクセスインターフェースです。
type PostRepository interface {
	// FindBySlug はスラッグでポストを検索します。
	FindBySlug(slug string) (*model.Post, error)

	// FindAll はポスト一覧を取得します。
	FindAll(opts model.ListOptions) (*model.PageResult[*model.Post], error)

	// FindByTag はタグ別ポスト一覧を取得します。
	FindByTag(tag string, opts model.ListOptions) (*model.PageResult[*model.Post], error)

	// FindByCategory はカテゴリ別ポスト一覧を取得します。
	FindByCategory(category string, opts model.ListOptions) (*model.PageResult[*model.Post], error)

	// AllTags は全タグをカウント付きで取得します。
	AllTags() ([]model.Tag, error)

	// AllCategories は全カテゴリをカウント付きで取得します。
	AllCategories() ([]model.Category, error)

	// Upsert はポストを作成または更新します。
	Upsert(post *model.Post) error

	// Delete はスラッグでポストを削除します。
	Delete(slug string) error

	// AllSlugs はDBに存在する全スラッグを取得します。
	AllSlugs() ([]string, error)
}
