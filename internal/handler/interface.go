package handler

import (
	"github.com/Linkify-Company/common_utils/errify"
	"github.com/Linkify-Company/common_utils/logger"
	"github.com/gorilla/mux"
)

type Handler interface {
	Init(router *mux.Router)
}

func Run(log logger.Logger, handlers ...Handler) *mux.Router {
	var router = mux.NewRouter()

	router.Use()

	main := router.PathPrefix("/srv-auth").Subrouter()
	api := main.PathPrefix("/api").Subrouter()

	for _, h := range handlers {
		h.Init(api)
	}
	registeredEndpoints(router, log)

	return router
}

func registeredEndpoints(router *mux.Router, log logger.Logger) {
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
			log.Debugf("%s %s", methods, t)
		}
		return nil
	})
	if err != nil {
		log.Error(errify.NewInternalServerError(err.Error(), "registeredEndpoints/Walk"))
	}
}
