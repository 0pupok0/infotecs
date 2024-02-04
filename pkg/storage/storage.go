package storage

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"sync"
)

type Storage struct {
	mutex *sync.Mutex
	pool  *pgxpool.Pool
}

func New(conn string) (*Storage, error) {
	p, err := pgxpool.Connect(context.Background(), conn)
	if err != nil {
		return nil, err
	}
	return &Storage{mutex: &sync.Mutex{}, pool: p}, nil
}

func (s *Storage) CreateTables() error {
	_, err := s.pool.Query(context.Background(), `
		CREATE TABLE IF NOT EXISTS wallets (
		    id SERIAL PRIMARY KEY,
		    balance FLOAT NOT NULL DEFAULT 100
		)
	`)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	_, err = s.pool.Query(context.Background(), `
		CREATE TABLE IF NOT EXISTS transactions_history (
		    id SERIAL PRIMARY KEY,
		    time TIMESTAMP NOT NULL,
		    from_id INTEGER NOT NULL DEFAULT 0,
		    to_id INTEGER NOT NULL DEFAULT 0,
		    amount FLOAT NOT NULL DEFAULT 0
		)
	`)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
