package v1

import (
	"auth/internal/domain"
	hr "auth/internal/handler"
	"auth/internal/service"
	"auth/pkg/errify"
	"auth/pkg/response"
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type CheckAuthResponse struct {
	User  *domain.AuthData `json:"user,omitempty"`
	Token string           `json:"token,omitempty"`
}

func initAuth(h *handler, router *mux.Router) {
	auth := router.PathPrefix("/auth").Subrouter()

	auth.HandleFunc("/login", h.Login).Methods(http.MethodPost)
	auth.HandleFunc("/logout", h.Logout).Methods(http.MethodDelete)
	auth.HandleFunc("/check", h.CheckAuth).Methods(http.MethodGet)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var req domain.Auth
	e := json.NewDecoder(r.Body).Decode(&req)
	if e != nil {
		response.Error(w, errify.NewBadRequestError(e.Error(), hr.ValidationError, "Login").
			JoinLoc("NewDecoder"), h.log)
		return
	}
	e = req.Valid()
	if e != nil {
		response.Error(w, errify.NewBadRequestError(e.Error(), hr.ValidationError, "Login").
			JoinLoc("Valid"), h.log)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.ContextTimeout)
	defer cancel()

	token, err := h.service.Authorization(ctx, &req, *h.tokenCfg)
	if err != nil {
		response.Error(w, err.JoinLoc("Authorization"), h.log)
		return
	}
	h.service.SetToken(w, token)

	response.Ok(w, response.NewSend(token, "Authorization successfully", http.StatusOK), h.log)
}

func (h *handler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	req, e := h.service.GetToken(r)
	if e != nil || req == "" {
		response.Error(w, errify.NewBadRequestError(e.Error(),
			service.ErrInvalidCredentials.Error(), "CheckAuth"), h.log)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.ContextTimeout)
	defer cancel()

	user, err := h.service.CheckAuthorization(ctx, req)
	if err != nil {
		if errors.Is(err, service.ErrTokenExpired) {
			token, err := h.service.RenewAuthorization(ctx, req, *h.tokenCfg)
			if err != nil {
				if err, ok := err.(*errify.InternalServerError); ok {
					response.Error(w, err.JoinLoc("CheckAuth"), h.log)
					return
				}
				response.Error(w, errify.NewInternalServerError(err.Error(), "CheckAuth").
					JoinLoc("RenewAuthorization").JoinLoc(err.Location()), h.log)
				return
			}
			user, err := h.service.CheckAuthorization(ctx, req)
			if err != nil {
				if err, ok := err.(*errify.InternalServerError); ok {
					response.Error(w, err.JoinLoc("CheckAuth"), h.log)
					return
				}
				response.Error(w, errify.NewInternalServerError(err.Error(), "CheckAuth").
					JoinLoc("RenewAuthorization").JoinLoc(err.Location()), h.log)
				return
			}
			response.Ok(w, response.NewSend(CheckAuthResponse{
				User:  user,
				Token: token,
			}, "Authorization successfully", http.StatusOK), h.log)
			return
		}
		response.Error(w, err.JoinLoc("CheckAuth"), h.log)
		return
	}
	response.Ok(w, response.NewSend(CheckAuthResponse{
		User: user,
	}, "Authorization successfully", http.StatusOK), h.log)
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	req, e := h.service.GetToken(r)
	if e != nil || req == "" {
		response.Error(w, errify.NewBadRequestError(e.Error(),
			service.ErrInvalidCredentials.Error(), "CheckAuth"), h.log)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.ContextTimeout)
	defer cancel()

	err := h.service.Logout(ctx, req)
	if err != nil {
		response.Error(w, err.JoinLoc("Logout"), h.log)
		return
	}
	response.Ok(w, response.NewSend("", "Logout successfully", http.StatusOK), h.log)
}
