package server

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/superles/yapmetrics/internal/metric"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func printValue(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func (s *Server) MainPage(res http.ResponseWriter, req *http.Request) {

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
			  <td>{{ printValue .Value }}</td>
			{{end}}
		</table>
	</body>
</html>`

	check := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	t, err := template.New("webpage").Funcs(
		template.FuncMap{
			"printValue": printValue,
		},
	).Parse(tpl)

	check(err)

	collection := s.storage.GetAll()

	err = t.Execute(res, collection)

	check(err)
}

func (s *Server) BadRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte(""))
}

func (s *Server) UpdateGauge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	floatVar, err := strconv.ParseFloat(value, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot parse gauge metric: %s", err), http.StatusBadRequest)
		return
	}
	s.storage.SetFloat(name, floatVar)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(""))
}

func (s *Server) UpdateCounter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	intVar, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot parse counter metric: %s", err), http.StatusBadRequest)
		return
	}
	s.storage.IncCounter(name, intVar)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(""))
}

func (s *Server) GetValue(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain")

	name := chi.URLParam(r, "name")
	mType := chi.URLParam(r, "type")

	metricItem, err := s.storage.Get(name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	metricType, metricTypeError := metric.StringToType(mType)

	if metricTypeError != nil || metricItem.Type != metricType {
		fmt.Println(errors.New("тип метрики не совпадает"), metricItem)
		http.Error(w, "тип метрики не совпадает", http.StatusBadRequest)
		return
	}

	value, err := metricItem.String()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, writeErr := w.Write([]byte(value))
	if writeErr != nil {
		return
	}

}
