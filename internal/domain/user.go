package domain

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"regexp"
	"unicode/utf8"
)

// User Сущность пользователя
type User struct {
	// Идентификатор пользователя
	ID int `json:"id"`
	// Email пользователя
	Email string `json:"email" validate:"required,email"`
	// Password пользователя
	Password string `json:"password"`
	// AuthorizationCode Авторизационный код для подтверждения почты (email)
	AuthorizationCode int `json:"authorization_code" validate:"required"`

	Role Role `json:"-"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) Valid() error {
	if u == nil {
		return errors.New("user empty")
	}
	err := u.validPassword()
	if err != nil {
		return err
	}
	vl := validator.New()
	err = vl.Struct(*u)
	if err != nil {
		return err.(validator.ValidationErrors)[0]
	}
	return nil
}

func (u *User) validPassword() error {
	length := utf8.RuneCountInString(u.Password)
	if length < 4 || length > 25 {
		return errors.New("password not valid")
	}

	ok, err := regexp.MatchString("^[A-Za-z0-9А-Яа-я]+$", u.Password)
	if !ok || err != nil {
		return errors.New("password not valid")
	}
	return nil
}

type UserFromDB struct {
	ID           int
	Email        string
	HashPassword []byte
	Role         Role
}
