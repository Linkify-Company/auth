package postgres

import (
	"auth/internal/config"
	"auth/pkg/errify"
	"auth/pkg/logger"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"strconv"
	"strings"
)

type tracer struct {
	log logger.Logger
}

func (t tracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	t.log.Infof(strings.ReplaceAll(data.SQL, "\n", " "))
	return ctx
}

func (t tracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {}

func New(ctx context.Context, log logger.Logger, debug bool) (*pgxpool.Pool, errify.IError) {
	port, err := strconv.Atoi(os.Getenv(config.PostgresPort))
	if err != nil {
		return nil, errify.NewInternalServerError(err.Error(), "New/Atoi")
	}

	cfg, err := pgxpool.ParseConfig(fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s pool_max_conns=32",
		os.Getenv(config.PostgresHost),
		port,
		os.Getenv(config.PostgresUser),
		os.Getenv(config.PostgresPassword),
		os.Getenv(config.PostgresDB),
	))
	if err != nil {
		return nil, errify.NewInternalServerError(err.Error(), "New/ParseConfig")
	}
	if debug {
		cfg.ConnConfig.Tracer = tracer{log: log}
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, errify.NewInternalServerError(err.Error(), "New/NewWithConfig")
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, errify.NewInternalServerError(err.Error(), "New/Ping")
	}
	return pool, nil
}
