package service

import (
	"auth/internal/domain"
	"auth/internal/repository"
	"auth/pkg/errify"
	"auth/pkg/logger"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	log         logger.Logger
	transaction repository.Transaction
	userRepos   repository.User
}

func NewUserService(
	log logger.Logger,
	transaction repository.Transaction,
	userRepos repository.User,
) User {
	return &UserService{
		log:         log,
		transaction: transaction,
		userRepos:   userRepos,
	}
}

func (m *UserService) AddUser(ctx context.Context, user *domain.User) (int, errify.IError) {
	tx, err := m.transaction.Begin(ctx)
	if err != nil {
		return 0, errify.NewInternalServerError(err.Error(), "AddUser/Begin")
	}
	defer tx.Rollback(ctx)

	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, errify.NewBadRequestError(err.Error(), "password not valid", "AddUser/GenerateFromPassword")
	}
	user.Role.SetDefault()

	id, err := m.userRepos.AddUser(ctx, tx, user.Email, user.Role, passHash)
	if err != nil {
		if errors.Is(err, repository.UserAlreadyExist) {
			return 0, errify.NewBadRequestError(err.Error(), UserIsAlreadyExist.Error(), "AddUser/AddUser")
		}
		return 0, errify.NewInternalServerError(err.Error(), "AddUser/AddUser")
	}
	err = tx.Commit(ctx)
	if err != nil {
		return 0, errify.NewInternalServerError(err.Error(), "AddUser/Commit")
	}
	return id, nil
}

func (m *UserService) GetUserByID(ctx context.Context, id int) (*domain.User, errify.IError) {
	tx, err := m.transaction.Begin(ctx)
	if err != nil {
		return nil, errify.NewInternalServerError(err.Error(), "GetUserByID/Begin")
	}

	user, err := m.userRepos.UserById(ctx, tx, id)
	if err != nil {
		if errors.Is(err, repository.UserNotExist) {
			return nil, errify.NewBadRequestError(err.Error(), UserNotExist.Error(), "GetUserByID/UserById")
		}
		return nil, errify.NewInternalServerError(err.Error(), "GetUserByID/UserById")
	}
	return &domain.User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (m *UserService) GetUserByEmail(ctx context.Context, email string) (*domain.User, errify.IError) {
	tx, err := m.transaction.Begin(ctx)
	if err != nil {
		return nil, errify.NewInternalServerError(err.Error(), "GetUserByEmail/Begin")
	}

	user, err := m.userRepos.UserByEmail(ctx, tx, email)
	if err != nil {
		if errors.Is(err, repository.UserNotExist) {
			return nil, errify.NewBadRequestError(err.Error(), UserNotExist.Error(), "GetUserByEmail/UserByEmail")
		}
		return nil, errify.NewInternalServerError(err.Error(), "GetUserByEmail/UserByEmail")
	}
	return &domain.User{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}
