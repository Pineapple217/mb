package rss

import (
	"encoding/xml"
	"strconv"
	"time"

	"github.com/Pineapple217/mb/config"
	"github.com/Pineapple217/mb/database"
	"github.com/labstack/echo/v4"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel Channel
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Generator   string `xml:"generator"`
	Language    string `xml:"language"`
	Copyright   string `xml:"copyright"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Category    string `xml:"category"`
}

func RSSFeed(c echo.Context) error {
	c.Response().Header().Add("Content-Type", "application/rss+xml")

	q := database.GetQueries()
	posts, err := q.ListPosts(c.Request().Context())
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
		}
		rssPosts = append(rssPosts, p)
	}

	feed := RSS{
		Version: "2.0",
		Channel: Channel{
			Title:       "Micro Blog of " + config.HomepageRights,
			Description: config.HomepageMessage,
			Link:        config.Host + "/index.xml",
			Items:       rssPosts,
			Generator:   "Golang",
			Language:    "en-uk",
			Copyright:   "Copyright " + strconv.Itoa(time.Now().Year()) + ", " + config.HomepageRights,
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
