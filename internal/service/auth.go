package service

import (
	"auth/internal/config"
	"auth/internal/domain"
	"auth/internal/repository"
	"auth/pkg/errify"
	"auth/pkg/logger"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"os"
)

type AuthService struct {
	log         logger.Logger
	transaction repository.Transaction
	userRepos   repository.User
	authRepos   repository.Auth
}

func NewAuthService(
	log logger.Logger,
	transaction repository.Transaction,
	userRepos repository.User,
	authRepos repository.Auth,
) Auth {
	return &AuthService{
		log:         log,
		transaction: transaction,
		userRepos:   userRepos,
		authRepos:   authRepos,
	}
}

func (m *AuthService) Authorization(ctx context.Context, auth *domain.Auth, cfg config.TokenConfig) (string, errify.IError) {
	tx, err := m.transaction.Begin(ctx)
	if err != nil {
		return "", errify.NewInternalServerError(err.Error(), "Authorization/Begin")
	}
	user, err := m.userRepos.UserByEmail(ctx, tx, auth.Email)
	if err != nil {
		if errors.Is(err, repository.UserNotExist) {
			return "", errify.NewBadRequestError(err.Error(), UserNotExist.Error(), "Authorization/UserByEmail")
		}
		return "", errify.NewInternalServerError(err.Error(), "Authorization/UserByEmail")
	}
	err = bcrypt.CompareHashAndPassword(user.HashPassword, []byte(auth.Password))
	if err != nil {
		return "", errify.NewBadRequestError(err.Error(), ErrInvalidCredentials.Error(), "Authorization/CompareHashAndPassword")
	}

	redisTx, err := m.transaction.RedisTx(ctx)
	if err != nil {
		return "", errify.NewInternalServerError(err.Error(), "Authorization/RedisTx")
	}
	defer m.transaction.RedisRollback(ctx, redisTx)

	token, err := m.authRepos.Authorization(ctx, redisTx, &domain.AuthData{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}, cfg.RefreshTTL, cfg.AccessTTL, os.Getenv(config.Secret))
	if err != nil {
		return "", errify.NewInternalServerError(err.Error(), "Authorization/authRepos.Authorization")
	}
	err = m.transaction.RedisCommit(redisTx)
	if err != nil {
		return "", errify.NewInternalServerError(err.Error(), "Authorization/RedisCommit")
	}
	return token, nil
}

func (m *AuthService) CheckAuthorization(ctx context.Context, accessToken string) (*domain.AuthData, errify.IError) {
	redisClient := m.transaction.RedisClient(ctx)

	user, err := m.authRepos.CheckAuthorization(ctx, redisClient, accessToken)
	if err != nil {
		if errors.Is(err, repository.TokenExpired) {
			return nil, errify.NewUnauthorizedError(err.Error(), ErrTokenExpired.Error(), "CheckAuthorization/CheckAuthorization")
		}
		if errors.Is(err, repository.TokenNotValid) {
			return nil, errify.NewUnauthorizedError(err.Error(), ErrInvalidCredentials.Error(), "CheckAuthorization/CheckAuthorization")
		}
		if errors.Is(err, repository.TokenNotExist) {
			return nil, errify.NewUnauthorizedError(err.Error(), ErrInvalidCredentials.Error(), "CheckAuthorization/CheckAuthorization")
		}
		return nil, errify.NewInternalServerError(err.Error(), "CheckAuthorization/CheckAuthorization")
	}
	return user, nil
}

func (m *AuthService) RenewAuthorization(ctx context.Context, accessToken string, cfg config.TokenConfig) (string, errify.IError) {
	redisClient := m.transaction.RedisClient(ctx)

	token, err := m.authRepos.RenewalAuthorization(ctx, redisClient, accessToken, cfg.RefreshTTL, os.Getenv(config.Secret))
	if err != nil {
		if errors.Is(err, repository.TokenExpired) {
			return "", errify.NewBadRequestError(err.Error(), ErrTokenExpired.Error(), "RenewAuthorization/RenewalAuthorization")
		}
		if errors.Is(err, repository.TokenNotValid) {
			return "", errify.NewBadRequestError(err.Error(), ErrInvalidCredentials.Error(), "RenewAuthorization/RenewalAuthorization")
		}
		return "", errify.NewInternalServerError(err.Error(), "CheckAuthorization/RenewalAuthorization")
	}
	return token, nil
}

func (m *AuthService) Logout(ctx context.Context, accessToken string) errify.IError {
	tx, err := m.transaction.RedisTx(ctx)
	if err != nil {
		return errify.NewInternalServerError(err.Error(), "Logout/RedisTx")
	}
	defer m.transaction.RedisRollback(ctx, tx)

	err = m.authRepos.RemoveAuthorization(ctx, tx, accessToken)
	if err != nil {
		return errify.NewInternalServerError(err.Error(), "Logout/RemoveAuthorization")
	}
	err = m.transaction.RedisCommit(tx)
	if err != nil {
		return errify.NewInternalServerError(err.Error(), "Logout/RedisCommit")
	}
	return nil
}
