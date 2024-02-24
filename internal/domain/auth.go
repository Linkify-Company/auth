package domain

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"regexp"
	"unicode/utf8"
)

type Auth struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"`
}

func (u *Auth) Valid() error {
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

func (u *Auth) validPassword() error {
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
