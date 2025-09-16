package db

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)


type Db struct {
	Pool *pgxpool.Pool
}

func NewDb(connstr string, logger *slog.Logger) *Db {
	const op = "internal.db.go"
	
	pool, err := pgxpool.New(context.Background(), connstr)
	if err != nil {
		logger.Error("Invalid connection Db", slog.String("File error: ", op))
	}
	
	if err := pool.Ping(context.Background()); err != nil {
		logger.Error("Invalid ping to Db", slog.String("File error:", op))
	}

	logger.Info("Db was run successfuly!", slog.String("Storage_path", connstr))

	return &Db{
		Pool: pool,
	}
}