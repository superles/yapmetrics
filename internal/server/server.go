package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/superles/yapmetrics/internal/server/config"
	pages "github.com/superles/yapmetrics/internal/server/handlers"
	"net/http"
)

func Run() {

	config.InitConfig()

	router := httprouter.New()

	router.POST(`/update/:type/:name/:value`, pages.UpdatePage)
	//для тестов
	router.GET(`/update/:type/:name/:value`, pages.UpdatePage)
	router.GET(`/value/:type/:name`, pages.ValuePage)
	router.GET(`/`, pages.MainPage)

	err := http.ListenAndServe(config.ServerConfig.Endpoint, router)
	if err != nil {
		panic(err)
	}
}
