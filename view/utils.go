package view

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"sync"
	"time"

	"github.com/Pineapple217/mb/config"
	"github.com/Pineapple217/mb/database"
	"github.com/Pineapple217/mb/embed"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var (
	reS             *regexp.Regexp = regexp.MustCompile(`https?://open\.spotify\.com/track/(\S+)`)
	reY             *regexp.Regexp = regexp.MustCompile(`https?://(?:www\.)?youtu(?:be\.com/watch\?v=)|(?:\.be/)(\S+)`)
	reYTID          *regexp.Regexp = regexp.MustCompile(`(?:youtube\.com\/watch\?v=|youtu\.be\/)([^&?/]+)`)
	reSID           *regexp.Regexp = regexp.MustCompile(`/track/(\w+)`)
	renderer        *html.Renderer = initRender()
	redendererMutex sync.Mutex
)

func renderSpotifyEmbed(ctx context.Context, w io.Writer, l *ast.Link, entering bool) {
	if entering {
		id := reSID.FindStringSubmatch(string(l.Destination))[1]
		queries := database.GetQueries()
		sc, err := queries.GetSpotifyCache(ctx, id)
		if err != nil {
			sc = embed.SpotifyScrape(ctx, string(l.Destination))
		}
		SpotifyEmbed(sc).Render(ctx, w)
		// setting the content to nil so the OG url wil show
		l.Children[0].AsLeaf().Literal = nil
	} else {
		// prevents string that are in the same p form being exleded
		// TODO: modify node tree to remove this fix
		// https://blog.kowalczyk.info/article/cxn3/advanced-markdown-processing-in-go.html
		io.WriteString(w, "<p/><p>")
	}

}

func renderYoutubeEmbed(ctx context.Context, w io.Writer, l *ast.Link, entering bool) {
	if entering {
		id := reYTID.FindStringSubmatch(string(l.Destination))[1]
		queries := database.GetQueries()
		ytc, err := queries.GetYoutubeCache(ctx, id)
		if err != nil {
			ytc = embed.YoutubeScrape(ctx, id)
		}
		YoutubeEmbed(ytc).Render(ctx, w)
		// setting the content to nil so the OG url wil show
		l.Children[0].AsLeaf().Literal = nil
	} else {
		// prevents string that are in the same p form being exleded
		// TODO: modify node tree to remove this fix
		// https://blog.kowalczyk.info/article/cxn3/advanced-markdown-processing-in-go.html
		io.WriteString(w, "<p/><p>")
	}
}

func makeEmbedRenderHook(ctx context.Context) html.RenderNodeFunc {
	// TODO: opt out of embed with [link text](url)
	return func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
		if link, ok := node.(*ast.Link); ok {
			if reS.MatchString(string(link.Destination)) {
				renderSpotifyEmbed(ctx, w, link, entering)
				return ast.GoToNext, true
			}
			if reY.MatchString(string(link.Destination)) {
				renderYoutubeEmbed(ctx, w, link, entering)
				return ast.GoToNext, true
			}
		}
		// TODO: add custom video and audio nodes
		return ast.GoToNext, false
	}
}

func MdToHTML(ctx context.Context, md string) string {
	// create markdown parser with extensions
	// TODO: make a gobal parser 1 time
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.FencedCode
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(md))

	// TODO: mutex reduces speed by 20%, add renderer pool speed up
	redendererMutex.Lock()
	defer redendererMutex.Unlock()
	renderer.Opts.RenderNodeHook = makeEmbedRenderHook(ctx)

	// TODO: syntax highlighter with github.com/alecthomas/chroma
	// https://blog.kowalczyk.info/article/cxn3/advanced-markdown-processing-in-go.html
	return string(markdown.Render(doc, renderer))
}

func initRender() *html.Renderer {
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{
		Flags: htmlFlags,
	}
	return html.NewRenderer(opts)

}

func UnixTimeToHTMLDateTime(unixTime int64) string {
	goTime := time.Unix(unixTime, 0).In(config.OutputTimezone)
	formattedTime := goTime.Format("2006-01-02T15:04:05.000Z")
	htmlDateTime := fmt.Sprintf(`<time datetime="%s">%s</time>`, formattedTime, goTime.Format("Mon, 2 Jan 2006 15:04:05 MST"))

	return htmlDateTime
}
