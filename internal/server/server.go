package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	types "github.com/superles/yapmetrics/internal/metric"
	"github.com/superles/yapmetrics/internal/server/config"
	"github.com/superles/yapmetrics/internal/server/middleware/compress"
	"github.com/superles/yapmetrics/internal/server/middleware/logging"
	"github.com/superles/yapmetrics/internal/utils/logger"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type metricProvider interface {
	GetAll(ctx context.Context) (map[string]types.Metric, error)
	Get(ctx context.Context, name string) (types.Metric, error)
	Set(ctx context.Context, data *types.Metric) error
	SetAll(ctx context.Context, data *[]types.Metric) error
	SetFloat(ctx context.Context, Name string, Value float64) error
	IncCounter(ctx context.Context, Name string, Value int64) error
	Ping(ctx context.Context) error
}

type Server struct {
	storage metricProvider
	router  *chi.Mux
	config  *config.Config
}

func New(s metricProvider, cfg *config.Config) *Server {
	router := chi.NewRouter()
	router.Use(compress.WithCompressGzip, logging.WithLogging)
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

func (s *Server) load() error {
	file, fileErr := os.OpenFile(s.config.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0666)
	if fileErr != nil {
		return fileErr
	}
	dec := json.NewDecoder(file)

	for {
		var m types.Metric
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if err := s.storage.Set(context.Background(), &m); err != nil {
			return err
		}
	}

	if err := file.Close(); err != nil {
		return err
	}
	logger.Log.Info("load success")
	return nil
}

func (s *Server) dump() error {
	file, fileErr := os.OpenFile(s.config.FileStoragePath, os.O_WRONLY|os.O_CREATE, 0666)
	if fileErr != nil {
		return fileErr
	}
	if err := file.Truncate(0); err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	metrics, err := s.storage.GetAll(context.Background())
	if err != nil {
		return err
	}
	for _, metric := range metrics {
		err := encoder.Encode(&metric)
		if err != nil {
			return fileErr
		}
	}
	if err := file.Close(); err != nil {
		return err
	}
	logger.Log.Info("dump success")
	return nil
}

func (s *Server) startDumpWatcher() {
	if s.config.StoreInterval > 0 {
		ticker := time.NewTicker(time.Second * time.Duration(s.config.StoreInterval))
		go func() {
			for t := range ticker.C {
				logger.Log.Debug(fmt.Sprintf("Tick at: %v\n", t.UTC()))
				if err := s.dump(); err != nil {
					logger.Log.Fatal(err.Error())
				}
			}
		}()
	}
}

func (s *Server) Run() error {

	if s.config.Restore && s.config.DatabaseDsn == "" {
		if err := s.load(); err != nil {
			logger.Log.Error(err.Error())
			return err
		}
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := &http.Server{
		Addr:    s.config.Endpoint,
		Handler: s.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Log.Error(fmt.Sprintf("не могу запустить сервер: %s", err))
		}
	}()

	s.startDumpWatcher()

	logger.Log.Info("Server Started")
	<-done
	logger.Log.Info("Server Stopped")

	if err := srv.Shutdown(context.Background()); err != nil {
		return err
	}

	logger.Log.Info("Server Exited Properly")

	if err := s.dump(); err != nil {
		return err
	}

	return nil

}
