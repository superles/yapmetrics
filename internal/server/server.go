package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	types "github.com/superles/yapmetrics/internal/metric"
	"github.com/superles/yapmetrics/internal/server/config"
	"log"
	"net/http"
)

type metricProvider interface {
	GetAll() map[string]types.Metric
	Get(name string) (types.Metric, error)
	SetFloat(Name string, Value float64)
	IncCounter(Name string, Value int64)
}

type Server struct {
	storage metricProvider
	router  *chi.Mux
	config  *config.Config
}

func New(s metricProvider) *Server {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	cfg := config.New()
	server := &Server{storage: s, router: router, config: cfg}
	server.registerRoutes()
	return server
}

func (s *Server) registerRoutes() {
	s.router.Post("/update/counter/{name}/{value}", s.UpdateCounter)
	s.router.Post("/update/gauge/{name}/{value}", s.UpdateGauge)
	s.router.Get("/update/counter/{name}/{value}", s.UpdateCounter)
	s.router.Get("/update/gauge/{name}/{value}", s.UpdateGauge)
	s.router.Post("/update/{type}/{name}/{value}", s.BadRequest)
	s.router.Get("/value/{type}/{name}", s.GetValue)
	s.router.Get("/", s.MainPage)
}

func (s *Server) Run() {

	if err := http.ListenAndServe(s.config.Endpoint, s.router); err != nil {
		log.Fatalf("не могу запустить сервер: %s", err)
	}
}
