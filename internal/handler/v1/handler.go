package v1

import (
	"auth/internal/config"
	hr "auth/internal/handler"
	"auth/internal/service"
	"auth/pkg/errify"
	"auth/pkg/logger"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type handler struct {
	cfg      *config.HandlerConfig
	tokenCfg *config.TokenConfig
	log      logger.Logger
	service  *service.Service
}

func NewHandler(
	cfg *config.HandlerConfig,
	tokenCfg *config.TokenConfig,
	log logger.Logger,
	service *service.Service,
) hr.Handler {
	return &handler{
		cfg:      cfg,
		log:      log,
		service:  service,
		tokenCfg: tokenCfg,
	}
}

func (h *handler) Init(router *mux.Router) {
	h.log.Infof("Initialization handler V1")

	version := router.PathPrefix("/v1").Subrouter()
	version.Use(h.panicMiddleware)

	initUser(h, version)
	initAuth(h, version)
}

func (h *handler) panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				h.log.Error(errify.NewInternalServerError(fmt.Sprint(err), r.RequestURI).SetDetails("There was a panic in the router under version No. 1"))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
