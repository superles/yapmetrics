package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/superles/yapmetrics/internal/server/config"
	"github.com/superles/yapmetrics/internal/storage"
	"log"
	"net/http"
)

type Server struct {
	Storage storage.Storage
	Router  *chi.Mux
	Config  *config.Config
}

func New(s storage.Storage) *Server {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	cfg := config.New()
	server := &Server{Storage: s, Router: router, Config: cfg}
	server.registerRoutes()
	return server
}

func (s *Server) registerRoutes() {
	s.Router.Post("/update/counter/{name}/{value}", s.UpdateCounter)
	s.Router.Post("/update/gauge/{name}/{value}", s.UpdateGauge)
	s.Router.Get("/update/counter/{name}/{value}", s.UpdateCounter)
	s.Router.Get("/update/gauge/{name}/{value}", s.UpdateGauge)
	s.Router.Post("/update/{type}/{name}/{value}", s.BadRequest)
	s.Router.Get("/value/{type}/{name}", s.GetValue)
	s.Router.Get("/", s.MainPage)
}

func (s *Server) Run() {

	if err := http.ListenAndServe(s.Config.Endpoint, s.Router); err != nil {
		log.Fatalf("не могу запустить сервер: %s", err)
	}
}
