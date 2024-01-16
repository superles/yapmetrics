package server

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	types "github.com/superles/yapmetrics/internal/metric"
	"github.com/superles/yapmetrics/internal/server/config"
	"github.com/superles/yapmetrics/internal/server/middleware/auth"
	"github.com/superles/yapmetrics/internal/server/middleware/compress"
	"github.com/superles/yapmetrics/internal/server/middleware/logging"
	"github.com/superles/yapmetrics/internal/utils/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type metricProvider interface {
	GetAll(ctx context.Context) (map[string]types.Metric, error)
	Get(ctx context.Context, name string) (types.Metric, error)
	Set(ctx context.Context, data types.Metric) error
	SetAll(ctx context.Context, data []types.Metric) error
	SetFloat(ctx context.Context, Name string, Value float64) error
	IncCounter(ctx context.Context, Name string, Value int64) error
	Ping(ctx context.Context) error
	Dump(ctx context.Context, path string) error
	Restore(ctx context.Context, path string) error
}

type Server struct {
	storage metricProvider
	router  *chi.Mux
	config  *config.Config
}

func New(s metricProvider, cfg *config.Config) *Server {
	router := chi.NewRouter()
	router.Use(compress.WithCompressGzip, logging.WithLogging, auth.WithAuth(cfg.SecretKey))
	server := &Server{storage: s, router: router, config: cfg}
	server.registerRoutes()
	return server
}

func (s *Server) registerRoutes() {
	s.router.Post("/update/", s.Update)
	s.router.Post("/updates/", s.Updates)
	s.router.Post("/value/", s.GetJSONValue)
	s.router.Post("/update/counter/{name}/{value}", s.UpdateCounter)
	s.router.Post("/update/gauge/{name}/{value}", s.UpdateGauge)
	s.router.Get("/update/counter/{name}/{value}", s.UpdateCounter)
	s.router.Get("/update/gauge/{name}/{value}", s.UpdateGauge)
	s.router.Post("/update/{type}/{name}/{value}", s.BadRequest)
	s.router.Get("/value/{type}/{name}", s.GetValue)
	s.router.Get("/ping", s.GetPing)
	s.router.Get("/", s.MainPage)
}

func (s *Server) load(ctx context.Context) error {
	if len(s.config.FileStoragePath) != 0 {
		return s.storage.Restore(ctx, s.config.FileStoragePath)
	}
	return nil
}

func (s *Server) dump(ctx context.Context) error {
	if len(s.config.FileStoragePath) != 0 {
		return s.storage.Dump(ctx, s.config.FileStoragePath)
	}
	return nil
}

func (s *Server) startDumpWatcher(ctx context.Context) {
	if s.config.StoreInterval > 0 {
		ticker := time.NewTicker(time.Second * time.Duration(s.config.StoreInterval))
		go func() {
			for t := range ticker.C {
				logger.Log.Debug(fmt.Sprintf("Tick at: %v\n", t.UTC()))
				if err := s.dump(ctx); err != nil {
					logger.Log.Fatal(err.Error())
				}
			}
		}()
	}
}

// Run Запуск сервера.
func (s *Server) Run(appContext context.Context) error {

	ctx, done := signal.NotifyContext(appContext, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer done()

	if s.config.Restore && s.config.DatabaseDsn == "" {
		if err := s.load(ctx); err != nil {
			logger.Log.Error(err.Error())
			return err
		} else {
			logger.Log.Debug("бд загружена успешно")
		}
	}

	srv := &http.Server{
		Addr:    s.config.Endpoint,
		Handler: s.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Log.Error(fmt.Sprintf("не могу запустить сервер: %s", err))
		}
	}()

	s.startDumpWatcher(ctx)

	logger.Log.Info("Server Started")
	<-ctx.Done()
	logger.Log.Info("Server Stopped")

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}

	logger.Log.Info("Server Exited Properly")

	if err := s.dump(ctx); err != nil {
		return err
	}

	return nil

}
