package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type Storage struct {
	conn *sqlx.DB
}

func NewStorage(dsn string) (*Storage, error) {
	conn, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Printf("[storage.NewStorage]: %v", err)
		return nil, err
	}

	return &Storage{
		conn: conn,
	}, nil
}

func (s *Storage) Save(ctx context.Context, symbol string, average float64, sources map[string]float64) error {
	tx, err := s.conn.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	row := tx.QueryRowx("INSERT INTO prices (symbol, average, recorded_at) VALUES ($1, $2, $3) RETURNING id",
		symbol, average, time.Now())

	var id int
	if err := row.Scan(&id); err != nil {
		return err
	}

	for source, price := range sources {
		_, err = tx.Exec("INSERT INTO source_prices (price_id, source, price) VALUES ($1, $2, $3)",
			id, source, price)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
