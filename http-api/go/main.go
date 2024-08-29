package main

import (
	"cmp"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"toasterexample/internal/handlers"
	"toasterexample/internal/services"
	"toasterexample/internal/stores"

	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGTERM)
	defer stop()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cmp.Or(os.Getenv("DATABASE_USER"), "toaster"),
		cmp.Or(os.Getenv("DATABASE_PASSWORD"), "password"),
		cmp.Or(os.Getenv("DATABASE_HOST"), "localhost"),
		cmp.Or(os.Getenv("DATABASE_PORT"), "5432"),
		cmp.Or(os.Getenv("DATABASE_NAME"), "toaster"),
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		slog.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}

	go func() {
		<-ctx.Done()
		slog.Info("closing database connection")
		_ = db.Close()
	}()

	// Init stores
	bookStore := stores.NewBookStore(db)

	// Init services
	libraryService := services.NewLibraryService(bookStore)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /books/{$}", handlers.ListBooks(libraryService))
	mux.HandleFunc("GET /books/{bookID}", handlers.GetBook(libraryService))
	mux.HandleFunc("POST /books/{$}", handlers.CreateBook(libraryService))

	var httpHandler http.Handler = mux
	httpHandler = middleware.Compress(5)(httpHandler)
	httpHandler = middleware.Logger(httpHandler)
	httpHandler = middleware.Recoverer(httpHandler)
	httpHandler = middleware.Heartbeat("/health")(httpHandler)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           httpHandler,
		ReadHeaderTimeout: 3 * time.Second,
	}

	go func() {
		slog.Info("starting server on :8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed unexpectedly", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutting down server")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			slog.Error("shutdown timed out, exiting immediately")
			os.Exit(1)
		}

		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to shutdown http server", slog.Any("error", err))
			os.Exit(1)
		}
	}

	slog.Info("server shutdown successfully")
}
