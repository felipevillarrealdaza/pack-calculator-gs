package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/felipevillarrealdaza/go-service-template/internal/controller/viewmodel"
	"github.com/felipevillarrealdaza/go-service-template/internal/mediator"
	"github.com/felipevillarrealdaza/go-service-template/internal/mediator/domain_model"
	"github.com/go-playground/validator/v10"
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
	validate      *validator.Validate
}

func NewHttpOrderController(deps ...OrderControllerDeps) OrderController {
	orderController := orderController{validate: validator.New(validator.WithRequiredStructEnabled())}
	for _, opt := range deps {
		opt(&orderController)
	}
	return orderController
}

func (oc orderController) AddOrder(w http.ResponseWriter, r *http.Request) {
	var requestBody viewmodel.OrderRequest

	// Validate JSON and request body
	jsonErr := json.NewDecoder(r.Body).Decode(&requestBody)
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusUnprocessableEntity)
		return
	}
	validationErr := oc.validate.Struct(&requestBody)
	if validationErr != nil {
		http.Error(w, validationErr.Error(), http.StatusBadRequest)
		return
	}

	// Create domain model from viewmodel and create order
	order := domain_model.Order{
		OrderId:  uuid.New(),
		Quantity: requestBody.OrderQuantity,
	}
	if createOrderErr := oc.orderMediator.CreateOrder(r.Context(), order); createOrderErr != nil {
		http.Error(w, createOrderErr.Error(), http.StatusInternalServerError)
		return
	}

	// Once order is created, calculate order packs needed
	orderPacks, calculateErr := oc.orderMediator.CalculateOrderPacks(r.Context(), order.OrderId)
	if calculateErr != nil {
		http.Error(w, calculateErr.Error(), http.StatusInternalServerError)
		return
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
