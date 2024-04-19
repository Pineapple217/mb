package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/Pineapple217/mb/pkg/config"
	"github.com/Pineapple217/mb/pkg/database"
	"github.com/Pineapple217/mb/pkg/handler"
	"github.com/Pineapple217/mb/pkg/media"
	"github.com/Pineapple217/mb/pkg/server"
)

const banner = `
• ▌ ▄ ·. ▄▄▄▄· 
·██ ▐███▪▐█ ▀█▪
▐█ ▌▐▌▐█·▐█▀▀█▄
██ ██▌▐█▌██▄▪▐█
▀▀  █▪▀▀▀·▀▀▀▀	v0.6.1
Minimal blog with no JavaScript
https://github.com/Pineapple217/mb/pkg
---------------------------------------------------`

func main() {
	slog.SetDefault(slog.New(slog.Default().Handler()))
	fmt.Println(banner)

	CreateDataDir()
	media.CreateUploadDir()

	fmt.Println("Loading database...")
	q := database.NewQueries("file:" + config.DataDir + "/database.db")
	h := handler.NewHandler(q)

	server := server.NewServer()
	server.RegisterRoutes(h)
	server.ApplyMiddleware(q)

	fmt.Println("Loading configs...")
	config.Load()
	fmt.Println("Loading middlewares...")

	// e.Static("/static", "./static/public")

	server.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	slog.Info("Received an interrupt signal, exiting...")

	server.Stop()
}

func CreateDataDir() {
	if _, err := os.Stat(config.DataDir); os.IsNotExist(err) {
		err := os.Mkdir(config.DataDir, 0755)
		if err != nil {
			slog.Error("Failed to create directory",
				"dir", config.DataDir,
				"error", err,
			)
		}
	}
}
