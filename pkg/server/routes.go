package server

import (
	"embed"
	"time"

	"github.com/Pineapple217/mb/pkg/config"
	"github.com/Pineapple217/mb/pkg/handler"
	"github.com/Pineapple217/mb/pkg/middleware"
	"github.com/labstack/echo/v4"
)

var (
	//go:embed static/public/*
	publicFS embed.FS
)

func (server *Server) RegisterRoutes(hdlr *handler.Handler) {
	e := server.e

	s := e.Group("/static")
	// TODO: post issue, StaticFS not getting cached
	bootTime := time.Now().Add(-2 * time.Hour)
	s.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Add("Last-Modified", bootTime.Local().UTC().Format("Mon, 2 Jan 2006 15:04:05 GMT"))
			return next(c)
		}
	})
	s.StaticFS("/", echo.MustSubFS(publicFS, "static/public"))

	e.GET("/index.xml", hdlr.RSSFeed)

	e.GET("robot.txt", hdlr.RobotTxt)

	//TODO better caching with http headers

	a := e.Group("", middleware.CheckAuth)

	a.Static("/backup", config.BackupDir)
	a.GET("/backup", hdlr.Backups)
	a.POST("/backup", hdlr.CreateBackup)

	e.GET("/auth", hdlr.AuthForm)
	e.POST("/auth", hdlr.Auth)

	e.GET("/post/:xid", hdlr.Post)
	a.GET("/post/:xid/edit", hdlr.EditPostForm)
	a.POST("/post/:xid/edit", hdlr.EditPost)
	a.GET("/post/:xid/delete", hdlr.DeletePostForm)
	a.POST("/post/:xid/delete", hdlr.DeletePost)
	a.POST("/post", hdlr.CreatePost)
	e.GET("/", hdlr.Posts)

	e.GET("/media/t/:name", hdlr.Thumbnail)
	e.GET("/media/:id", hdlr.Mediafile)
	a.GET("/media/:id/edit", hdlr.MediaEditForm)
	a.POST("/media/:id/edit", hdlr.MediaEdit)
	a.GET("/media/:id/delete", hdlr.MediaDeleteForm)
	a.POST("/media/:id/delete", hdlr.MediaDelete)
	a.GET("/media", hdlr.Media)
	a.POST("/media", hdlr.MediaUpload)

	e.Static("/m", config.UploadDir)
}
