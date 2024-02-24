package service

import "errors"

var (
	UserIsAlreadyExist    = errors.New("user is already exist")
	UserNotExist          = errors.New("user is not exist")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenExpired       = errors.New("token expired")
)
