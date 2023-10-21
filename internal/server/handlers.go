package server

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mailru/easyjson"
	"github.com/superles/yapmetrics/internal/metric"
	"github.com/superles/yapmetrics/internal/utils/logger"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
)

func printValue(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func (s *Server) dumpStorage() {
	if s.config.StoreInterval == 0 {
		go func() {
			if err := s.dump(); err != nil {
				logger.Log.Fatal(err.Error())
			}
		}()
	}
}

func (s *Server) MainPage(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
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

	collection := s.storage.GetAll(context.Background())

	var tplBuf bytes.Buffer
	err = t.Execute(&tplBuf, collection)

	check(err)

	res.WriteHeader(http.StatusOK)
	_, _ = res.Write(tplBuf.Bytes())
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
	s.storage.SetFloat(context.Background(), name, floatVar)
	s.dumpStorage()
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(""))
}

func (s *Server) Update(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}(r.Body)

	body, _ := io.ReadAll(r.Body)

	updateData := metric.JSONData{}

	if err := easyjson.Unmarshal(body, &updateData); err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, fmt.Sprintf("ошибка парсинга json: %s", err), http.StatusBadRequest)
		return
	}

	switch updateData.MType {
	case metric.GaugeMetricTypeName:
		if updateData.Value == nil {
			http.Error(w, "отсутсвует значение метрики", http.StatusBadRequest)
			return
		}
		s.storage.SetFloat(context.Background(), updateData.ID, *updateData.Value)
	case metric.CounterMetricTypeName:
		if updateData.Delta == nil {
			http.Error(w, "отсутсвует значение метрики", http.StatusBadRequest)
			return
		}
		s.storage.IncCounter(context.Background(), updateData.ID, *updateData.Delta)
	default:
		logger.Log.Error("неверный тип метрики")
		http.Error(w, "неверный тип метрики", http.StatusBadRequest)
		return
	}

	s.dumpStorage()

	updatedData, _ := s.storage.Get(context.Background(), updateData.ID)

	if updatedJSON, err := updatedData.ToJSON(); err == nil {
		rawBytes, _ := easyjson.Marshal(updatedJSON)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(rawBytes)
	} else {
		http.Error(w, fmt.Sprintf("ошибка конвертации метрики в json: %s", err), http.StatusBadRequest)
	}

}

func (s *Server) Updates(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}(r.Body)

	body, _ := io.ReadAll(r.Body)

	item := metric.JSONDataCollection{}

	if err := easyjson.Unmarshal(body, &item); err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, fmt.Sprintf("ошибка парсинга json: %s", err), http.StatusBadRequest)
		return
	}

	if err := s.storage.SetAll(context.Background(), item.ToMetrics()); err != nil {
		http.Error(w, fmt.Sprintf("ошибка: %s", err), http.StatusBadRequest)
		return
	}

	if true {
		rawBytes, _ := easyjson.Marshal(item)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(rawBytes)
	} else {
		http.Error(w, fmt.Sprintf("ошибка конвертации метрики в json: %s", "error"), http.StatusBadRequest)
	}
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
	s.storage.IncCounter(context.Background(), name, intVar)
	s.dumpStorage()
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(""))
}

func (s *Server) GetPing(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	dsn, _ := pgx.ParseConfig(s.config.DatabaseDsn)

	ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		dsn.Host, dsn.User, dsn.Password, dsn.Database)

	db, err := sql.Open("pgx", ps)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = db.Ping()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	w.WriteHeader(http.StatusOK)

}

func (s *Server) GetValue(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain")

	name := chi.URLParam(r, "name")
	mType := chi.URLParam(r, "type")

	metricItem, ok := s.storage.Get(context.Background(), name)

	if !ok {
		http.Error(w, "метрика не найдена", http.StatusNotFound)
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

func (s *Server) GetJSONValue(w http.ResponseWriter, r *http.Request) {

	logger.Log.Debug("GetJSONValue")

	w.Header().Set("Content-Type", "application/json")

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Log.Error(err.Error())
		}
	}(r.Body)

	body, _ := io.ReadAll(r.Body)

	getData := metric.JSONData{}

	if err := easyjson.Unmarshal(body, &getData); err != nil {
		logger.Log.Error(err.Error())
		http.Error(w, fmt.Sprintf("ошибка парсинга json: %s", err), http.StatusBadRequest)
		return
	}

	metricItem, ok := s.storage.Get(context.Background(), getData.ID)

	logger.Log.Debug(metricItem)

	if !ok {
		http.Error(w, "метрика не найдена", http.StatusNotFound)
		return
	}

	metricType, metricTypeError := metric.StringToType(getData.MType)

	if metricTypeError != nil || metricItem.Type != metricType {
		fmt.Println(errors.New("тип метрики не совпадает"), metricItem)
		http.Error(w, "тип метрики не совпадает", http.StatusBadRequest)
		return
	}

	if updatedJSON, err := metricItem.ToJSON(); err == nil {
		rawBytes, _ := easyjson.Marshal(updatedJSON)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(rawBytes)
	} else {
		http.Error(w, fmt.Sprintf("ошибка конвертации метрики в json: %s", err), http.StatusBadRequest)
		return
	}

}
