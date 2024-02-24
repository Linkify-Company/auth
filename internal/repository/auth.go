package repository

import (
	"auth/internal/config"
	"auth/internal/domain"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type AuthRepo struct{}

func NewAuthRepo() Auth {
	return &AuthRepo{}
}

const (
	accessKey  = "access_key:"
	refreshKey = "refresh_key:"

	userID    = "uid"
	userEmail = "email"
	userRole  = "role"
	userExp   = "exp"
)

func (m *AuthRepo) Authorization(ctx context.Context, tx redis.Pipeliner, user *domain.AuthData, refreshTTL time.Duration, accessTTL time.Duration, secret string) (string, error) {
	accessToken, err := newToken(*user, accessTTL, secret)
	if err != nil {
		return "", fmt.Errorf("Authorization/newToken: %w", err)
	}
	refreshToken, err := newToken(*user, refreshTTL, secret)
	if err != nil {
		return "", fmt.Errorf("Authorization/newToken: %w", err)
	}
	// Устанавливаем у access токена такое же время, как и у refresh, чтобы access не удалился с redis раньше нужного
	status := tx.Set(fmt.Sprint(accessKey, user.ID), accessToken, refreshTTL)
	if status.Err() != nil {
		return "", fmt.Errorf("Authorization/Set: %w", status.Err())
	}
	status = tx.Set(fmt.Sprint(refreshKey, user.ID), refreshToken, refreshTTL)
	if status.Err() != nil {
		return "", fmt.Errorf("Authorization/Set: %w", status.Err())
	}
	return accessToken, nil
}

func (m *AuthRepo) RenewalAuthorization(ctx context.Context, redisClient *redis.Client, accessToken string, accessTTL time.Duration, secret string) (string, error) {
	user, err := userData(accessToken)
	if err != nil {
		return "", err
	}

	// Проверяем, валидный ли access
	err = validAccess(redisClient, accessToken, user.ID)
	if err != nil {
		return "", err
	}

	// Если да, то удаляем его и получаем refresh токен, и если нет ошибок, то создаем новый access и отправляем
	redisClient.Del(fmt.Sprint(accessKey, userID))

	user, err = getUserData(redisClient, refreshKey, user.ID)
	if err != nil {
		return "", err
	}

	token, err := newToken(*user, accessTTL, secret)
	if err != nil {
		return "", fmt.Errorf("RenewalAuthorization/newToken: %w", err)
	}
	// Сохраняем в redis
	status := redisClient.Set(fmt.Sprint(accessKey, user.ID), token, accessTTL)
	if status.Err() != nil {
		return "", fmt.Errorf("RenewalAuthorization/Set: %w", err)
	}
	return token, nil
}

func (m *AuthRepo) CheckAuthorization(ctx context.Context, redisClient *redis.Client, accessToken string) (*domain.AuthData, error) {
	user, err := userData(accessToken)
	if err != nil {
		return nil, err
	}
	err = validAccess(redisClient, accessToken, user.ID)
	if err != nil {
		return nil, err
	}
	user, err = getUserData(redisClient, refreshKey, user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (m *AuthRepo) RemoveAuthorization(ctx context.Context, tx redis.Pipeliner, accessToken string) error {
	user, err := userData(accessToken)
	if err != nil {
		return err
	}
	inc := tx.Del(fmt.Sprint(accessKey, user.ID))
	if inc.Err() != nil {
		return fmt.Errorf("RemoveAuthorization/Del: %w", inc.Err())
	}
	inc = tx.Del(fmt.Sprint(refreshKey, user.ID))
	if inc.Err() != nil {
		return fmt.Errorf("RemoveAuthorization/Del: %w", inc.Err())
	}
	return nil
}

func validAccess(redisClient *redis.Client, token string, userID int) error {
	// Проверяем, есть ли access в redis
	s := redisClient.Get(fmt.Sprint(accessKey, userID))
	if s.Err() != nil {
		if errors.Is(s.Err(), redis.Nil) {
			// Ключ не найден в Redis
			return TokenNotExist
		}
		return fmt.Errorf("validAccess: %w", s.Err())
	}
	var val string
	err := s.Scan(&val)
	if err != nil {
		return fmt.Errorf("validAccess/Scan: %w", err)
	}
	// проверяем, это такой же токен, что и в redis
	if val == "" || val != token {
		return TokenNotValid
	}
	return nil
}

func newToken(user domain.AuthData, duration time.Duration, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		userID:    user.ID,
		userEmail: user.Email,
		userRole:  user.Role,
		userExp:   time.Now().Add(duration).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func getUserData(redisClient *redis.Client, key string, userID int) (*domain.AuthData, error) {
	s := redisClient.Get(fmt.Sprint(key, userID))
	user, err := userData(s.Val())
	return user, err
}

func userData(token string) (user *domain.AuthData, err error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv(config.Secret)), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			claims, ok := t.Claims.(jwt.MapClaims)
			if !ok {
				return nil, TokenInvalidClaims
			}

			return &domain.AuthData{
				ID: int(claims[userID].(float64)),
			}, TokenExpired
		}
		return nil, TokenNotValid
	}
	if !t.Valid {
		return nil, TokenNotValid
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, TokenInvalidClaims
	}

	return &domain.AuthData{
		ID:    int(claims[userID].(float64)),
		Email: fmt.Sprint(claims[userEmail]),
		Role:  domain.Role(int(claims[userRole].(float64))),
	}, nil
}
