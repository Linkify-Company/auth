package service

import (
	"net/http"
	"time"
)

type CookiesService struct{}

func NewCookiesService() Cookies {
	return &CookiesService{}
}

const (
	Authorization = "Authorization"
	TokenNil      = ""
)

func (m *CookiesService) SetToken(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     Authorization,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Domain:   "localhost",
		Expires:  time.Now().Add(365 * 24 * time.Hour),
	})
}

func (m *CookiesService) GetToken(r *http.Request) (string, error) {
	token := r.Header.Get(Authorization)
	if token != "" {
		return token, nil
	}

	c, err := r.Cookie(Authorization)
	if err != nil {
		return "", err
	}
	return c.Value, nil
}
