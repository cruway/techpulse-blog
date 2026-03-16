package repository

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite" // SQLiteドライバ登録
)

// スキーマ定義（個別ステートメント — SQLiteは複数ステートメント一括実行に制限あり）
var schemaMigrations = []string{
	`CREATE TABLE IF NOT EXISTS posts (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		title       TEXT NOT NULL,
		slug        TEXT NOT NULL UNIQUE,
		content     TEXT NOT NULL,
		html        TEXT NOT NULL,
		excerpt     TEXT NOT NULL DEFAULT '',
		date        DATETIME NOT NULL,
		category    TEXT NOT NULL DEFAULT '',
		status      TEXT NOT NULL DEFAULT 'draft',
		source_url  TEXT NOT NULL DEFAULT '',
		mermaid     BOOLEAN NOT NULL DEFAULT 0,
		created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`,
	`CREATE TABLE IF NOT EXISTS post_tags (
		post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
		tag     TEXT NOT NULL,
		PRIMARY KEY (post_id, tag)
	)`,
	`CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug)`,
	`CREATE INDEX IF NOT EXISTS idx_posts_date ON posts(date DESC)`,
	`CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status)`,
	`CREATE INDEX IF NOT EXISTS idx_posts_category ON posts(category)`,
	`CREATE INDEX IF NOT EXISTS idx_post_tags_tag ON post_tags(tag)`,
}

// NewDB はSQLiteデータベース接続を初期化します。
// WALモードを有効にし、スキーマを自動マイグレーションします。
func NewDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("DB接続失敗: %w", err)
	}

	// WALモード有効化
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return nil, fmt.Errorf("WALモード設定失敗: %w", err)
	}

	// 外部キー制約有効化
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		return nil, fmt.Errorf("外部キー設定失敗: %w", err)
	}

	// スキーマ適用（個別ステートメント実行）
	for _, stmt := range schemaMigrations {
		if _, err := db.Exec(stmt); err != nil {
			return nil, fmt.Errorf("スキーマ適用失敗: %w", err)
		}
	}

	return db, nil
}

// NewTestDB はテスト用のインメモリSQLiteを初期化します。
// 共有キャッシュモードで接続プール内の全接続が同一DBを参照します。
func NewTestDB() (*sql.DB, error) {
	return NewDB("file::memory:?cache=shared")
}
