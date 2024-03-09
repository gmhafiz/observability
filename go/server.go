package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

type Server struct {
	Cfg

	Mux *http.ServeMux
	DB  *sql.DB
}

func (s *Server) init() {
	s.initLog()
	s.newOpenTelemetry()
	s.initDatabase()
}

func (s *Server) initLog() {
	slog.SetDefault(slog.New(NewTraceHandler(
		os.Stdout,
		&slog.HandlerOptions{},
	)))
}

func (s *Server) initDatabase() {
	ctx := context.Background()

	db, err := NewDB(ctx, s.Cfg.Database)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	s.DB = db
}

func (s *Server) newOpenTelemetry() {
	ctx := context.Background()
	_ = SetupOTLPExporter(ctx, s.Cfg.OpenTelemetry)
}

func (s *Server) Readiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	if err := s.DB.PingContext(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"status":500}`))
		return
	}

	_, _ = w.Write([]byte(`{"status":200}`))
}

func New() *Server {
	cfg := Config()
	mux := http.NewServeMux()

	return &Server{
		Cfg: cfg,
		Mux: mux,
	}
}
