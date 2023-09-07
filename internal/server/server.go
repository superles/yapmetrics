package server

import (
	pages "github.com/superles/yapmetrics/internal/server/handlers"
	"net/http"
)

func Run() {

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, pages.UpdatePage)
	mux.HandleFunc(`/`, pages.MainPage)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
