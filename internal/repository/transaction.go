package repository

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type TransactionRepos struct {
	pool        *pgxpool.Pool
	redisClient *redis.Client
}

func NewTransactionsRepos(pool *pgxpool.Pool, redisClient *redis.Client) Transaction {
	return &TransactionRepos{pool: pool, redisClient: redisClient}
}

func (r *TransactionRepos) Begin(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}

func (r *TransactionRepos) Rollback(ctx context.Context, tx pgx.Tx) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	return tx.Rollback(ctx)
}

func (r *TransactionRepos) RedisTx(ctx context.Context) (redis.Pipeliner, error) {
	clientCtx := r.redisClient.WithContext(ctx)
	return clientCtx.TxPipeline(), nil
}
func (r *TransactionRepos) RedisRollback(ctx context.Context, tx redis.Pipeliner) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	return tx.Discard()
}
func (r *TransactionRepos) RedisCommit(tx redis.Pipeliner) error {
	_, err := tx.Exec()
	return err
}

func (r *TransactionRepos) RedisClient(ctx context.Context) *redis.Client {
	return r.redisClient.WithContext(ctx)
}
