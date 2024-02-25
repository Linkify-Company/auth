package handler

import (
	"github.com/Linkify-Company/common_utils/errify"
	"github.com/Linkify-Company/common_utils/logger"
	"github.com/Linkify-Company/common_utils/response"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler interface {
	Init(router *mux.Router)
}

func Run(log logger.Logger, handlers ...Handler) *mux.Router {
	var router = mux.NewRouter()
	var hr = handler{log: log}

	router.Use()

	main := router.PathPrefix("/srv-auth").Subrouter()

	main.HandleFunc("/ping", hr.ping).Methods(http.MethodGet)

	api := main.PathPrefix("/api").Subrouter()

	for _, h := range handlers {
		h.Init(api)
	}
	hr.registeredEndpoints(router)

	return router
}

type handler struct {
	log logger.Logger
}

func (h *handler) registeredEndpoints(router *mux.Router) {
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		// Проверяем, имеет ли путь обработчик (является ли конечным)
		if route.GetHandler() != nil {
			t, err := route.GetPathTemplate()
			if err != nil {
				return err
			}
			methods, err := route.GetMethods()
			if err != nil {
				return err
			}
			h.log.Debugf("%s %s", methods, t)
		}
		return nil
	})
	if err != nil {
		h.log.Error(errify.NewInternalServerError(err.Error(), "registeredEndpoints/Walk"))
	}
}

func (h *handler) ping(w http.ResponseWriter, r *http.Request) {
	response.Ok(w, response.NewSend("", "pong", http.StatusOK), h.log)
}
