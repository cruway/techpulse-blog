package repository

import (
	"fmt"
	"testing"
	"time"

	"github.com/cruway/techpulse-blog/internal/model"
)

// newTestPost はテスト用のPost構造体を生成します。
func newTestPost(slug, title, category string, tags []string, status model.PostStatus) *model.Post {
	return &model.Post{
		Title:    title,
		Slug:     slug,
		Content:  "# " + title,
		HTML:     "<h1>" + title + "</h1>",
		Excerpt:  title,
		Date:     time.Date(2026, 3, 16, 0, 0, 0, 0, time.UTC),
		Tags:     tags,
		Category: category,
		Status:   status,
	}
}

// setupTestDB はテストDBとリポジトリを初期化します。
func setupTestDB(t *testing.T) PostRepository {
	t.Helper()
	db, err := NewTestDB()
	if err != nil {
		t.Fatalf("テストDB初期化失敗: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return NewSQLitePostRepository(db)
}

// mustUpsert はテスト用のUpsertヘルパーです。エラー時にt.Fatalを呼びます。
func mustUpsert(t *testing.T, repo PostRepository, post *model.Post) {
	t.Helper()
	if err := repo.Upsert(post); err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}
}

// TestUpsert_新規作成 はポストの新規作成をテストします。
func TestUpsert_新規作成(t *testing.T) {
	repo := setupTestDB(t)
	post := newTestPost("test-post", "テスト記事", "技術", []string{"go", "echo"}, model.PostStatusPublished)

	if err := repo.Upsert(post); err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}

	got, err := repo.FindBySlug("test-post")
	if err != nil {
		t.Fatalf("FindBySlug() error = %v", err)
	}
	if got.Title != "テスト記事" {
		t.Errorf("Title = %q, 期待値 %q", got.Title, "テスト記事")
	}
	if len(got.Tags) != 2 {
		t.Errorf("Tags len = %d, 期待値 2", len(got.Tags))
	}
}

// TestUpsert_更新 は同一slugでの更新をテストします。
func TestUpsert_更新(t *testing.T) {
	repo := setupTestDB(t)
	post := newTestPost("update-post", "初回", "技術", []string{"go"}, model.PostStatusDraft)
	if err := repo.Upsert(post); err != nil {
		t.Fatalf("Upsert(1回目) error = %v", err)
	}

	// 更新
	post.Title = "更新後"
	post.Tags = []string{"go", "htmx"}
	if err := repo.Upsert(post); err != nil {
		t.Fatalf("Upsert(2回目) error = %v", err)
	}

	got, err := repo.FindBySlug("update-post")
	if err != nil {
		t.Fatalf("FindBySlug() error = %v", err)
	}
	if got.Title != "更新後" {
		t.Errorf("Title = %q, 期待値 %q", got.Title, "更新後")
	}
	if len(got.Tags) != 2 {
		t.Errorf("Tags len = %d, 期待値 2", len(got.Tags))
	}
}

// TestFindBySlug_存在しない は存在しないスラッグの検索をテストします。
func TestFindBySlug_存在しない(t *testing.T) {
	repo := setupTestDB(t)
	_, err := repo.FindBySlug("nonexistent")
	if err != model.ErrPostNotFound {
		t.Errorf("FindBySlug() error = %v, 期待値 %v", err, model.ErrPostNotFound)
	}
}

// TestFindAll_ページネーション はページネーション動作をテストします。
func TestFindAll_ページネーション(t *testing.T) {
	repo := setupTestDB(t)

	// 5件作成
	for i := 0; i < 5; i++ {
		post := newTestPost(
			fmt.Sprintf("post-%d", i),
			fmt.Sprintf("記事 %d", i),
			"技術",
			[]string{"go"},
			model.PostStatusPublished,
		)
		post.Date = time.Date(2026, 3, 16-i, 0, 0, 0, 0, time.UTC)
		if err := repo.Upsert(post); err != nil {
			t.Fatalf("Upsert() error = %v", err)
		}
	}

	result, err := repo.FindAll(model.ListOptions{Page: 1, PageSize: 2})
	if err != nil {
		t.Fatalf("FindAll() error = %v", err)
	}

	if result.TotalItems != 5 {
		t.Errorf("TotalItems = %d, 期待値 5", result.TotalItems)
	}
	if result.TotalPages != 3 {
		t.Errorf("TotalPages = %d, 期待値 3", result.TotalPages)
	}
	if len(result.Items) != 2 {
		t.Errorf("Items len = %d, 期待値 2", len(result.Items))
	}
}

