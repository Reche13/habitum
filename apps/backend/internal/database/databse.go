package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/reche13/habitum/internal/config"
)

type Database struct {
	Pool *pgxpool.Pool
}

func New(cfg *config.Config) (*Database, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("parse db config: %w", err)
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create db pool: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return &Database{Pool: pool}, nil
}

func (db *Database) Close() {
	db.Pool.Close()
}
