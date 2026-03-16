package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/cruway/techpulse-blog/internal/model"
)

// @MX:ANCHOR: [AUTO] Service層から呼び出されるデータアクセス基盤
// @MX:REASON: fan_in >= 3（PostService, SyncContent, Handler経由）

// SQLitePostRepository はSQLiteベースのPostRepository実装です。
type SQLitePostRepository struct {
	db *sql.DB
}

// NewSQLitePostRepository は新しいSQLitePostRepositoryを生成します。
func NewSQLitePostRepository(db *sql.DB) *SQLitePostRepository {
	return &SQLitePostRepository{db: db}
}

// FindBySlug はスラッグでポストを検索します。
func (r *SQLitePostRepository) FindBySlug(slug string) (*model.Post, error) {
	post := &model.Post{}
	err := r.db.QueryRow(
		`SELECT id, title, slug, content, html, excerpt, date, category, status, source_url, mermaid, created_at, updated_at
		 FROM posts WHERE slug = ?`, slug,
	).Scan(
		&post.ID, &post.Title, &post.Slug, &post.Content, &post.HTML, &post.Excerpt,
		&post.Date, &post.Category, &post.Status, &post.SourceURL, &post.Mermaid,
		&post.CreatedAt, &post.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, model.ErrPostNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("ポスト検索失敗 (slug=%s): %w", slug, err)
	}

	tags, err := r.findTagsByPostID(post.ID)
	if err != nil {
		return nil, err
	}
	post.Tags = tags

	return post, nil
}

// FindAll はポスト一覧を取得します。
func (r *SQLitePostRepository) FindAll(opts model.ListOptions) (*model.PageResult[*model.Post], error) {
	where, args := buildWhereClause(opts)
	return r.findPosts(where, args, opts)
}

// FindByTag はタグ別ポスト一覧を取得します。
func (r *SQLitePostRepository) FindByTag(tag string, opts model.ListOptions) (*model.PageResult[*model.Post], error) {
	where, args := buildWhereClause(opts)
	where = append(where, "p.id IN (SELECT post_id FROM post_tags WHERE tag = ?)")
	args = append(args, tag)
	return r.findPosts(where, args, opts)
}

// FindByCategory はカテゴリ別ポスト一覧を取得します。
func (r *SQLitePostRepository) FindByCategory(category string, opts model.ListOptions) (*model.PageResult[*model.Post], error) {
	where, args := buildWhereClause(opts)
	where = append(where, "p.category = ?")
	args = append(args, category)
	return r.findPosts(where, args, opts)
}

// AllTags は全タグをカウント付きで取得します（カウント降順）。
func (r *SQLitePostRepository) AllTags() ([]model.Tag, error) {
	rows, err := r.db.Query(
		`SELECT tag, COUNT(*) as cnt FROM post_tags GROUP BY tag ORDER BY cnt DESC, tag ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("タグ一覧取得失敗: %w", err)
	}
	defer rows.Close()

	var tags []model.Tag
	for rows.Next() {
		var tag model.Tag
		if err := rows.Scan(&tag.Name, &tag.Count); err != nil {
			return nil, fmt.Errorf("タグ読み取り失敗: %w", err)
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}

// AllCategories は全カテゴリをカウント付きで取得します（カウント降順）。
func (r *SQLitePostRepository) AllCategories() ([]model.Category, error) {
	rows, err := r.db.Query(
		`SELECT category, COUNT(*) as cnt FROM posts WHERE category != '' GROUP BY category ORDER BY cnt DESC, category ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("カテゴリ一覧取得失敗: %w", err)
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var cat model.Category
		if err := rows.Scan(&cat.Name, &cat.Count); err != nil {
			return nil, fmt.Errorf("カテゴリ読み取り失敗: %w", err)
		}
		categories = append(categories, cat)
	}
	return categories, rows.Err()
}

