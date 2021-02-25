package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/google/uuid"
	"github.com/jobbox-tech/recruiter-api/api"
	"github.com/jobbox-tech/recruiter-api/logging"
	"github.com/spf13/viper"
)

// Server provides an http.Server.
type server struct {
	svr    *http.Server
	logger logging.Logger
	txID   string
}

// New creates and configures an APIServer serving all application routes.
func New() Server {
	txID := uuid.New().String()

	log.Println("configuring server...")
	apiHandler := api.New(viper.GetBool("host.enable_cors"))

	var addr string
	port := viper.GetString("host.port")

	// allow port to be set as localhost:3000 in env during development to avoid "accept incoming network connection" request on restarts
	if strings.Contains(port, ":") {
		addr = port
	} else {
		addr = ":" + port
	}

	srv := http.Server{
		Addr:    addr,
		Handler: apiHandler,
	}

	return &server{
		svr:    &srv,
		logger: logging.NewLogger(),
		txID:   txID,
	}
}

// Start runs ListenAndServe on the http.Server with graceful shutdown.
func (s *server) Start() {
	s.logger.Info(s.txID).Infof("starting server at %v", s.svr.Addr)
	go func() {
		if err := s.svr.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Fatal(s.txID, FailedToStartServer).Errorf("Failed to start server with error %v", err)
		}
	}()

	s.logger.Info(s.txID).Infof("Server listening on %s", s.svr.Addr)
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit

	s.logger.Info(s.txID).Infof("Shutting down server... Reason:%v", sig)
	if err := s.svr.Shutdown(context.Background()); err != nil {
		s.logger.Fatal(s.txID, FailedToStopServer).Errorf("Failed to stop server with error %v", err)
	}

	s.logger.Info(s.txID).Info("Server gracefully stopped")
}
