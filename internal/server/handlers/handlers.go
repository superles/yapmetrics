package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

func UnknownPage(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusBadRequest)
}

func UpdatePage(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	parts := strings.Split(strings.TrimLeft(strings.Trim(req.RequestURI, " "), "/"), "/")
	if len(parts) < 4 {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	metricType := parts[1]

	match, err := regexp.MatchString("^\\w", parts[2])
	if !match && err == nil {
		fmt.Println("метрика должна начинаться с буквы")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		fmt.Println(err.Error())
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	switch metricType {
	case "counter":
		_, err := counter(parts[2], parts[3])
		if err != nil {
			fmt.Println(err.Error())
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		res.WriteHeader(http.StatusOK)
	case "gauge":
		_, err := gauge(parts[2], parts[3])
		if err != nil {
			fmt.Println(err.Error())
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		res.WriteHeader(http.StatusOK)
	default:
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusBadRequest)
}
