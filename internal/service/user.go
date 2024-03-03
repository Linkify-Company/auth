package service

import (
	"auth/internal/domain"
	"auth/internal/repository"
	"auth/pkg/html_template"
	"context"
	"errors"
	"fmt"
	"github.com/Linkify-Company/common_utils/errify"
	"github.com/Linkify-Company/common_utils/logger"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

type UserService struct {
	log         logger.Logger
	transaction repository.Transaction
	userRepos   repository.User
	emailRepos  repository.Email
}

func NewUserService(
	log logger.Logger,
	transaction repository.Transaction,
	userRepos repository.User,
	emailRepos repository.Email,
) User {
	return &UserService{
		log:         log,
		transaction: transaction,
		userRepos:   userRepos,
		emailRepos:  emailRepos,
	}
}

func (m *UserService) AddUser(ctx context.Context, emailService Email, user *domain.User) (int, errify.IError) {
	tx, err := m.transaction.Begin(ctx)
	if err != nil {
		return 0, errify.NewInternalServerError(err.Error(), "AddUser/Begin")
	}
	defer m.transaction.Rollback(ctx, tx)

	userExist, err := m.userRepos.UserByEmail(ctx, tx, user.Email)
	if err != nil && !errors.Is(err, repository.UserNotExist) {
		return 0, errify.NewInternalServerError(err.Error(), "AddUser/UserByEmail")
	}
	if userExist != nil {
		return 0, errify.NewBadRequestError(err.Error(), UserIsAlreadyExist.Error(), "AddUser/UserByEmail")
	}
	if !m.emailRepos.IsValid(ctx, user.Email, user.AuthorizationCode) {
		return 0, errify.NewBadRequestError(err.Error(), MailConfirmationError.Error(), "AddUser/IsValid")
	}

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

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		err = emailService.Send(ctx, "Успешная регистрация в Linkify", user.Email, fmt.Sprintf(html_template.RegistgrationSuccessfully, user.Email))
		if err != nil {
			m.log.Error(err.(errify.IError).JoinLoc("AddUser"))
		}
	}()

	return id, nil
}

func (m *UserService) PushCodeInEmail(ctx context.Context, emailService Email, email string) errify.IError {
	tx, err := m.transaction.Begin(ctx)
	if err != nil {
		return errify.NewInternalServerError(err.Error(), "PushCodeInEmail/Begin")
	}
	defer m.transaction.Rollback(ctx, tx)

	userExist, err := m.userRepos.UserByEmail(ctx, tx, email)
	if err != nil && !errors.Is(err, repository.UserNotExist) {
		return errify.NewInternalServerError(err.Error(), "PushCodeInEmail/UserByEmail")
	}
	if userExist != nil {
		return errify.NewBadRequestError(UserIsAlreadyExist.Error(), UserIsAlreadyExist.Error(), "PushCodeInEmail/UserByEmail")
	}
	var mins = 1000000
	var maxs = mins * 10

	var authorizationCode = rand.Intn(maxs-mins) + mins

	err = emailService.Send(ctx, "Подтверждение регистрации в Linkify", email, fmt.Sprintf(html_template.PushAuthCode, authorizationCode, authorizationCode))
	if err != nil {
		return err.(errify.IError).JoinLoc("PushCodeInEmail")
	}
	err = m.emailRepos.Set(ctx, email, authorizationCode)
	if err != nil {
		return err.(errify.IError).JoinLoc("PushCodeInEmail")
	}
	return nil
}

func (m *UserService) GetUserByID(ctx context.Context, id int) (*domain.User, errify.IError) {
	tx, err := m.transaction.Begin(ctx)
	if err != nil {
		return nil, errify.NewInternalServerError(err.Error(), "GetUserByID/Begin")
	}
	defer m.transaction.Rollback(ctx, tx)

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
	defer m.transaction.Rollback(ctx, tx)

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
