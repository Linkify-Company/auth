package v1

import (
	"auth/internal/domain"
	hr "auth/internal/handler"
	"context"
	"encoding/json"
	"github.com/Linkify-Company/common_utils/errify"
	"github.com/Linkify-Company/common_utils/response"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func initUser(h *handler, router *mux.Router) {
	user := router.PathPrefix("/user").Subrouter()
	user.HandleFunc("", h.AddUser).Methods(http.MethodPost)

	user.HandleFunc("/{value}", h.GetUser).Methods(http.MethodGet)
}

func (h *handler) AddUser(w http.ResponseWriter, r *http.Request) {
	var req domain.User
	e := json.NewDecoder(r.Body).Decode(&req)
	if e != nil {
		response.Error(w, errify.NewBadRequestError(e.Error(), hr.ValidationError, "AddUser").
			JoinLoc("NewDecoder"), h.log)
		return
	}
	e = req.Valid()
	if e != nil {
		response.Error(w, errify.NewBadRequestError(e.Error(), hr.ValidationError, "AddUser").
			JoinLoc("Valid"), h.log)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.ContextTimeout)
	defer cancel()

	id, err := h.service.AddUser(ctx, h.service.Email, &req)
	if err != nil {
		response.Error(w, err.JoinLoc("AddUser"), h.log)
		return
	}
	response.Ok(w, response.NewSend(id, "Create user successfully", http.StatusCreated), h.log)
}

func (h *handler) GetUser(w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["value"]
	id, _ := strconv.Atoi(mux.Vars(r)["value"])
	if email == "" && id <= 0 {
		response.Error(w, errify.NewBadRequestError(hr.ValidationError, hr.ValidationError, "GetUser").
			JoinLoc("Atoi"), h.log)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.ContextTimeout)
	defer cancel()

	var err errify.IError
	var user *domain.User

	if id > 0 {
		user, err = h.service.GetUserByID(ctx, id)
	} else {
		user, err = h.service.GetUserByEmail(ctx, email)
	}
	if err != nil {
		response.Error(w, err.JoinLoc("GetUser"), h.log)
		return
	}
	response.Ok(w, response.NewSend(user, "Get user successfully", http.StatusOK), h.log)
}
