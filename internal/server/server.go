package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ashershnyov/gophermart-loyalty-program/internal/loyalty"
	"github.com/ashershnyov/gophermart-loyalty-program/internal/middleware"
	"github.com/ashershnyov/gophermart-loyalty-program/internal/server/config"
	"github.com/ashershnyov/gophermart-loyalty-program/pkg/db"
	"github.com/ashershnyov/gophermart-loyalty-program/pkg/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/pressly/goose"
)

type server struct {
	router chi.Router
	config *config.Config
	db     db.DB
}

// ListenAndServe launches listening loop on the address provided in the config.
func (s *server) ListenAndServe() {
	if err := http.ListenAndServe(s.config.Address, s.router); err != nil {
		log.Fatal(err)
	}
}

// New creates a new server.
func New(opts ...config.Option) (*server, error) {
	cfg := config.New(opts...)

	logHandler := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(logHandler)

	router := chi.NewRouter()
	router.Use(
		middleware.Logging(logger),
		middleware.Gzip(),
	)

	if cfg.DBAddress == "" {
		return nil, errors.New("SQL DB address can not be empty")
	}

	var (
		pg  *db.Postgres
		err error
	)

	pg, err = db.NewPostgres(context.Background(), cfg.DBAddress, cfg.MaxRetries)
	if err != nil {
		return nil, fmt.Errorf("an error occured when connecting to the DB: %w", err)
	}

	return &server{
		router: router,
		config: cfg,
		db:     pg,
	}, nil
}

// Run launches the server's loop.
func (s *server) Run() error {
	jwtProvider, err := jwt.NewGenerator()
	if err != nil {
		return fmt.Errorf("could not create JWT provider: %w", err)
	}

	service := loyalty.NewService(s.db, jwtProvider, s.config.AccrualAddress)

	handler := loyalty.NewHandler(service, middleware.JWTAuth(jwtProvider))

	handler.RegisterRoutes(s.router)

	err = s.db.PingContext(context.Background())
	if err != nil {
		return fmt.Errorf("an error trying to ping the DB: %w", err)
	}

	err = goose.Up(s.db.SQLDB(), "./migrations")
	if err != nil {
		return fmt.Errorf("an error occurred when starting Server: %w", err)
	}

	go s.ListenAndServe()

	term := make(chan os.Signal, 1)
	signal.Notify(term, syscall.SIGTERM, syscall.SIGINT)
	<-term

	goose.Down(s.db.SQLDB(), "./migrations")

	return nil
}
