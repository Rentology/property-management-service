package server

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"property-managment-service/internal/config"
	"property-managment-service/lib/sl"
	"strconv"
	"syscall"
	"time"
)

const (
	ctxTimeout = 5
	PprofPort  = "5556"
)

type Server struct {
	echo *echo.Echo
	cfg  *config.Config
	db   *sqlx.DB
	log  *slog.Logger
}

func NewServer(cfg *config.Config, db *sqlx.DB, log *slog.Logger) *Server {
	return &Server{echo: echo.New(), cfg: cfg, db: db, log: log}

}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:         ":" + strconv.Itoa(s.cfg.Server.Port),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}

	go func() {
		sl.Infof(s.log, "Server is listening on PORT: %d", s.cfg.Server.Port)
		if err := s.echo.StartServer(server); err != nil {
			s.log.Error("Error starting Server: ", err)
			os.Exit(1)
		}
	}()

	go func() {
		sl.Infof(s.log, "Starting Debug Server on PORT: %s", PprofPort)
		if err := http.ListenAndServe(":"+PprofPort, http.DefaultServeMux); err != nil {
			s.log.Error("Error PPROF ListenAndServe: %s", err)
			os.Exit(1)
		}
	}()

	if err := s.MapHandlers(s.echo); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	s.log.Info("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}
