package handler

import (
	"encoding/xml"
	"strconv"
	"time"

	"github.com/Pineapple217/mb/pkg/config"
	"github.com/labstack/echo/v4"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Xmlns   string   `xml:"xmlns:atom,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Link        string   `xml:"link"`
	Generator   string   `xml:"generator"`
	Language    string   `xml:"language"`
	Copyright   string   `xml:"copyright"`
	AtomLink    AtomLink `xml:"atom:link"`
	Items       []Item   `xml:"item"`
}

type AtomLink struct {
	XMLName xml.Name `xml:"atom:link"`
	Href    string   `xml:"href,attr"`
	Rel     string   `xml:"rel,attr"`
	Type    string   `xml:"type,attr"`
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Category    string `xml:"category"`
	Guid        string `xml:"guid"`
}

func (h *Handler) RSSFeed(c echo.Context) error {
	c.Response().Header().Add("Content-Type", "application/rss+xml")

	posts, err := h.Q.ListPublicPosts(c.Request().Context())
	if err != nil {
		return err
	}
	rssPosts := make([]Item, 0)
	for _, post := range posts {
		unixTime := strconv.FormatInt(post.CreatedAt, 10)
		p := Item{
			Title:       unixTime,
			Link:        config.Host + "/post/" + unixTime,
			PubDate:     time.Unix(post.CreatedAt, 0).Format(time.RFC1123Z),
			Category:    post.Tags.String,
			Description: truncateString(post.Content),
			Guid:        config.Host + "/post/" + unixTime,
		}
		rssPosts = append(rssPosts, p)
	}

	feed := RSS{
		Version: "2.0",
		Xmlns:   "http://www.w3.org/2005/Atom",
		Channel: Channel{
			Title:       "Micro Blog of " + config.HomepageRights,
			Description: config.HomepageMessage,
			Link:        config.Host,
			Items:       rssPosts,
			Generator:   "Golang",
			Language:    "en-uk",
			Copyright:   "Copyright " + strconv.Itoa(time.Now().Year()) + ", " + config.HomepageRights,
			AtomLink:    AtomLink{Href: config.Host + "/index.xml", Rel: "self", Type: "application/rss+xml"},
		},
	}
	xmlData, err := xml.MarshalIndent(feed, "", "    ")
	if err != nil {
		return err
	}
	c.Response().Write([]byte(xml.Header))
	c.Response().Write(xmlData)
	return nil
}

const descriptionMaxLength = 500

func truncateString(s string) string {
	if len(s) <= descriptionMaxLength {
		return s
	}
	return s[:descriptionMaxLength] + "..."
}
