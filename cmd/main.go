package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"github.com/bread-amammay/unitofwork/gen/api/blogs/v1/blogsv1connect"
	"github.com/bread-amammay/unitofwork/internal/api"
	"github.com/bread-amammay/unitofwork/internal/identity"
	"github.com/bread-amammay/unitofwork/internal/storage"
	"github.com/bread-amammay/unitofwork/internal/storage/migrations/postgres"
	"github.com/bread-amammay/unitofwork/internal/usecase"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {

	logger := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	logger.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	logger.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	logger.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}

	log := zerolog.New(logger).With().Timestamp().Logger()

	if err := run(log); err != nil {
		log.Fatal().Err(err).Msg("error running service")
	}

	log.Info().Msg("service stopped")

}

func run(z zerolog.Logger) error {
	z.Info().Msg("starting service")
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// connect to database
	open, err := sqlx.Open("pgx", "user=postgres password=postgres dbname=demo host=localhost port=5432 sslmode=disable")
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	err = open.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}

	goose.SetBaseFS(postgres.EmbedMigrations)
	err = goose.Up(open.DB, "migrations")
	if err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	z.Info().Msg("database connection established")

	work := storage.New(open)
	controller := usecase.NewController(work)
	blogHandler := api.NewServer(controller)

	mux := http.NewServeMux()
	reflector := grpcreflect.NewStaticReflector(
		blogsv1connect.BlogServiceName,
	)
	mux.Handle(blogsv1connect.NewBlogServiceHandler(blogHandler, connect.WithInterceptors(identity.ConnectInterceptor(z))))
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	srv := http.Server{
		Addr:    "localhost:8080",
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	go func() {
		z.Info().Msgf("listening on %s", srv.Addr)
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			z.Fatal().Err(err).Msg("error running server")
		}
	}()

	<-ctx.Done()

	z.Info().Msg("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = srv.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("error shutting down server: %w", err)
	}

	z.Info().Msg("server shutdown complete")

	err = open.Close()
	if err != nil {
		return fmt.Errorf("error closing database: %w", err)
	}

	z.Info().Msg("database connection closed")

	return nil
}
