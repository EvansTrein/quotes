package server

import (
	"context"
	"log/slog"
	"net/http"
	"quotes/config"
	"quotes/internal/controller"
	"quotes/internal/service"
	"quotes/pkg/middleware"
	"time"
)

const (
	gracefulShutdownTimer = time.Second * 5
)

type Server struct {
	log  *slog.Logger
	http *http.Server

	quote *controller.QuotesController
}

type ServerDeps struct {
	*slog.Logger
}

func NewServer(deps *ServerDeps) *Server {
	addr := config.GetAddress() + ":" + config.GetPort()
	engine := http.NewServeMux()

	quoteServ := service.NewQuoteService(&service.QuoteServiceDeps{
		Logger: deps.Logger,
	})

	base := controller.NewBaseController(&controller.BaseControllerDeps{
		Logger: deps.Logger,
	})

	quote := controller.NewQuotesController(&controller.QuotesControllerDeps{
		Router:         engine,
		BaseController: base,
		IQuoteService:  quoteServ,
	})

	httpLog := middleware.NewMiddlewareLogging(&middleware.MiddlewareLoggingDeps{
		Logger: deps.Logger,
	})

	server := &http.Server{
		Addr: addr,
		Handler: middleware.ChainMiddleware(
			httpLog.HandlersLog(),
		)(engine),
	}

	return &Server{
		log:   deps.Logger,
		http:  server,
		quote: quote,
	}
}

func (s *Server) Start() error {

	s.log.Info("HTTP server: successfully started", "address", s.http.Addr)
	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	s.log.Debug("HTTP server: stop started")

	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimer)
	defer cancel()

	if err := s.http.Shutdown(ctx); err != nil {
		s.log.Error("Server shutdown failed", "error", err)
		return err
	}

	s.http = nil
	s.quote = nil
	s.log.Info("HTTP server: stop successful")

	return nil
}
