package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	handlerhttp "github.com/lisvindanu/anaphase-cli/internal/adapter/handler/http"
	"github.com/lisvindanu/anaphase-cli/internal/adapter/repository/postgres"
)

// App holds all application dependencies
type App struct {
	logger *slog.Logger
	db     *pgxpool.Pool

	customerHandler *handlerhttp.CustomerHandler
}

// InitializeApp initializes all application dependencies
func InitializeApp(logger *slog.Logger) (*App, error) {
	// Database connection
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/anaphase?sslmode=disable"
	}

	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	// Ping database
	if err := db.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	logger.Info("database connected")

	// Initialize customer dependencies
	customerRepo := postgres.NewCustomerRepository(db)
	_ = customerRepo // TODO: Pass to service when implemented
	// TODO: Create customer service implementation
	// customerService := service.NewCustomerService(customerRepo)
	customerHandler := handlerhttp.NewCustomerHandler(nil, logger) // Pass service when implemented

	return &App{
		logger: logger,
		db:     db,
		customerHandler: customerHandler,
	}, nil
}

// RegisterRoutes registers all HTTP routes
func (a *App) RegisterRoutes(r chi.Router) {
	a.customerHandler.RegisterRoutes(r)
}

// Cleanup cleans up application resources
func (a *App) Cleanup() {
	if a.db != nil {
		a.db.Close()
		a.logger.Info("database connection closed")
	}
}