// TestFindAll_ステータスフィルタ はステータスフィルタをテストします。
func TestFindAll_ステータスフィルタ(t *testing.T) {
	repo := setupTestDB(t)

	mustUpsert(t, repo, newTestPost("pub-1", "公開1", "技術", nil, model.PostStatusPublished))
	mustUpsert(t, repo, newTestPost("draft-1", "下書き1", "技術", nil, model.PostStatusDraft))
	mustUpsert(t, repo, newTestPost("pub-2", "公開2", "技術", nil, model.PostStatusPublished))

	result, err := repo.FindAll(model.ListOptions{
		Page:     1,
		PageSize: 10,
		Status:   model.PostStatusPublished,
	})
	if err != nil {
		t.Fatalf("FindAll() error = %v", err)
	}

	if result.TotalItems != 2 {
		t.Errorf("TotalItems = %d, 期待値 2（publishedのみ）", result.TotalItems)
	}
}

// TestFindByTag はタグ別検索をテストします。
func TestFindByTag(t *testing.T) {
	repo := setupTestDB(t)

	mustUpsert(t, repo, newTestPost("go-1", "Go記事", "技術", []string{"go", "echo"}, model.PostStatusPublished))
	mustUpsert(t, repo, newTestPost("rust-1", "Rust記事", "技術", []string{"rust"}, model.PostStatusPublished))
	mustUpsert(t, repo, newTestPost("go-2", "Go記事2", "技術", []string{"go"}, model.PostStatusPublished))

	result, err := repo.FindByTag("go", model.ListOptions{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("FindByTag() error = %v", err)
	}

	if result.TotalItems != 2 {
		t.Errorf("TotalItems = %d, 期待値 2", result.TotalItems)
	}
}

// TestFindByCategory はカテゴリ別検索をテストします。
func TestFindByCategory(t *testing.T) {
	repo := setupTestDB(t)

	mustUpsert(t, repo, newTestPost("tech-1", "技術記事", "技術", nil, model.PostStatusPublished))
	mustUpsert(t, repo, newTestPost("life-1", "ライフ記事", "ライフ", nil, model.PostStatusPublished))

	result, err := repo.FindByCategory("技術", model.ListOptions{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("FindByCategory() error = %v", err)
	}

	if result.TotalItems != 1 {
		t.Errorf("TotalItems = %d, 期待値 1", result.TotalItems)
	}
}

// TestAllTags はタグ一覧取得をテストします。
func TestAllTags(t *testing.T) {
	repo := setupTestDB(t)

	mustUpsert(t, repo, newTestPost("p1", "記事1", "技術", []string{"go", "echo"}, model.PostStatusPublished))
	mustUpsert(t, repo, newTestPost("p2", "記事2", "技術", []string{"go", "htmx"}, model.PostStatusPublished))

	tags, err := repo.AllTags()
	if err != nil {
		t.Fatalf("AllTags() error = %v", err)
	}

	if len(tags) != 3 {
		t.Fatalf("Tags len = %d, 期待値 3", len(tags))
	}

	// goが2件で最初に来るはず（カウント降順）
	if tags[0].Name != "go" || tags[0].Count != 2 {
		t.Errorf("Tags[0] = {%q, %d}, 期待値 {\"go\", 2}", tags[0].Name, tags[0].Count)
	}
}

// TestAllCategories はカテゴリ一覧取得をテストします。
func TestAllCategories(t *testing.T) {
	repo := setupTestDB(t)

	mustUpsert(t, repo, newTestPost("p1", "記事1", "技術", nil, model.PostStatusPublished))
	mustUpsert(t, repo, newTestPost("p2", "記事2", "ライフ", nil, model.PostStatusPublished))
	mustUpsert(t, repo, newTestPost("p3", "記事3", "技術", nil, model.PostStatusPublished))

	categories, err := repo.AllCategories()
	if err != nil {
		t.Fatalf("AllCategories() error = %v", err)
	}

	if len(categories) != 2 {
		t.Fatalf("Categories len = %d, 期待値 2", len(categories))
	}
}

// TestDelete は削除とカスケードをテストします。
func TestDelete(t *testing.T) {
	repo := setupTestDB(t)

	mustUpsert(t, repo, newTestPost("del-post", "削除対象", "技術", []string{"go"}, model.PostStatusPublished))

	if err := repo.Delete("del-post"); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err := repo.FindBySlug("del-post")
	if err != model.ErrPostNotFound {
		t.Errorf("削除後のFindBySlug() error = %v, 期待値 %v", err, model.ErrPostNotFound)
	}
}

// TestAllSlugs は全スラッグ取得をテストします。
func TestAllSlugs(t *testing.T) {
	repo := setupTestDB(t)

	mustUpsert(t, repo, newTestPost("slug-a", "A", "技術", nil, model.PostStatusPublished))
	mustUpsert(t, repo, newTestPost("slug-b", "B", "技術", nil, model.PostStatusPublished))

	slugs, err := repo.AllSlugs()
	if err != nil {
		t.Fatalf("AllSlugs() error = %v", err)
	}

	if len(slugs) != 2 {
		t.Errorf("Slugs len = %d, 期待値 2", len(slugs))
	}
}
