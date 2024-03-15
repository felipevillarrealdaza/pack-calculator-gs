package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/felipevillarrealdaza/go-service-template/internal/controller/viewmodel"
	"github.com/felipevillarrealdaza/go-service-template/internal/mediator"
	"github.com/felipevillarrealdaza/go-service-template/internal/mediator/domain_model"
	"github.com/google/uuid"
)

// Dependency injection using optional pattern
type OrderControllerDeps func(controller *orderController)

func WithOrderMediator(mediator mediator.OrderMediator) OrderControllerDeps {
	return func(controller *orderController) {
		controller.orderMediator = mediator
	}
}

type OrderController interface {
	AddOrder(w http.ResponseWriter, r *http.Request)
}

type orderController struct {
	orderMediator mediator.OrderMediator
}

func NewHttpOrderController(deps ...OrderControllerDeps) OrderController {
	orderController := orderController{}
	for _, opt := range deps {
		opt(&orderController)
	}
	return orderController
}

func (oc orderController) AddOrder(w http.ResponseWriter, r *http.Request) {
	// Parse request to viewmodel
	var requestBody viewmodel.OrderRequest
	jsonErr := json.NewDecoder(r.Body).Decode(&requestBody)
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusBadRequest)
		return
	}

	// Create domain model from viewmodel and create order
	order := domain_model.Order{
		OrderId:  uuid.New(),
		Quantity: requestBody.OrderQuantity,
	}
	if createOrderErr := oc.orderMediator.CreateOrder(r.Context(), order); createOrderErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint(createOrderErr)))
	}

	// Once order is created, calculate order packs needed
	orderPacks, calculateErr := oc.orderMediator.CalculateOrderPacks(r.Context(), order.OrderId)
	if calculateErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint(calculateErr)))
	}

	// Translate the order packs to view model and return to client
	response, marshalErr := json.Marshal(orderPacks.ToViewModel())
	if marshalErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprint(calculateErr)))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
