package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/superles/yapmetrics/internal/storage"
	"testing"
)

func Test_gauge(t *testing.T) {
	type args struct {
		name  string
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    float64
		isError bool
	}{
		{
			name: "positive #1",
			args: args{
				name:  "testSetGet1256",
				value: "1256.000",
			},
			want: 1256,
		},
		{
			name: "negative #2",
			args: args{
				name:  "testSetGet1256",
				value: "blablabla",
			},
			want:    1256,
			isError: true,
		},
		{
			name: "positive #3",
			args: args{
				name:  "test354",
				value: "1234.asf",
			},
			want:    1234,
			isError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := gauge(test.args.name, test.args.value)
			// Если присутсвует проверка на ошибку - выходим
			if test.isError && err != nil {
				assert.Error(t, err)
				return
			}
			assert.NoErrorf(t, err, "Ошибка записи в репозиторий")
			assert.Equalf(t, test.want, got, "gauge(%v, %v)", test.args.name, test.args.value)
			item, getErr := storage.MetricRepository.Get(test.args.name)
			assert.NoErrorf(t, getErr, "Ошибка получения из репозитория")
			assert.Equalf(t, item.Value.(float64), got, "Не совпадают данные в репозитории и указанные(%v, %v)", test.args.name, test.args.value)

		})
	}
}
