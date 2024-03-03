package v1

import (
	hr "auth/internal/handler"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Linkify-Company/common_utils/errify"
	"github.com/Linkify-Company/common_utils/response"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
)

func initEmail(h *handler, router *mux.Router) {
	email := router.PathPrefix("/email").Subrouter()
	email.HandleFunc("/push_auth", h.PushCodeInEmail).Methods(http.MethodPost)
}

func (h *handler) PushCodeInEmail(w http.ResponseWriter, r *http.Request) {
	var email = struct {
		Email string `json:"email"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&email)
	if err != nil {
		response.Error(w, errify.NewBadRequestError(err.Error(), hr.ValidationError, "PushCodeInEmail/Decode"), h.log)
		return
	}
	err = validator.New().Var(email.Email, "required,email")
	if err != nil {
		response.Error(w, errify.NewBadRequestError(err.Error(), hr.ValidationError, "PushCodeInEmail/Var"), h.log)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.ContextTimeout)
	defer cancel()

	e := h.service.User.PushCodeInEmail(ctx, h.service.Email, email.Email)
	if e != nil {
		response.Error(w, e.JoinLoc("PushCodeInEmail"), h.log)
		return
	}
	response.Ok(w, response.NewSend("", fmt.Sprintf("the confirmation code has been sent to the email: %s", email.Email), http.StatusOK), h.log)
}
