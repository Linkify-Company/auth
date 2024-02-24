package service

import (
	"auth/internal/config"
	"auth/internal/domain"
	"auth/internal/repository"
	"auth/pkg/errify"
	"auth/pkg/logger"
	"context"
	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type User interface {
	AddUser(ctx context.Context, user *domain.User) (int, errify.IError)
	GetUserByID(ctx context.Context, id int) (*domain.User, errify.IError)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, errify.IError)
}

type Auth interface {
	Authorization(ctx context.Context, auth *domain.Auth, cfg config.TokenConfig) (string, errify.IError)
	CheckAuthorization(ctx context.Context, accessToken string) (*domain.AuthData, errify.IError)
	RenewAuthorization(ctx context.Context, accessToken string, cfg config.TokenConfig) (string, errify.IError)
	Logout(ctx context.Context, accessToken string) errify.IError
}

type Cookies interface {
	SetToken(w http.ResponseWriter, token string)
	GetToken(r *http.Request) (string, error)
}

type Service struct {
	User
	Auth
	Cookies

	log logger.Logger
}

func NewService(
	log logger.Logger,
	pool *pgxpool.Pool,
	redisClient *redis.Client,
	repos *repository.Repository,
) *Service {
	transaction := repository.NewTransactionsRepos(pool, redisClient)

	return &Service{
		User:    NewUserService(log, transaction, repos),
		Auth:    NewAuthService(log, transaction, repos, repos),
		Cookies: NewCookiesService(),
		log:     log,
	}
}
