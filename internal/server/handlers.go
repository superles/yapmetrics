package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/mailru/easyjson"
	"github.com/superles/yapmetrics/internal/metric"
	"github.com/superles/yapmetrics/internal/utils/logger"
	"html/template"
	"io"
	"net/http"
	"strconv"
)

func printValue(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

func setError(w http.ResponseWriter, resError error, resErrorText string, resStatus int) {
	logger.Log.Error(fmt.Sprintf("ошибка чтения body: %s", resError))
	http.Error(w, resErrorText, resStatus)
}

func (s *Server) dumpStorage(ctx context.Context) {
	if s.config.StoreInterval == 0 && len(s.config.DatabaseDsn) == 0 {
		go func() {
			if err := s.dump(ctx); err != nil {
				logger.Log.Fatal(err.Error())
			}
		}()
	}
}

func (s *Server) MainPage(res http.ResponseWriter, r *http.Request) {
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

	t, err := template.New("webpage").Funcs(
		template.FuncMap{
			"printValue": printValue,
		},
	).Parse(tpl)

	if err != nil {
		logger.Log.Error(fmt.Sprintf("ошибка: %s", err))
		http.Error(res, "ошибка сервера", http.StatusBadRequest)
		return
	}

	var collection map[string]metric.Metric

	collection, err = s.storage.GetAll(r.Context())

	if err != nil {
		logger.Log.Error(fmt.Sprintf("ошибка: %s", err))
		http.Error(res, "ошибка сервера", http.StatusBadRequest)
		return
	}

	var tplBuf bytes.Buffer
	err = t.Execute(&tplBuf, collection)

	if err != nil {
		logger.Log.Error(fmt.Sprintf("ошибка: %s", err))
		http.Error(res, "ошибка сервера", http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusOK)
	_, err = res.Write(tplBuf.Bytes())
	if err != nil {
		logger.Log.Error(fmt.Sprintf("ошибка записи body: %s", err))
	}
}

func (s *Server) BadRequest(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	_, err := w.Write([]byte(""))
	if err != nil {
		logger.Log.Error(fmt.Sprintf("ошибка записи body: %s", err))
	}
}

func (s *Server) UpdateGauge(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	floatVar, err := strconv.ParseFloat(value, 64)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("cannot parse gauge metric: %s", err))
		http.Error(w, "ошибка сервера", http.StatusBadRequest)
		return
	}
	err = s.storage.SetFloat(r.Context(), name, floatVar)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("ошибка обновления gauge: %s", err))
		http.Error(w, "ошибка сервера", http.StatusBadRequest)
		return
	}
	s.dumpStorage(r.Context())
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(""))
	if err != nil {
		logger.Log.Error(fmt.Sprintf("ошибка записи body: %s", err))
	}
}

func (s *Server) Update(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)

	if err != nil {
		setError(w, err, "ошибка запроса", http.StatusBadRequest)
		return
	}

	updateData := metric.JSONData{}

	if err := easyjson.Unmarshal(body, &updateData); err != nil {
		setError(w, err, "ошибка запроса", http.StatusBadRequest)
		return
	}

	switch updateData.MType {
	case metric.GaugeMetricTypeName:
		if updateData.Value == nil {
			setError(w, errors.New("отсутствует значение метрики"), "ошибка запроса", http.StatusBadRequest)
			return
		}
		if err := s.storage.SetFloat(r.Context(), updateData.ID, *updateData.Value); err != nil {
			logger.Log.Error(err.Error())
			http.Error(w, "ошибка сервера", http.StatusInternalServerError)
			return
		}
	case metric.CounterMetricTypeName:
		if updateData.Delta == nil {
			setError(w, errors.New("отсутствует значение метрики"), "ошибка запроса", http.StatusBadRequest)
			return
		}
		if err := s.storage.IncCounter(r.Context(), updateData.ID, *updateData.Delta); err != nil {
			logger.Log.Error(err.Error())
			http.Error(w, "ошибка сервера", http.StatusInternalServerError)
			return
		}
	default:
		logger.Log.Error(fmt.Sprintf("неверный тип метрики: %s", updateData.MType))
		http.Error(w, "неверный тип метрики", http.StatusBadRequest)
		return
	}

	s.dumpStorage(r.Context())

	updatedData, _ := s.storage.Get(r.Context(), updateData.ID)

	if updatedJSON, err := updatedData.ToJSON(); err == nil {
		rawBytes, _ := easyjson.Marshal(updatedJSON)
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(rawBytes)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("ошибка записи body: %s", err))
		}
	} else {
		logger.Log.Error(fmt.Sprintf("ошибка конвертации метрики в json: %s", err))
		http.Error(w, "ошибка сервера", http.StatusBadRequest)
	}

}

