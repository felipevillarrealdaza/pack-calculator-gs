package controller

import (
	"encoding/json"
	"net/http"

	"github.com/felipevillarrealdaza/go-service-template/internal/controller/viewmodel"
	"github.com/felipevillarrealdaza/go-service-template/internal/mediator"
)

// Dependency injection using optional pattern
type PackControllerDeps func(controller *packController)

func WithPackMediator(mediator mediator.PackMediator) PackControllerDeps {
	return func(controller *packController) {
		controller.packMediator = mediator
	}
}

type PackController interface {
	AddPack(w http.ResponseWriter, r *http.Request)
	RemovePack(w http.ResponseWriter, r *http.Request)
}

type packController struct {
	packMediator mediator.PackMediator
}

func NewHttpPackController(deps ...PackControllerDeps) PackController {
	packController := packController{}
	for _, opt := range deps {
		opt(&packController)
	}
	return packController
}

func (pc packController) AddPack(w http.ResponseWriter, r *http.Request) {
	// Parse request to viewmodel
	var requestBody viewmodel.Pack
	jsonErr := json.NewDecoder(r.Body).Decode(&requestBody)
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusBadRequest)
		return
	}

	if addPackErr := pc.packMediator.AddPack(r.Context(), requestBody.Size); addPackErr != nil {
		http.Error(w, addPackErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(""))
}

func (pc packController) RemovePack(w http.ResponseWriter, r *http.Request) {
	// Parse request to viewmodel
	var requestBody viewmodel.Pack
	jsonErr := json.NewDecoder(r.Body).Decode(&requestBody)
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusBadRequest)
		return
	}

	// Remove pack
	if addPackErr := pc.packMediator.RemovePack(r.Context(), requestBody.Size); addPackErr != nil {
		http.Error(w, addPackErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(""))
}
