// Package model はTechPulse Blogのドメインモデルを定義します。
package model

import "errors"

// ドメインエラー定義
var (
	// ErrPostNotFound はポストが見つからない場合のエラーです。
	ErrPostNotFound = errors.New("ポストが見つかりません")

	// ErrInvalidTitle は無効なタイトルのエラーです。
	ErrInvalidTitle = errors.New("無効なタイトルです")

	// ErrInvalidSlug は無効なスラッグのエラーです。
	ErrInvalidSlug = errors.New("無効なスラッグです")

	// ErrInvalidStatus は無効なステータスのエラーです。
	ErrInvalidStatus = errors.New("無効なステータスです")

	// ErrInvalidPage は無効なページ番号のエラーです。
	ErrInvalidPage = errors.New("無効なページ番号です")

	// ErrInvalidPageSize は無効なページサイズのエラーです。
	ErrInvalidPageSize = errors.New("無効なページサイズです")

	// ErrInvalidFeedItem は無効なフィードアイテムのエラーです。
	ErrInvalidFeedItem = errors.New("無効なフィードアイテムです")
)
