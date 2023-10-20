package pgstorage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	types "github.com/superles/yapmetrics/internal/metric"
	"github.com/superles/yapmetrics/internal/utils/logger"
	"strconv"
)

const tableName = "metrics"

type PgStorage struct {
	db *sql.DB
}

func New(dsn string) *PgStorage {

	connCfg, _ := pgx.ParseConfig(dsn)

	ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		connCfg.Host, connCfg.User, connCfg.Password, connCfg.Database)

	db, err := sql.Open("pgx", ps)

	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	checkTable(db)

	return &PgStorage{db}
}

func checkTable(db *sql.DB) {
	db.QueryRow(`create table if not exists public.metrics
(
    id    varchar          not null
        constraint metrics_pk
            primary key,
    type  smallint         not null,
    value double precision not null
);`)
}

func (s *PgStorage) GetAll(ctx context.Context) map[string]types.Metric {

	items := make(map[string]types.Metric)

	rows, err := s.db.QueryContext(ctx, `
SELECT id, type, value from `+tableName+`	
`)
	if rows.Err() != nil {
		logger.Log.Error(err.Error())
		return items
	}

	if err != nil {
		logger.Log.Error(err.Error())
		return items
	}

	for rows.Next() {
		var item types.Metric
		err = rows.Scan(&item.Name, &item.Type, &item.Value)
		if err != nil {
			logger.Log.Error(err.Error())
			continue
		}
		items[item.Name] = item
	}

	return items
}

func (s *PgStorage) Get(ctx context.Context, name string) (types.Metric, bool) {

	item := types.Metric{}

	row := s.db.QueryRowContext(ctx, `
SELECT id, type, value from `+tableName+` where id='`+name+`'	
`)
	if row == nil {
		return item, false
	}

	if err := row.Scan(&item.Name, &item.Type, &item.Value); err != nil {
		logger.Log.Warn(err.Error())
		return item, false
	}

	return item, true
}

func (s *PgStorage) Set(ctx context.Context, data *types.Metric) {
	query := fmt.Sprintf("update %s set value=%f, type=%d; where id = '%s'", tableName, data.Value, data.Type, data.Name)

	result, err := s.db.ExecContext(ctx, query)

	if err != nil {
		logger.Log.Warn(err.Error())
	}

	_, err = result.RowsAffected()

	if err != nil {
		logger.Log.Warn(err.Error())
	}
}

func (s *PgStorage) SetAll(ctx context.Context, data *[]types.Metric) error {
	begin, bErr := s.db.BeginTx(ctx, &sql.TxOptions{})
	if bErr != nil {
		return bErr
	}
	for _, item := range *data {

		switch item.Type {
		case types.GaugeMetricType:
			s.SetFloat(ctx, item.Name, item.Value)
		case types.CounterMetricType:
			s.IncCounter(ctx, item.Name, int64(item.Value))
		default:
			logger.Log.Error("неверный тип метрики")
			if err := begin.Rollback(); err != nil {
				return err
			}
			return nil
		}
	}

	if err := begin.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *PgStorage) SetFloat(ctx context.Context, Name string, Value float64) {
	valueStr := strconv.FormatFloat(Value, 'f', -1, 64)
	query := `
INSERT INTO ` + tableName + ` (id, type, value)
VALUES('` + Name + `',1, ` + valueStr + `)
on conflict on constraint metrics_pk
do UPDATE SET type = 1, value = ` + valueStr + `;
`

	result, err := s.db.ExecContext(ctx, query)

	if err != nil {
		logger.Log.Warn(err.Error())
	}

	_, err = result.RowsAffected()

	if err != nil {
		logger.Log.Warn(err.Error())
	}
}

func (s *PgStorage) IncCounter(ctx context.Context, Name string, Value int64) {
	valueStr := strconv.FormatInt(Value, 10)
	query := `
INSERT INTO ` + tableName + ` (id, type, value)
VALUES('` + Name + `',2, ` + valueStr + `)
on conflict on constraint metrics_pk
do UPDATE SET type = 2, value = ` + tableName + `.value + ` + valueStr + `;
`

	result, err := s.db.ExecContext(ctx, query)

	if err != nil {
		logger.Log.Error(err.Error())
	}

	_, err = result.RowsAffected()

	if err != nil {
		logger.Log.Error(err.Error())
	}
}
