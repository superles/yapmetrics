package pgstorage

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	types "github.com/superles/yapmetrics/internal/metric"
	"github.com/superles/yapmetrics/internal/utils/logger"
)

const tableName = "metrics"

type PgStorage struct {
	db *pgxpool.Pool
}

func New(dsn string) (*PgStorage, error) {

	dbConfig, dbErr := pgxpool.ParseConfig(dsn)

	if dbErr != nil {
		return nil, dbErr
	}

	db, err := pgxpool.NewWithConfig(context.Background(), dbConfig)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, err
	}

	checkTable(db)

	return &PgStorage{db}, nil
}

func checkTable(db *pgxpool.Pool) {
	db.QueryRow(context.Background(), `create table if not exists `+pgx.Identifier{tableName}.Sanitize()+`
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

	rows, err := s.db.Query(ctx, `
SELECT id, type, value from `+pgx.Identifier{tableName}.Sanitize()+`	
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

	row := s.db.QueryRow(ctx, `SELECT id, type, value from `+pgx.Identifier{tableName}.Sanitize()+` where id=$1`, name)
	if row == nil {
		return item, false
	}

	if err := row.Scan(&item.Name, &item.Type, &item.Value); err != nil {
		if errors.Is(err, pgx.ErrNoRows) == false {
			logger.Log.Error(err.Error())
		}
		return item, false
	}

	return item, true
}

func (s *PgStorage) Set(ctx context.Context, data *types.Metric) {

	_, err := s.db.Exec(ctx, "update "+pgx.Identifier{tableName}.Sanitize()+" set value=$1, type=$2; where id = $3", data.Value, data.Type, data.Name)

	if err != nil {
		logger.Log.Error(err.Error())
	}
}

func (s *PgStorage) SetAll(ctx context.Context, data *[]types.Metric) error {
	begin, bErr := s.db.BeginTx(ctx, pgx.TxOptions{})
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
			if err := begin.Rollback(ctx); err != nil {
				return err
			}
			return nil
		}
	}

	if err := begin.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (s *PgStorage) SetFloat(ctx context.Context, Name string, Value float64) {

	query := `INSERT INTO  ` + pgx.Identifier{tableName}.Sanitize() + `  (id, type, value) VALUES($1, 1, $2) on conflict on constraint metrics_pk do UPDATE SET type = 1 , value = $2`

	_, err := s.db.Exec(ctx, query, Name, Value)

	if err != nil {
		logger.Log.Error(err.Error())
	}
}

func (s *PgStorage) IncCounter(ctx context.Context, Name string, Value int64) {
	query := `INSERT INTO ` + pgx.Identifier{tableName}.Sanitize() + ` (id, type, value)
VALUES($1, 2, $2)
on conflict on constraint metrics_pk
do UPDATE SET type = 2, value = ` + pgx.Identifier{tableName, "value"}.Sanitize() + ` + $2;`

	_, err := s.db.Exec(ctx, query, Name, Value)
	if err != nil {
		logger.Log.Error(err.Error())
	}
}
