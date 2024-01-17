package view

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/Pineapple217/mb/database"
	"github.com/Pineapple217/mb/embed"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func renderSpotifyEmbed(ctx context.Context, w io.Writer, l *ast.Link, entering bool) {
	if entering {
		re := regexp.MustCompile(`/track/(\w+)`)
		id := re.FindStringSubmatch(string(l.Destination))[1]
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
		re := regexp.MustCompile(`(?:youtube\.com\/watch\?v=|youtu\.be\/)([^&?/]+)`)
		id := re.FindStringSubmatch(string(l.Destination))[1]
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
	return func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
		if link, ok := node.(*ast.Link); ok {
			reS := regexp.MustCompile(`https?://open\.spotify\.com/track/(\S+)`)
			if reS.MatchString(string(link.Destination)) {
				renderSpotifyEmbed(ctx, w, link, entering)
				return ast.GoToNext, true
			}
			reY := regexp.MustCompile(`https?://(?:www\.)?youtu(?:be\.com/watch\?v=)|(?:\.be/)(\S+)`)
			if reY.MatchString(string(link.Destination)) {
				renderYoutubeEmbed(ctx, w, link, entering)
				return ast.GoToNext, true
			}
		}
		return ast.GoToNext, false
	}
}

func MdToHTML(ctx context.Context, md string) string {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock | parser.FencedCode
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse([]byte(md))

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{
		Flags:          htmlFlags,
		RenderNodeHook: makeEmbedRenderHook(ctx),
	}
	renderer := html.NewRenderer(opts)

	// TODO: syntax highlighter with github.com/alecthomas/chroma
	// https://blog.kowalczyk.info/article/cxn3/advanced-markdown-processing-in-go.html
	return string(markdown.Render(doc, renderer))
}

func UnixTimeToHTMLDateTime(unixTime int64) string {
	// TODO: make it an env
	loc, err := time.LoadLocation("Europe/Brussels")
	if err != nil {
		loc = time.UTC
	}
	goTime := time.Unix(unixTime, 0).In(loc)
	formattedTime := goTime.Format("2006-01-02T15:04:05.000Z")
	htmlDateTime := fmt.Sprintf(`<time datetime="%s">%s</time>`, formattedTime, goTime.Format("Mon, 2 Jan 2006 15:04:05 MST"))

	return htmlDateTime
}
