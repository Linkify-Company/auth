package repository

import "errors"

var (
	UserAlreadyExist   = errors.New("user already exists")
	UserNotExist       = errors.New("user not exists")
	TokenNotExist      = errors.New("token not exists")
	TokenExpired       = errors.New("token expired")
	TokenNotValid      = errors.New("token not valid")
	TokenInvalidClaims = errors.New("token invalid claims")
)
