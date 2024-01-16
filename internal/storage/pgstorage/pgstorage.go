package pgstorage

import (
	"context"
	"errors"
	"fmt"
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

// New Создание объекта PgStorage
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

func (s *PgStorage) GetAll(ctx context.Context) (map[string]types.Metric, error) {

	items := make(map[string]types.Metric)

	rows, err := s.db.Query(ctx, `
SELECT id, type, value from `+pgx.Identifier{tableName}.Sanitize()+`	
`)
	if err != nil {
		return nil, err
	}

	if rows.Err() != nil {
		return nil, err
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

	return items, nil
}

func (s *PgStorage) Get(ctx context.Context, name string) (types.Metric, error) {

	item := types.Metric{}

	row := s.db.QueryRow(ctx, `SELECT id, type, value from `+pgx.Identifier{tableName}.Sanitize()+` where id=$1`, name)

	if row == nil {
		return item, errors.New("объект row пустой")
	}

	if err := row.Scan(&item.Name, &item.Type, &item.Value); err != nil {
		return item, err
	}

	return item, nil
}

func (s *PgStorage) Set(ctx context.Context, data types.Metric) error {

	_, err := s.db.Exec(ctx, "update "+pgx.Identifier{tableName}.Sanitize()+" set value=$1, type=$2; where id = $3", data.Value, data.Type, data.Name)

	return err
}
func (s *PgStorage) SetAll(ctx context.Context, data []types.Metric) error {
	begin, bErr := s.db.BeginTx(ctx, pgx.TxOptions{})
	if bErr != nil {
		return bErr
	}
	defer func(begin pgx.Tx, ctx context.Context) {
		err := begin.Rollback(ctx)
		if err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			logger.Log.Error(fmt.Sprintf("rollback error: %s", err))
		}
	}(begin, ctx)
	b := &pgx.Batch{}
	for _, item := range data {
		switch item.Type {
		case types.CounterMetricType:
			b.Queue(`INSERT INTO  `+pgx.Identifier{tableName}.Sanitize()+`  (id, type, value)
VALUES($1, 2, $2) 
on conflict on constraint metrics_pk do 
UPDATE SET type = 2 , value = `+pgx.Identifier{tableName, "value"}.Sanitize()+` + $2`, item.Name, item.Value)
		case types.GaugeMetricType:
			b.Queue(`INSERT INTO  `+pgx.Identifier{tableName}.Sanitize()+`  (id, type, value) 
VALUES($1, 1, $2) 
on conflict on constraint metrics_pk do 
UPDATE SET type = 1 , value = $2`, item.Name, item.Value)
		default:
			logger.Log.Error("неверный тип метрики")
			return errors.New("неверный тип метрики")
		}
	}
	if err := begin.SendBatch(ctx, b).Close(); err != nil {
		logger.Log.Error(fmt.Sprintf("batch close error: %s", err))
		return err
	}
	if err := begin.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (s *PgStorage) setAllOld(ctx context.Context, data []types.Metric) error {
	begin, err := s.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	for _, item := range data {

		switch item.Type {
		case types.GaugeMetricType:
			if err := s.SetFloat(ctx, item.Name, item.Value); err != nil {
				return err
			}
		case types.CounterMetricType:
			if err := s.IncCounter(ctx, item.Name, int64(item.Value)); err != nil {
				return err
			}
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

func (s *PgStorage) SetFloat(ctx context.Context, Name string, Value float64) error {

	query := `INSERT INTO  ` + pgx.Identifier{tableName}.Sanitize() + `  (id, type, value) VALUES($1, 1, $2) on conflict on constraint metrics_pk do UPDATE SET type = 1 , value = $2`

	_, err := s.db.Exec(ctx, query, Name, Value)

	return err
}

func (s *PgStorage) IncCounter(ctx context.Context, Name string, Value int64) error {
	query := `INSERT INTO ` + pgx.Identifier{tableName}.Sanitize() + ` (id, type, value)
VALUES($1, 2, $2)
on conflict on constraint metrics_pk
do UPDATE SET type = 2, value = ` + pgx.Identifier{tableName, "value"}.Sanitize() + ` + $2;`

	_, err := s.db.Exec(ctx, query, Name, Value)
	return err
}

func (s *PgStorage) Ping(ctx context.Context) error {
	return s.db.Ping(ctx)
}

func (s *PgStorage) Dump(ctx context.Context, path string) error {
	return nil
}

func (s *PgStorage) Restore(ctx context.Context, path string) error {
	return nil
}
