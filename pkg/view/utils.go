package view

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Pineapple217/mb/pkg/config"
	"github.com/Pineapple217/mb/pkg/database"
	"github.com/Pineapple217/mb/pkg/embed"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type NullTags struct {
	Valid bool
	Tags  []database.GetAllTagsRow
}

type spotifyEmbed struct {
	ast.Leaf
	cache database.SpotifyCache
}

type navidromeEmbed struct {
	ast.Leaf
	cache database.NavidromeCache
}

type youtubeEmbed struct {
	ast.Leaf
	cache database.YoutubeCache
}

var (
	reY             *regexp.Regexp = regexp.MustCompile(`https?://(?:www\.)?youtu(?:be\.com/watch\?v=|\.be/)([\w\-]+)`)
	reYTID          *regexp.Regexp = regexp.MustCompile(`(?:youtube\.com\/watch\?v=|youtu\.be\/)([^&?/]+)`)
	renderer        *html.Renderer = initRender()
	redendererMutex sync.Mutex
)

func embedRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	// context is not needed for embed templates
	// no context removes need for closure and simplifies code
	if s, ok := node.(*spotifyEmbed); ok {
		SpotifyEmbed(s.cache).Render(context.Background(), w)
		return ast.GoToNext, true
	}
	if s, ok := node.(*navidromeEmbed); ok {
		NavidromeEmbed(s.cache).Render(context.Background(), w)
		return ast.GoToNext, true
	}
	if s, ok := node.(*youtubeEmbed); ok {
		YoutubeEmbed(s.cache).Render(context.Background(), w)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func makeParserHook(ctx context.Context, q *database.Queries) parser.BlockFunc {
	return func(data []byte) (ast.Node, []byte, int) {
		if bytes.HasPrefix(data, []byte(embed.SpotifyUrlPrefix)) {
			i := bytes.IndexByte(data, '\n')
			var d string
			if i == -1 {
				d = string(data)
			} else {
				d = string(data[:i])
			}

			id, _ := strings.CutPrefix(d, embed.SpotifyUrlPrefix)
			id = strings.Split(id, "?si=")[0]

			c, err := q.GetSpotifyCache(ctx, id)
			if errors.Is(err, sql.ErrNoRows) {
				c, err = embed.SpotifyScrape(ctx, q, id)
				if err != nil {
					slog.Error("Failed to scrape track", "id", id, "err", err)
					return nil, nil, 0
				}
			} else if err != nil {
				slog.Error("Failed to fetch track cache", "id", id, "err", err)
				return nil, nil, 0
			}
			node := spotifyEmbed{cache: c}
			return &node, nil, len(d)
		}

		if config.NavidromePrefix != "" && bytes.HasPrefix(data, []byte(config.NavidromePrefix)) {
			i := bytes.IndexByte(data, '\n')
			var d string
			if i == -1 {
				d = string(data)
			} else {
				d = string(data[:i])
			}
			id, _ := strings.CutPrefix(d, config.NavidromePrefix)

			c, err := q.GetNavidromeCache(ctx, id)
			if errors.Is(err, sql.ErrNoRows) {
				c, err = embed.NavidromeScrape(ctx, q, id)
				if err != nil {
					slog.Error("Failed to scrape navidrome track", "id", id, "err", err)
					return nil, nil, 0
				}
			} else if err != nil {
				slog.Error("Failed to fetch navidrome track cache", "id", id, "err", err)
				return nil, nil, 0
			}
			node := navidromeEmbed{cache: c}
			return &node, nil, len(d)
		}

		// TODO: cleanup regex
		if bytes.HasPrefix(data, []byte("http")) && reY.Match(data) {
			i := bytes.IndexByte(data, '\n')
			var d string
			if i == -1 {
				d = string(data)
			} else {
				d = string(data[:i])
			}
			id := reYTID.FindStringSubmatch(d)[1]
			c, err := q.GetYoutubeCache(ctx, id)
			if errors.Is(err, sql.ErrNoRows) {
				c, err = embed.YoutubeScrape(ctx, q, id)
				if err != nil {
					slog.Error("Failed to scrape video", "id", id, "err", err)
					return nil, nil, 0
				}
			} else if err != nil {
				slog.Error("Failed to fetch video cache", "id", id, "err", err)
				return nil, nil, 0
			}
			node := youtubeEmbed{cache: c}
			return &node, nil, len(d)
		}
		return nil, nil, 0
	}
}

func MdToHTML(ctx context.Context, q *database.Queries, md string) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.FencedCode
	p := parser.NewWithExtensions(extensions)
	p.Opts.ParserHook = makeParserHook(ctx, q)

	doc := p.Parse([]byte(md))

	// TODO: mutex reduces speed by 20%, add renderer pool speed up
	// TODO: syntax highlighter with github.com/alecthomas/chroma
	// https://blog.kowalczyk.info/article/cxn3/advanced-markdown-processing-in-go.html

	redendererMutex.Lock()
	defer redendererMutex.Unlock()
	return strings.TrimSpace(string(markdown.Render(doc, renderer)))
}

func initRender() *html.Renderer {
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{
		Flags:          htmlFlags,
		RenderNodeHook: embedRenderHook,
	}
	return html.NewRenderer(opts)

}

func UnixTimeToHTMLDateTime(unixTime int64) string {
	goTime := time.Unix(unixTime, 0).In(config.OutputTimezone)
	formattedTime := goTime.Format("2006-01-02T15:04:05.000Z")
	htmlDateTime := fmt.Sprintf(`<time datetime="%s">%s</time>`, formattedTime, goTime.Format("Mon, 2 Jan 2006 15:04:05 MST"))

	return htmlDateTime
}
