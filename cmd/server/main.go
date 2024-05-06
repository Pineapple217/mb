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
)

// TODO: banner does not get printed first
const banner = `
• ▌ ▄ ·. ▄▄▄▄· 
·██ ▐███▪▐█ ▀█▪
▐█ ▌▐▌▐█·▐█▀▀█▄
██ ██▌▐█▌██▄▪▐█
▀▀  █▪▀▀▀·▀▀▀▀	v0.7.1
Minimal blog with no JavaScript
https://github.com/Pineapple217/mb/pkg
-----------------------------------------------------------------------------`

func main() {
	slog.SetDefault(slog.New(slog.Default().Handler()))
	fmt.Println(banner)

	CreateDataDir()
	media.CreateUploadDir()

	config.Load()

	q := database.NewQueries("file:" + config.DataDir + "/database.db")
	h := handler.NewHandler(q)

	server := server.NewServer()
	server.RegisterRoutes(h)
	server.ApplyMiddleware(q)

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
