package handlers

import (
	"errors"
	"fmt"
	"github.com/superles/yapmetrics/internal/memstorage"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
)

func MainPage(res http.ResponseWriter, req *http.Request) {

	if req.RequestURI != "/" {
		UnknownPage(res, req)
		return
	}

	const tpl = `
<!DOCTYPE html>
<html lang="ru">
	<head>
		<meta charset="UTF-8">
		<title>Таблица метрик</title>
	</head>
	<body>
		<h4>Таблица метрик</h4>
		<table>
			{{range .}}    
			  <tr>
			  <td>{{.Name}}</td><td>
			  <td>{{.Value  }}</td>
			{{end}}
		</table>
	</body>
</html>`

	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}
	t, err := template.New("webpage").Parse(tpl)

	check(err)

	collection := memstorage.Storage.Collection

	err = t.Execute(res, collection)

	check(err)
}

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
	fmt.Println("update", parts)

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
		return
	case "gauge":
		_, err := gauge(parts[2], parts[3])
		if err != nil {
			fmt.Println(err.Error())
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		res.WriteHeader(http.StatusOK)
		return
	default:
		res.WriteHeader(http.StatusBadRequest)
		return
	}
}

func ValuePage(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "text/plain; charset=utf-8")

	parts := strings.Split(strings.TrimLeft(strings.Trim(req.RequestURI, " "), "/"), "/")

	if len(parts) < 3 {
		res.WriteHeader(http.StatusNotFound)
		return
	}
	fmt.Println("value", parts)
	metricType := parts[1]
	metricName := parts[2]
	if metricType != "counter" && metricType != "gauge" {
		fmt.Println("метрика должна начинаться с буквы")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	match, err := regexp.MatchString("^\\w", metricName)
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

	metric, metricError := memstorage.Storage.Get(metricName)
	if metricError != nil {
		fmt.Println(metricError.Error())
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if metric.Type != metricType {
		fmt.Println(errors.New("тип метрики не совпадает"))
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	_, writeErr := res.Write([]byte(metric.Value))
	if writeErr != nil {
		return
	}
	res.WriteHeader(http.StatusOK)
}
