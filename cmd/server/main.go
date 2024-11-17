package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/Pineapple217/mb/pkg/config"
	"github.com/Pineapple217/mb/pkg/database"
	"github.com/Pineapple217/mb/pkg/handler"
	"github.com/Pineapple217/mb/pkg/media"
	"github.com/Pineapple217/mb/pkg/scheduler"
	"github.com/Pineapple217/mb/pkg/server"
	"github.com/Pineapple217/mb/pkg/static"
)

const banner = `
• ▌ ▄ ·. ▄▄▄▄· 
·██ ▐███▪▐█ ▀█▪
▐█ ▌▐▌▐█·▐█▀▀█▄
██ ██▌▐█▌██▄▪▐█
▀▀  █▪▀▀▀·▀▀▀▀	v0.8.0
Minimal blog with no JavaScript
https://github.com/Pineapple217/mb
-----------------------------------------------------------------------------`

func main() {
	slog.SetDefault(slog.New(slog.Default().Handler()))
	fmt.Println(banner)
	os.Stdout.Sync()

	CreateDataDir()
	media.CreateUploadDir()

	config.Load()

	rr := static.HashPublicFS()

	q := database.NewQueries("file:" + config.DataDir + "/database.db")
	h := handler.NewHandler(q)

	server := server.NewServer()
	server.RegisterRoutes(h)
	server.ApplyMiddleware(q, rr)

	server.Start()

	s := scheduler.NewScheduler()
	s.Schedule(time.Hour*24, func() {
		scheduler.CleanCache(q)
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	slog.Info("Received an interrupt signal, exiting...")

	s.Stop()
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
