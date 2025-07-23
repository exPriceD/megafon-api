package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

// Config описывает всё, что нужно для подключения к Postgres.
type Config struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	MigrationsDir   string
}

// New открывает пул *sql.DB и применяет настройки.
func New(ctx context.Context, cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	cctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err = db.PingContext(cctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
