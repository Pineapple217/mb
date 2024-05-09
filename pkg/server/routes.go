package server

import (
	"log/slog"

	"github.com/Pineapple217/mb/pkg/config"
	"github.com/Pineapple217/mb/pkg/handler"
	"github.com/Pineapple217/mb/pkg/middleware"
	"github.com/Pineapple217/mb/pkg/static"
	"github.com/labstack/echo/v4"
)

func (server *Server) RegisterRoutes(hdlr *handler.Handler) {
	slog.Info("Registering routes")
	e := server.e

	s := e.Group("/static")

	s.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Add("Cache-Control", "public, max-age=31536000, immutable")
			return next(c)
		}
	})
	s.StaticFS("/", echo.MustSubFS(static.PublicFS, "public"))

	e.GET("/index.xml", hdlr.RSSFeed)

	e.GET("robot.txt", hdlr.RobotTxt)
	e.GET("/site.webmanifest", hdlr.Manifest)

	//TODO better caching with http headers

	a := e.Group("", middleware.CheckAuth)

	a.Static("/backup", config.BackupDir)
	a.GET("/backup", hdlr.Backups)
	a.POST("/backup", hdlr.CreateBackup)

	e.GET("/auth", hdlr.AuthForm)
	e.POST("/auth", hdlr.Auth)

	e.GET("/post/:xid", hdlr.Post)
	e.GET("post/latest", hdlr.PostLatest)
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
