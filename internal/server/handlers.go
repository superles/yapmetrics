package server

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

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
			  {{if eq .Type "counter"}}
			 	 <td>{{ printf "%d" .ValueInt }}</td>
			  {{else}}
			 	 <td>{{ printf "%.2f" .ValueFloat }}</td>
			  {{end}}
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

	collection := s.Storage.GetAll()

	err = t.Execute(res, collection)

	check(err)
}

func (s *Server) BadRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Server) UpdateGauge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	floatVar, err := strconv.ParseFloat(value, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot parse gauge metric: %s", err), http.StatusBadRequest)
	}
	s.Storage.SetFloat(name, floatVar)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) UpdateCounter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	intVar, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("cannot parse counter metric: %s", err), http.StatusBadRequest)
	}
	s.Storage.IncCounter(name, intVar)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetValue(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain")

	name := chi.URLParam(r, "name")
	mType := chi.URLParam(r, "type")

	metric, err := s.Storage.Get(name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	if metric.Type != mType {
		fmt.Println(errors.New("тип метрики не совпадает"), metric)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	value, err := metric.String()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)

	_, writeErr := w.Write([]byte(value))
	if writeErr != nil {
		return
	}

}
