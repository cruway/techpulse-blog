// Package markdown はマークダウンのパースとHTML変換を提供します。
package markdown

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/cruway/techpulse-blog/internal/model"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// パーサーエラー定義
var (
	// ErrFileNotFound はファイルが見つからない場合のエラーです。
	ErrFileNotFound = errors.New("ファイルが見つかりません")

	// ErrInvalidFrontmatter は無効なfrontmatterのエラーです。
	ErrInvalidFrontmatter = errors.New("無効なfrontmatterです")

	// ErrEmptyContent はコンテンツが空の場合のエラーです。
	ErrEmptyContent = errors.New("コンテンツが空です")
)

// Excerpt最大長（rune単位）
const maxExcerptLength = 200

// mermaidコードブロック置換用正規表現
var mermaidCodeBlockRe = regexp.MustCompile(`<code class="language-mermaid">`)

// @MX:ANCHOR: [AUTO] Service/Handlerから呼び出される基盤パーサー
// @MX:REASON: fan_in >= 3（PostService, ContentSync, Handler）

// Parser はマークダウンをPost構造体に変換するパーサーです。
type Parser struct {
	md goldmark.Markdown
}

// NewParser は新しいParserを生成します。
//
// goldmarkをGFM拡張 + frontmatter対応で初期化します。
func NewParser() *Parser {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			meta.Meta,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithUnsafe(),
		),
	)
	return &Parser{md: md}
}

// ParseFile はファイルパスからPost構造体を生成します。
func (p *Parser) ParseFile(path string) (*model.Post, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", ErrFileNotFound, path)
		}
		return nil, fmt.Errorf("ファイル読み込み失敗: %w", err)
	}
	return p.ParseContent(content)
}

// ParseContent はバイト列からPost構造体を生成します。
func (p *Parser) ParseContent(content []byte) (*model.Post, error) {
	if len(content) == 0 {
		return nil, ErrEmptyContent
	}

	// goldmarkでパース
	var buf bytes.Buffer
	ctx := parser.NewContext()
	if err := p.md.Convert(content, &buf, parser.WithContext(ctx)); err != nil {
		return nil, fmt.Errorf("マークダウン変換失敗: %w", err)
	}

	// frontmatter取得
	metaData := meta.Get(ctx)
	if len(metaData) == 0 {
		return nil, ErrInvalidFrontmatter
	}

	// titleは必須
	title, ok := metaData["title"].(string)
	if !ok || title == "" {
		return nil, fmt.Errorf("%w: titleが必要です", ErrInvalidFrontmatter)
	}

	htmlStr := buf.String()

	// Mermaidコードブロック検知・置換
	hasMermaid := strings.Contains(htmlStr, `class="language-mermaid"`)
	if hasMermaid {
		htmlStr = replaceMermaidBlocks(htmlStr)
	}

	// frontmatterからmermaidフラグも確認
	if mermaidFlag, ok := metaData["mermaid"].(bool); ok && mermaidFlag {
		hasMermaid = true
	}

	post := &model.Post{
		Title:     title,
		Content:   string(content),
		HTML:      htmlStr,
		Tags:      extractTags(metaData),
		Category:  extractString(metaData, "category"),
		Status:    model.PostStatus(extractString(metaData, "status")),
		SourceURL: extractString(metaData, "source_url"),
		Mermaid:   hasMermaid,
		Date:      extractDate(metaData),
		Excerpt:   generateExcerpt(htmlStr),
	}

	return post, nil
}

// replaceMermaidBlocks はMermaidコードブロックをクライアントレンダリング用に置換します。
func replaceMermaidBlocks(htmlStr string) string {
	// <pre><code class="language-mermaid">...</code></pre>
	// → <pre class="mermaid">...</pre>
	result := mermaidCodeBlockRe.ReplaceAllString(htmlStr, `<pre class="mermaid">`)
	// mermaidブロック内の</code></pre>を</pre>に置換
	// 注意: mermaidブロック以外の</code></pre>は変更しない
	// mermaidブロックは<pre class="mermaid">で始まるため、その直後の</code></pre>のみ対象
	parts := strings.Split(result, `<pre class="mermaid">`)
	if len(parts) <= 1 {
		return result
	}
	var sb strings.Builder
	sb.WriteString(parts[0])
	for _, part := range parts[1:] {
		sb.WriteString(`<pre class="mermaid">`)
		// このpart内の最初の</code></pre>のみ</pre>に置換
		part = strings.Replace(part, "</code></pre>", "</pre>", 1)
		sb.WriteString(part)
	}
	return sb.String()
}

// extractTags はmetaDataからタグを抽出します。
func extractTags(metaData map[string]interface{}) []string {
	rawTags, ok := metaData["tags"]
	if !ok {
		return nil
	}

	switch tags := rawTags.(type) {
	case []interface{}:
		result := make([]string, 0, len(tags))
		for _, tag := range tags {
			if s, ok := tag.(string); ok {
				result = append(result, s)
			}
		}
		return result
	case []string:
		return tags
	default:
		return nil
	}
}

// extractString はmetaDataから文字列を取得します。
func extractString(metaData map[string]interface{}, key string) string {
	if v, ok := metaData[key].(string); ok {
		return v
	}
	return ""
}

// extractDate はmetaDataから日付を取得します。
func extractDate(metaData map[string]interface{}) time.Time {
	raw, ok := metaData["date"]
	if !ok {
		return time.Time{}
	}

	switch v := raw.(type) {
	case time.Time:
		return v
	case string:
		// YYYY-MM-DD形式を試行
		t, err := time.Parse("2006-01-02", v)
		if err != nil {
			return time.Time{}
		}
		return t
	default:
		return time.Time{}
	}
}

// htmlTagRe はHTMLタグ除去用の正規表現です。
var htmlTagRe = regexp.MustCompile(`<[^>]*>`)

// generateExcerpt はHTMLからプレーンテキストの抜粋を生成します。
func generateExcerpt(htmlStr string) string {
	// HTMLタグ除去
	text := htmlTagRe.ReplaceAllString(htmlStr, "")
	// 改行をスペースに置換
	text = strings.ReplaceAll(text, "\n", " ")
	// 連続スペースを1つに
	text = strings.Join(strings.Fields(text), " ")
	text = strings.TrimSpace(text)

	runes := []rune(text)
	if len(runes) > maxExcerptLength {
		return string(runes[:maxExcerptLength])
	}
	return text
}
