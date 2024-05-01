package server

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
)

var (
	listen = flag.String("listen", "127.0.0.1", "Where to listen, 0.0.0.0 is needed for docker")
	port   = flag.String("port", ":3000", "Port to listen on")
)

type Server struct {
	e *echo.Echo
}

func NewServer() *Server {
	slog.Info("Creating server")
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Debug = true
	NewServer := &Server{
		e: e,
	}

	return NewServer
}

// Starts the server in a new routine
func (s *Server) Start() {
	flag.Parse()
	slog.Info("Starting server")
	go func() {
		if err := s.e.Start(*listen + *port); err != nil && err != http.ErrServerClosed {
			slog.Error("Shutting down the server", "error", err.Error())
		}
	}()
	slog.Info("Server started", "bind", *listen, "port", *port)
}

// Tries to the stops the server gracefully
func (s *Server) Stop() {
	slog.Info("Stopping server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.e.Shutdown(ctx); err != nil {
		slog.Error(err.Error())
	}
}
