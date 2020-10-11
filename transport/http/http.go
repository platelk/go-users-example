package http

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/go-chi/chi"

	"go-users-example/infra/logger"
)

// Config hold configuration for http server
type Config struct {
	Addr string `env:"HTTP_ADDR" env-default:"0.0.0.0:8080"`
}

// Builder will construct the all http server, setup middleware correctly, etc.
type Builder struct {
	c      Config
	log    logger.Logger
	router chi.Router
}

// NewBuilder will initialise Builder
func NewBuilder(log logger.Logger, c Config) *Builder {
	return &Builder{c: c, log: log, router: chi.NewRouter()}
}

// Build will construct the final Server
func (b *Builder) Build() *Server {
	return &Server{
		log:    b.log.With().Str("component", "http_server").Logger(),
		server: &http.Server{Addr: b.c.Addr, Handler: b.router},
	}
}

// Server is the configured http server to run
type Server struct {
	log    logger.Logger
	server *http.Server
}

// Run will block to accept http request and interrupt on os.Interrupt signal
func (s *Server) Run() error {
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := s.server.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := s.server.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
		return err
	}
	<-idleConnsClosed
	return nil
}