func (s *Server) Updates(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)

	if err != nil {
		logger.Log.Error(fmt.Sprintf("ошибка чтения body: %s", err))
		http.Error(w, "ошибка сервера", http.StatusBadRequest)
		return
	}

	item := metric.JSONDataCollection{}

	if err := easyjson.Unmarshal(body, &item); err != nil {
		logger.Log.Error(fmt.Sprintf("ошибка парсинга json: %s", err))
		http.Error(w, "ошибка сервера", http.StatusBadRequest)
		return
	}

	metrics := item.ToMetrics()

	if err := s.storage.SetAll(r.Context(), &metrics); err != nil {
		logger.Log.Error(fmt.Sprintf("ошибка setall: %s", err))
		http.Error(w, "ошибка сервера", http.StatusInternalServerError)
		return
	}

	s.dumpStorage(r.Context())

	if rawBytes, err := easyjson.Marshal(item); err == nil {
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(rawBytes)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("ошибка записи body: %s", err))
		}
	} else {
		logger.Log.Error(fmt.Sprintf("ошибка сериализации: %s", err))
		http.Error(w, "ошибка сервера", http.StatusBadRequest)
	}
}

func (s *Server) UpdateCounter(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	name := chi.URLParam(r, "name")
	value := chi.URLParam(r, "value")

	intVar, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("cannot parse counter metric: %s", err))
		http.Error(w, "ошибка сервера", http.StatusBadRequest)
		return
	}
	err = s.storage.IncCounter(r.Context(), name, intVar)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("ошибка обновления counter: %s", err))
		http.Error(w, "ошибка сервера", http.StatusBadRequest)
		return
	}
	s.dumpStorage(r.Context())
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(""))
	if err != nil {
		logger.Log.Error(fmt.Sprintf("ошибка записи body: %s", err))
	}
}

func (s *Server) GetPing(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	err := s.storage.Ping(r.Context())

	if err != nil {
		logger.Log.Error(fmt.Sprintf("ошибка подключения к бд: %s", err))
		http.Error(w, "ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func (s *Server) GetValue(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/plain")

	name := chi.URLParam(r, "name")
	mType := chi.URLParam(r, "type")

	metricItem, err := s.storage.Get(r.Context(), name)

	if err != nil {
		logger.Log.Warn(fmt.Sprintf("метрика не найдена: %s, ошибка: %s", name, err))
		http.Error(w, "метрика не найдена", http.StatusNotFound)
		return
	}

	metricType, metricTypeError := metric.StringToType(mType)

	if metricTypeError != nil || metricItem.Type != metricType {
		logger.Log.Error(fmt.Sprintf("тип метрики не совпадает: %d != %d, имя: %s", metricItem.Type, metricType, name))
		http.Error(w, "тип метрики не совпадает", http.StatusBadRequest)
		return
	}

	var value string

	value, err = metricItem.String()

	if err != nil {
		logger.Log.Error(fmt.Sprintf("ошибка конвертирования metric->string: %s", err))
		http.Error(w, "ошибка сервера", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(value)); err != nil {
		logger.Log.Error(fmt.Sprintf("ошибка записи body: %s", err))
	}

}

func (s *Server) GetJSONValue(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)

	if err != nil {
		logger.Log.Error(fmt.Sprintf("ошибка чтения body: %s", err))
		http.Error(w, "ошибка сервера", http.StatusBadRequest)
		return
	}

	getData := metric.JSONData{}

	if err := easyjson.Unmarshal(body, &getData); err != nil {
		logger.Log.Error(fmt.Sprintf("ошибка парсинга json: %s", err))
		http.Error(w, "ошибка сервера", http.StatusBadRequest)
		return
	}

	metricItem, metricErr := s.storage.Get(r.Context(), getData.ID)

	if metricErr != nil {
		logger.Log.Warn(fmt.Sprintf("метрика не найдена: %s, ошибка: %s", getData.ID, err))
		http.Error(w, "метрика не найдена", http.StatusNotFound)
		return
	}

	metricType, metricTypeError := metric.StringToType(getData.MType)

	if metricTypeError != nil || metricItem.Type != metricType {
		logger.Log.Warn(fmt.Sprintf("тип метрики не совпадает: %d != %d, имя: %s", metricItem.Type, metricType, getData.ID))
		http.Error(w, "тип метрики не совпадает", http.StatusBadRequest)
		return
	}

	var updatedJSON *metric.JSONData

	updatedJSON, err = metricItem.ToJSON()

	if err != nil {
		logger.Log.Error(fmt.Sprintf("ошибка конвертации метрики в json: %s", err))
		http.Error(w, "ошибка сервера", http.StatusBadRequest)
		return
	}

	if rawBytes, err := easyjson.Marshal(updatedJSON); err == nil {
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(rawBytes)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("ошибка записи body: %s", err))
		}
	} else {
		logger.Log.Error(fmt.Sprintf("ошибка сериализации: %s", err))
	}

}
