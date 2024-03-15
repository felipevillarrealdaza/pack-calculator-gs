package controller

import (
	"net/http"
)

type HealthControllerDeps func(controller *healthController)

type HealthController interface {
	Health(w http.ResponseWriter, r *http.Request)
}

type healthController struct {
}

func NewHttpHealthController(deps ...HealthControllerDeps) HealthController {
	controller := healthController{}
	for _, opt := range deps {
		opt(&controller)
	}
	return controller
}

func (hc healthController) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Health endpoint called!"))
}
