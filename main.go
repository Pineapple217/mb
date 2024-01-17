package main

import (
	"context"
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Pineapple217/mb/auth"
	"github.com/Pineapple217/mb/database"
	"github.com/Pineapple217/mb/handler"
	"github.com/Pineapple217/mb/middleware"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
)

//go:embed static/public/*
var publicFS embed.FS

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
	e := echo.New()
	e.HideBanner = true
	fmt.Println(banner)

	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found")
	}
	password, isSet := os.LookupEnv("AUTH_PASSWORD")
	if !isSet {
		fmt.Println("AUHT_PASSWORD is not set. Using random password")
		fmt.Printf("Auth password: %s\n", auth.SecretPassword)
	} else {
		fmt.Println("Auth password succesfully set")
		auth.SecretPassword = password
	}

	fmt.Println("Loading middlewares...")
	e.Use(echoMw.RequestLoggerWithConfig(echoMw.RequestLoggerConfig{
		LogStatus: true,
		LogURI:    true,
		LogMethod: true,
		LogValuesFunc: func(c echo.Context, v echoMw.RequestLoggerValues) error {
			slog.Info("request",
				"method", v.Method,
				"status", v.Status,
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
		if err := e.Start("127.0.0.1:3000"); err != nil && err != http.ErrServerClosed {
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
