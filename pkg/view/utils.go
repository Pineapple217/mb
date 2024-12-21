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
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

var (
	htmlFormatter  *html.Formatter
	highlightStyle *chroma.Style
)

func init() {
	htmlFormatter = html.New(html.WithClasses(false), html.TabWidth(4))
	if htmlFormatter == nil {
		panic("couldn't create html formatter")
	}
	styleName := "xcode-dark"
	highlightStyle = styles.Get(styleName)
	if highlightStyle == nil {
		panic(fmt.Sprintf("didn't find style '%s'", styleName))
	}
	builder := highlightStyle.Builder()
	bg := builder.Get(chroma.Background)
	bg.Background = 0
	bg.NoInherit = true

	builder.AddEntry(chroma.Background, bg)
	style, err := builder.Build()
	if err != nil {
		panic(err)
	}
	highlightStyle = style
}

type NullTags struct {
	Valid bool
	Tags  []database.GetAllTagsRow
}

type spotifyEmbed struct {
	ast.Leaf
	cache database.SpotifyCache
}

type youtubeEmbed struct {
	ast.Leaf
	cache database.YoutubeCache
}

var (
	reY             *regexp.Regexp   = regexp.MustCompile(`https?://(?:www\.)?youtu(?:be\.com/watch\?v=|\.be/)([\w\-]+)`)
	reYTID          *regexp.Regexp   = regexp.MustCompile(`(?:youtube\.com\/watch\?v=|youtu\.be\/)([^&?/]+)`)
	renderer        *mdhtml.Renderer = initRender()
	redendererMutex sync.Mutex
)

func embedRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	// context is not needed for embed templates
	// no context removes need for closure and simplifies code
	if s, ok := node.(*spotifyEmbed); ok {
		SpotifyEmbed(s.cache).Render(context.Background(), w)
		return ast.GoToNext, true
	}
	if s, ok := node.(*youtubeEmbed); ok {
		YoutubeEmbed(s.cache).Render(context.Background(), w)
		return ast.GoToNext, true
	}
	if code, ok := node.(*ast.CodeBlock); ok {
		renderCode(w, code, entering)
		return ast.GoToNext, true
	}
	return ast.GoToNext, false
}

func renderCode(w io.Writer, codeBlock *ast.CodeBlock, entering bool) {
	defaultLang := ""
	lang := string(codeBlock.Info)
	htmlHighlight(w, string(codeBlock.Literal), lang, defaultLang)
}

func htmlHighlight(w io.Writer, source, lang, defaultLang string) error {
	if lang == "" {
		lang = defaultLang
	}
	l := lexers.Get(lang)
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return err
	}
	return htmlFormatter.Format(w, highlightStyle, it)
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

	redendererMutex.Lock()
	defer redendererMutex.Unlock()
	return string(markdown.Render(doc, renderer))
}

func initRender() *mdhtml.Renderer {
	htmlFlags := mdhtml.CommonFlags | mdhtml.HrefTargetBlank
	opts := mdhtml.RendererOptions{
		Flags:          htmlFlags,
		RenderNodeHook: embedRenderHook,
	}
	return mdhtml.NewRenderer(opts)

}

func UnixTimeToHTMLDateTime(unixTime int64) string {
	goTime := time.Unix(unixTime, 0).In(config.OutputTimezone)
	formattedTime := goTime.Format("2006-01-02T15:04:05.000Z")
	htmlDateTime := fmt.Sprintf(`<time datetime="%s">%s</time>`, formattedTime, goTime.Format("Mon, 2 Jan 2006 15:04:05 MST"))

	return htmlDateTime
}
