package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/WilliamDeBruin/nps_alerts/src/config"
	"github.com/WilliamDeBruin/nps_alerts/src/nps"
	"github.com/WilliamDeBruin/nps_alerts/src/twilio"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/leosunmo/zapchi"
	"go.uber.org/zap"

	"github.com/pkg/errors"
)

type Server struct {
	twilioClient *twilio.Client
	npsClient    nps.Client
	httpServer   *http.Server
	port         string
	logger       *zap.Logger
}

func NewServer(
	cfg *config.Configuration,
	logger *zap.Logger,
) (*Server, error) {

	twilioClient, err := twilio.NewClient(cfg.TwilioFromNumber)
	if err != nil {
		log.Fatalf("error initializing twilio client: %s", err)
	}

	npsClient, err := nps.NewClient()
	if err != nil {
		log.Fatalf("error initializing nps client: %s", err)
	}

	s := &Server{
		twilioClient: twilioClient,
		npsClient:    npsClient,
		port:         cfg.Port,
		logger:       logger,
	}

	return s, nil
}

func (s *Server) Serve() {
	defer s.Close()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		panic(fmt.Sprintf("unable to serve: %s", err))
	}

	s.listen(listener)
}

func (s *Server) listen(listener net.Listener) {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(zapchi.Logger(s.logger, "router"))

	router.Get("/health", s.HealthHandler)
	router.Post("/incoming-sms", s.IncomingSmsHandler)

	// re-discover what port we are on. If config was to port :0, this will allow us to know what port we bound to
	port := listener.Addr().(*net.TCPAddr).Port
	s.port = fmt.Sprintf("%d", port)
	s.logger.Info(fmt.Sprintf("listening on %d", port))

	s.httpServer = &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: router}
	// TODO add timeouts to config
	s.httpServer.WriteTimeout = 1 * time.Minute
	s.httpServer.ReadTimeout = 1 * time.Minute

	if err := s.httpServer.Serve(listener); err != nil {
		if err != http.ErrServerClosed {
			s.logger.Error(fmt.Sprintf("server crash: %v", err))
			os.Exit(1)
		}
	}
}

// Close closes all db connections or any other clean up
func (srv *Server) Close() error {
	// potentially doing many things that could error. Keep all errors and return at the end.
	var errs error
	// close socket to stop new requests from coming in
	err := srv.httpServer.Close()
	if err != nil {
		errs = errors.Wrap(err, "error closing http server")
	}

	return errs
}
