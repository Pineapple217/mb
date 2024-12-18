package server

import (
	"log/slog"
	"time"

	"github.com/Pineapple217/mb/pkg/database"
	"github.com/Pineapple217/mb/pkg/handler"
	"github.com/Pineapple217/mb/pkg/middleware"
	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
)

func (s *Server) ApplyMiddleware(q *database.Queries, reRoutes map[string]string) {
	slog.Info("Applying middlewares")
	s.e.Pre(echoMw.Rewrite(reRoutes))
	s.e.Use(echoMw.RequestLoggerWithConfig(echoMw.RequestLoggerConfig{
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
	s.e.Use(middleware.Path)

	s.e.Use(echoMw.RateLimiterWithConfig(echoMw.RateLimiterConfig{
		Store: echoMw.NewRateLimiterMemoryStoreWithConfig(
			echoMw.RateLimiterMemoryStoreConfig{Rate: 10, Burst: 30, ExpiresIn: 3 * time.Minute},
		),
	}))

	s.e.Use(echoMw.GzipWithConfig(echoMw.GzipConfig{
		Level: 5,
	}))

	echo.NotFoundHandler = handler.NotFound
	s.e.Use(middleware.Auth)
	s.e.Use(func(next echo.HandlerFunc) echo.HandlerFunc { return middleware.Stats(next, q) })
}
