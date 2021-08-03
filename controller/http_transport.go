package controller

import (
	"fmt"
	"net/http"
	"stratum-server/service"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

const (
	apiResource    = "api"
	v1Resource     = "v1"
	healthResource = "health"
	wsResource     = "ws"
)

var (
	healthEndpoint = fmt.Sprintf("/%s/%s/%s", apiResource, v1Resource, healthResource)
	wsEndpoint     = fmt.Sprintf("/%s/%s/%s", apiResource, v1Resource, wsResource)
)

// NewHandler: create handlers
func NewHandler(svc service.Service) http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(middleware.Recoverer, middleware.StripSlashes, middleware.Logger)

		r.Get(healthEndpoint, health(svc))
		r.Get(wsEndpoint, ws(svc))

	})

	return r
}