// Upsert はポストを作成または更新します。
func (r *SQLitePostRepository) Upsert(post *model.Post) error {
	now := time.Now()

	result, err := r.db.Exec(
		`INSERT INTO posts (title, slug, content, html, excerpt, date, category, status, source_url, mermaid, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		 ON CONFLICT(slug) DO UPDATE SET
		   title=excluded.title, content=excluded.content, html=excluded.html,
		   excerpt=excluded.excerpt, date=excluded.date, category=excluded.category,
		   status=excluded.status, source_url=excluded.source_url, mermaid=excluded.mermaid,
		   updated_at=excluded.updated_at`,
		post.Title, post.Slug, post.Content, post.HTML, post.Excerpt,
		post.Date, post.Category, post.Status, post.SourceURL, post.Mermaid,
		now, now,
	)
	if err != nil {
		return fmt.Errorf("ポスト保存失敗 (slug=%s): %w", post.Slug, err)
	}

	// post IDを取得（INSERT時はLastInsertId、UPDATE時はSELECT）
	postID, err := result.LastInsertId()
	if err != nil || postID == 0 {
		err = r.db.QueryRow("SELECT id FROM posts WHERE slug = ?", post.Slug).Scan(&postID)
		if err != nil {
			return fmt.Errorf("ポストID取得失敗 (slug=%s): %w", post.Slug, err)
		}
	}

	// タグ更新（全削除 → 再挿入）
	if _, err := r.db.Exec("DELETE FROM post_tags WHERE post_id = ?", postID); err != nil {
		return fmt.Errorf("タグ削除失敗: %w", err)
	}

	for _, tag := range post.Tags {
		if _, err := r.db.Exec("INSERT INTO post_tags (post_id, tag) VALUES (?, ?)", postID, tag); err != nil {
			return fmt.Errorf("タグ挿入失敗 (tag=%s): %w", tag, err)
		}
	}

	return nil
}

// Delete はスラッグでポストを削除します。
func (r *SQLitePostRepository) Delete(slug string) error {
	_, err := r.db.Exec("DELETE FROM posts WHERE slug = ?", slug)
	if err != nil {
		return fmt.Errorf("ポスト削除失敗 (slug=%s): %w", slug, err)
	}
	return nil
}

// AllSlugs はDBに存在する全スラッグを取得します。
func (r *SQLitePostRepository) AllSlugs() ([]string, error) {
	rows, err := r.db.Query("SELECT slug FROM posts")
	if err != nil {
		return nil, fmt.Errorf("スラッグ一覧取得失敗: %w", err)
	}
	defer rows.Close()

	var slugs []string
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return nil, fmt.Errorf("スラッグ読み取り失敗: %w", err)
		}
		slugs = append(slugs, slug)
	}
	return slugs, rows.Err()
}

// findTagsByPostID はポストIDに紐づくタグを取得します。
func (r *SQLitePostRepository) findTagsByPostID(postID int64) ([]string, error) {
	rows, err := r.db.Query("SELECT tag FROM post_tags WHERE post_id = ? ORDER BY tag", postID)
	if err != nil {
		return nil, fmt.Errorf("タグ取得失敗 (post_id=%d): %w", postID, err)
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, fmt.Errorf("タグ読み取り失敗: %w", err)
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}

// findPosts は条件に合致するポスト一覧を取得する内部メソッドです。
func (r *SQLitePostRepository) findPosts(where []string, args []interface{}, opts model.ListOptions) (*model.PageResult[*model.Post], error) {
	// 件数取得
	countQuery := "SELECT COUNT(*) FROM posts p"
	if len(where) > 0 {
		countQuery += " WHERE " + strings.Join(where, " AND ")
	}

	var totalItems int
	if err := r.db.QueryRow(countQuery, args...).Scan(&totalItems); err != nil {
		return nil, fmt.Errorf("件数取得失敗: %w", err)
	}

	// データ取得
	query := `SELECT p.id, p.title, p.slug, p.content, p.html, p.excerpt, p.date, p.category, p.status, p.source_url, p.mermaid, p.created_at, p.updated_at
	          FROM posts p`
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	query += " ORDER BY p.date DESC"

	offset := (opts.Page - 1) * opts.PageSize
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", opts.PageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("ポスト一覧取得失敗: %w", err)
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		post := &model.Post{}
		if err := rows.Scan(
			&post.ID, &post.Title, &post.Slug, &post.Content, &post.HTML, &post.Excerpt,
			&post.Date, &post.Category, &post.Status, &post.SourceURL, &post.Mermaid,
			&post.CreatedAt, &post.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("ポスト読み取り失敗: %w", err)
		}

		tags, err := r.findTagsByPostID(post.ID)
		if err != nil {
			return nil, err
		}
		post.Tags = tags

		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ポスト一覧走査失敗: %w", err)
	}

	result := model.NewPageResult(posts, totalItems, opts.Page, opts.PageSize)
	return &result, nil
}

// buildWhereClause はListOptionsからWHERE条件を構築します。
func buildWhereClause(opts model.ListOptions) ([]string, []interface{}) {
	var where []string
	var args []interface{}

	if opts.Status != "" {
		where = append(where, "p.status = ?")
		args = append(args, string(opts.Status))
	}

	return where, args
}
