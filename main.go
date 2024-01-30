package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Pineapple217/mb/config"
	"github.com/Pineapple217/mb/database"
	"github.com/Pineapple217/mb/handler"
	"github.com/Pineapple217/mb/middleware"

	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
)

var (
	//go:embed static/public/*
	publicFS embed.FS
	listen   = flag.String("listen", "127.0.0.1", "Where to listen, 0.0.0.0 is needed for docker")
	port     = flag.String("port", ":3000", "Port to listen on")
)

const banner = `
• ▌ ▄ ·. ▄▄▄▄· 
·██ ▐███▪▐█ ▀█▪
▐█ ▌▐▌▐█·▐█▀▀█▄
██ ██▌▐█▌██▄▪▐█
▀▀  █▪▀▀▀·▀▀▀▀	v0.1.0
Minimal blog with no JavaScript
https://github.com/Pineapple217/mb
---------------------------------------------------`

func main() {
	slog.SetDefault(slog.New(slog.Default().Handler()))
	flag.Parse()
	e := echo.New()
	e.HideBanner = true
	fmt.Println(banner)

	fmt.Println("Loading configs...")
	config.Load()
	fmt.Println("Loading middlewares...")
	e.Use(echoMw.RequestLoggerWithConfig(echoMw.RequestLoggerConfig{
		LogStatus:  true,
		LogURI:     true,
		LogMethod:  true,
		LogLatency: true,
		LogValuesFunc: func(c echo.Context, v echoMw.RequestLoggerValues) error {
			slog.Info("request",
				"method", v.Method,
				"status", v.Status,
				"latency", v.Latency,
				"path", v.URI,
			)
			return nil

		},
	}))

	echo.NotFoundHandler = handler.NotFound
	e.Use(middleware.Stats)
	e.Use(middleware.Auth)

	// TODO: post issue, StaticFS not getting cached
	s := e.Group("/static")
	bootTime := time.Now().Add(-2 * time.Hour)
	s.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Add("Last-Modified", bootTime.Local().UTC().Format("Mon, 2 Jan 2006 15:04:05 GMT"))
			return next(c)
		}
	})

	fmt.Println("Loading database...")
	database.Init("file:database.db")

	// e.Static("/static", "./static/public")
	s.StaticFS("/", echo.MustSubFS(publicFS, "static/public"))

	//TODO RSS

	//TODO better caching with http headers

	//TODO backup

	e.GET("/auth", handler.AuthForm)
	e.POST("/auth", handler.Auth)

	e.GET("/post/:xid", handler.Post)
	e.GET("/post/:xid/edit", middleware.CheckAuth(handler.EditPostForm))
	e.POST("/post/:xid/edit", middleware.CheckAuth(handler.EditPost))
	e.GET("/post/:xid/delete", middleware.CheckAuth(handler.DeletePostForm))
	e.POST("/post/:xid/delete", middleware.CheckAuth(handler.DeletePost))
	e.POST("/post", middleware.CheckAuth(handler.CreatePost))
	e.GET("/", handler.Posts)

	go func() {
		if err := e.Start(*listen + *port); err != nil && err != http.ErrServerClosed {
			slog.Error("Shutting down the server", "error", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		slog.Error(err.Error())
	}
}
