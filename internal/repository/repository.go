package repository

import (
	"auth/internal/domain"
	"context"
	"github.com/Linkify-Company/common_utils/errify"
	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v5"
	"time"
)

type User interface {
	AddUser(ctx context.Context, tx pgx.Tx, email string, role domain.Role, passHash []byte) (int, error)
	UserById(ctx context.Context, tx pgx.Tx, id int) (*domain.UserFromDB, error)
	UserByEmail(ctx context.Context, tx pgx.Tx, email string) (*domain.UserFromDB, error)
}

type Auth interface {
	Authorization(ctx context.Context, tx redis.Pipeliner, user *domain.AuthData, refreshTTL time.Duration, accessTTL time.Duration, secret string) (string, error)
	RenewalAuthorization(ctx context.Context, redisClient *redis.Client, accessToken string, accessTTL time.Duration, secret string) (string, error)
	CheckAuthorization(ctx context.Context, redisClient *redis.Client, accessToken string) (*domain.AuthData, error)
	RemoveAuthorization(ctx context.Context, tx redis.Pipeliner, accessToken string) error
}

type Email interface {
	IsExist(ctx context.Context, email string) bool
	Set(ctx context.Context, email string, code int) errify.IError
	IsValid(ctx context.Context, email string, code int) bool
}

type Transaction interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Rollback(ctx context.Context, tx pgx.Tx) error

	RedisTx(ctx context.Context) (redis.Pipeliner, error)
	RedisRollback(ctx context.Context, tx redis.Pipeliner) error
	RedisCommit(tx redis.Pipeliner) error
	RedisClient(ctx context.Context) *redis.Client
}

type Repository struct {
	User
	Auth
	Email
}

func NewRepository() *Repository {
	return &Repository{
		User:  NewUserRepos(),
		Auth:  NewAuthRepo(),
		Email: NewEmailRepos(),
	}
}
