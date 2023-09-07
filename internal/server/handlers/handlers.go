package handlers

import (
	"fmt"
	"net/http"
	"strings"
)

func MainPage(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Привет!"))
}

func UpdatePage(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	parts := strings.Split(strings.TrimLeft(strings.Trim(req.RequestURI, " "), "/"), "/")
	if len(parts) < 4 {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	metricType := parts[1]
	switch metricType {
	case "counter":
		counter(parts[1], parts[2])
	case "gauge":
		fmt.Println("Linux.")
	default:
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
}
