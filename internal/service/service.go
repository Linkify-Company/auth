package service

import (
	"auth/internal/config"
	"auth/internal/domain"
	"auth/internal/repository"
	"context"
	"github.com/Linkify-Company/common_utils/errify"
	"github.com/Linkify-Company/common_utils/logger"
	"github.com/go-redis/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

type User interface {
	AddUser(ctx context.Context, emailService Email, user *domain.User) (int, errify.IError)
	GetUserByID(ctx context.Context, id int) (*domain.User, errify.IError)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, errify.IError)
	PushCodeInEmail(ctx context.Context, emailService Email, email string) errify.IError
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

type Email interface {
	Send(ctx context.Context, title string, toEmail string, message string) errify.IError
}

type Service struct {
	User
	Auth
	Cookies
	Email

	log logger.Logger
}

func NewService(
	log logger.Logger,
	pool *pgxpool.Pool,
	redisClient *redis.Client,
	repos *repository.Repository,
	emailConfig *config.EmailServiceConfig,
) *Service {
	transaction := repository.NewTransactionsRepos(pool, redisClient)

	return &Service{
		User:    NewUserService(log, transaction, repos, repos),
		Auth:    NewAuthService(log, transaction, repos, repos, repos),
		Cookies: NewCookiesService(),
		Email:   NewEmailService(log, repos, *emailConfig),
		log:     log,
	}
}
