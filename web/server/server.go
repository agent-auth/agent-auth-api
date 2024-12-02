package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/agent-auth/agent-auth-api/web/router"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Server provides an http.Server.
type server struct {
	svr               *http.Server
	logger            *zap.Logger
	startTimeStampUTC time.Time
}

// NewServer creates and configures an APIServer serving all application routes.
func NewServer() Server {
	var addr string
	port := viper.GetString("host.port")
	apiHandler := router.NewRouter().Router(viper.GetBool("host.enable_cors"))

	// allow port to be set as localhost:8001 in env during development to avoid "accept incoming network connection" request on restarts
	if strings.Contains(port, ":") {
		addr = port
	} else {
		addr = ":" + port
	}

	srv := http.Server{
		Addr:    addr,
		Handler: apiHandler,
	}

	logger, _ := zap.NewProduction()

	return &server{
		svr:    &srv,
		logger: logger,
	}
}

// Start runs ListenAndServe on the http.Server with graceful shutdown.
func (s *server) Start() {
	s.logger.Info("starting server",
		zap.String("address", s.svr.Addr),
	)
	go func() {
		if err := s.svr.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Fatal("failed to start server",
				zap.Error(err))
		}
	}()

	s.startTimeStampUTC = time.Now().UTC()
	viper.Set("service_started_timestamp_utc", time.Now().UTC())

	s.logger.Info("server listening",
		zap.String("address", s.svr.Addr))

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit

	s.logger.Info("shutting down server",

		zap.String("reason", sig.String()))
	if err := s.svr.Shutdown(context.Background()); err != nil {
		s.logger.Fatal("failed to stop server",

			zap.Error(err))
	}

	s.logger.Info("server gracefully stopped")
}

func (s *server) StartTimeStampUTC() time.Time {
	return s.startTimeStampUTC
}
