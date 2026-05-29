package di

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/duddy57/cotrip-api/internal/api"
	"github.com/duddy57/cotrip-api/internal/api/spec"
	"github.com/duddy57/cotrip-api/internal/config"
	"github.com/duddy57/cotrip-api/internal/infra"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/klauspost/compress/gzhttp"
	"github.com/phenpessoa/gutils/netutils/httputils"
	"go.uber.org/zap"
)

type Application struct {
	r    *chi.Mux
	pool *pgxpool.Pool
	l    *zap.Logger
	m    *infra.Mail
}

func Run(ctx context.Context) error {
	cfg := config.LoadEnv()

	lg, err := config.NewLogger(cfg.ServerConfig.Environment)
	if err != nil {
		return err
	}
	defer func() { _ = lg.Sync() }()

	pool, err := infra.NewDB(
		ctx,
		cfg.DatabaseConfig.CotripDatabaseUser,
		cfg.DatabaseConfig.CotripDatabasePassword,
		cfg.DatabaseConfig.CotripDatabaseHost,
		cfg.DatabaseConfig.CotripDatabasePort,
		cfg.DatabaseConfig.CotripDatabaseName,
	)
	if err != nil {
		return err
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		return err
	}

	infra.NewMail(pool, cfg.MailerConfig.ApiKey, cfg.MailerConfig.EmailFrom)

	r := chi.NewMux()
	app := Application{
		r:    r,
		pool: pool,
		l:    lg,
		m:    infra.NewMail(pool, cfg.MailerConfig.ApiKey, cfg.MailerConfig.EmailFrom),
	}

	app.mount()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.ServerConfig.Port),
		Handler:      gzhttp.GzipHandler(r),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     zap.NewStdLog(lg.Named("journey_http_error_logger")),
	}

	defer func() {
		const timeout = 30 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		lg.Info("Sending shutdown signal to HTTP server", zap.Duration("timeout", timeout))
		if shutdownErr := srv.Shutdown(ctx); shutdownErr != nil {
			lg.Error("Failed to shutdown HTTP server", zap.Error(err))
			err = errors.Join(err, shutdownErr)
		}
		lg.Info("HTTP Server shutdown")
	}()

	errChan := make(chan error, 1)

	go func() {
		lg.Info("Cotrip is starting up")
		if err := srv.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				err = nil
			}
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		lg.Info("Received shutdown signal, shutting down...")
	case err = <-errChan:
		if err != nil {
			lg.Error("HTTP Server error, shutting down...", zap.Error(err))
		}
	}

	return nil

}

func (app *Application) mount() {
	apiHandlers := api.NewAPI(app.pool, app.m, app.l)
	app.r.Use(middleware.RequestID)
	app.r.Use(middleware.Recoverer)
	app.r.Use(httputils.ChiLogger(app.l.Named("mux")))
	app.r.Mount("/", spec.Handler(&apiHandlers, spec.WithErrorHandler(apiHandlers.ErrorHandlerFunc)))
	app.l.Info("Handlers mounted")
}
