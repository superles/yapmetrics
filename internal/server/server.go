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
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type metricProvider interface {
	GetAll() map[string]types.Metric
	Get(name string) (types.Metric, bool)
	Set(data *types.Metric)
	SetFloat(Name string, Value float64)
	IncCounter(Name string, Value int64)
}

type Server struct {
	storage metricProvider
	router  *chi.Mux
	config  *config.Config
}

func New(s metricProvider) *Server {
	cfg := config.New()
	router := chi.NewRouter()
	err := logger.Initialize(cfg.LogLevel)
	if err != nil {
		log.Panicln("ошибка инициализации логера", err.Error())
	}

	router.Use(compress.WithCompressGzip, logging.WithLogging)
	server := &Server{storage: s, router: router, config: cfg}
	server.registerRoutes()
	return server
}

func (s *Server) registerRoutes() {
	s.router.Post("/update/", s.Update)
	s.router.Post("/value/", s.GetJSONValue)
	s.router.Post("/update/counter/{name}/{value}", s.UpdateCounter)
	s.router.Post("/update/gauge/{name}/{value}", s.UpdateGauge)
	s.router.Get("/update/counter/{name}/{value}", s.UpdateCounter)
	s.router.Get("/update/gauge/{name}/{value}", s.UpdateGauge)
	s.router.Post("/update/{type}/{name}/{value}", s.BadRequest)
	s.router.Get("/value/{type}/{name}", s.GetValue)
	s.router.Get("/", s.MainPage)
}

func (s *Server) Load() error {
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
			log.Fatal(err)
		}
		s.storage.Set(&m)
	}

	if err := file.Close(); err != nil {
		return err
	}
	logger.Log.Info("load success")
	return nil
}

func (s *Server) Dump() error {
	file, fileErr := os.OpenFile(s.config.FileStoragePath, os.O_WRONLY|os.O_CREATE, 0666)
	if fileErr != nil {
		return fileErr
	}
	if err := file.Truncate(0); err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	for _, metric := range s.storage.GetAll() {
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

func (s *Server) Run() {

	if s.config.Restore {
		if err := s.Load(); err != nil {
			logger.Log.Fatal(err.Error())
		}
	}

	if s.config.StoreInterval > 0 {
		go func() {
			for range time.Tick(time.Second * time.Duration(s.config.StoreInterval)) {
				if err := s.Dump(); err != nil {
					logger.Log.Fatal(err.Error())
				}
			}
		}()
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

	logger.Log.Info("Server Started")
	<-done
	logger.Log.Info("Server Stopped")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}

	logger.Log.Info("Server Exited Properly")

	if err := s.Dump(); err != nil {
		logger.Log.Fatal(err.Error())
	}

}
