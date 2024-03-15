package api

import (
	"database/sql"
	"net/http"

	"github.com/felipevillarrealdaza/go-service-template/internal/controller"
	"github.com/felipevillarrealdaza/go-service-template/internal/mediator"
	"github.com/felipevillarrealdaza/go-service-template/internal/repository"
	"github.com/gorilla/mux"
)

func NewRouter(dbCtx *sql.DB) http.Handler {
	router := mux.NewRouter().PathPrefix("/api/v1").Subrouter()

	// Add middlewares for the router

	// Create repository, which is a dependency for mediators
	repository := repository.New(dbCtx)

	// Create mediators, which are dependencies for controllers
	packMediator := mediator.NewPackMediator(mediator.WithPackRepository(repository))
	orderMediator := mediator.NewOrderMediator(mediator.WithOrderRepository(repository))

	// Create controllers
	healthController := controller.NewHttpHealthController()
	orderController := controller.NewHttpOrderController(controller.WithOrderMediator(orderMediator))
	packController := controller.NewHttpPackController(controller.WithPackMediator(packMediator))

	// Match routes to controller's methods
	router.Path("/health").Methods(http.MethodGet).HandlerFunc(healthController.Health)
	router.Path("/pack").Methods(http.MethodPost).HandlerFunc(packController.AddPack)
	router.Path("/pack").Methods(http.MethodDelete).HandlerFunc(packController.RemovePack)
	router.Path("/order").Methods(http.MethodPost).HandlerFunc(orderController.AddOrder)

	return router
}
