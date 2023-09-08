package server

import (
	"github.com/superles/yapmetrics/internal/server/config"
	pages "github.com/superles/yapmetrics/internal/server/handlers"
	"net/http"
)

func Run() {

	config.InitConfig()

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, pages.UpdatePage)
	mux.HandleFunc(`/value/`, pages.ValuePage)
	mux.HandleFunc(`/`, pages.MainPage)

	err := http.ListenAndServe(config.ServerConfig.Endpoint, mux)
	if err != nil {
		panic(err)
	}
}
