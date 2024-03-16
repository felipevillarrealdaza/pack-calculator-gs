package controller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	api "github.com/felipevillarrealdaza/go-service-template/internal/api/http"
	"github.com/felipevillarrealdaza/go-service-template/internal/controller/viewmodel"
	"github.com/felipevillarrealdaza/go-service-template/internal/mediator/domain_model"
	mediator_mocks "github.com/felipevillarrealdaza/go-service-template/internal/mediator/mocks"
	repository_mocks "github.com/felipevillarrealdaza/go-service-template/internal/repository/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_AddOrder_OK(t *testing.T) {
	// Set Up
	orderMediatorMock := mediator_mocks.NewOrderMediator(t)
	packMediatorMock := mediator_mocks.NewPackMediator(t)
	repositoryMock := repository_mocks.NewQuerier(t)
	router := api.NewRouter(packMediatorMock, orderMediatorMock, repositoryMock)
	httpRecorder := httptest.NewRecorder()

	// Arrange
	reqBody := viewmodel.OrderRequest{
		OrderQuantity: 500,
	}
	requestBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/order", bytes.NewBuffer(requestBytes))
	orderMediatorMock.On("CreateOrder", mock.Anything, mock.Anything).Return(nil)
	orderMediatorMock.On("CalculateOrderPacks", mock.Anything, mock.Anything).Return(domain_model.OrderPacks{OptimalOrderPack: domain_model.OrderPack(map[int]int{2: 2})}, nil)

	// Act
	router.ServeHTTP(httpRecorder, req)

	// Assert
	orderMediatorMock.AssertExpectations(t)
	require.Equal(t, http.StatusCreated, httpRecorder.Code)

	// Clean up
	orderMediatorMock.ExpectedCalls = make([]*mock.Call, 0)
}

func Test_AddOrder_Errors(t *testing.T) {
	// Set up
	orderMediatorMock := mediator_mocks.NewOrderMediator(t)
	packMediatorMock := mediator_mocks.NewPackMediator(t)
	repositoryMock := repository_mocks.NewQuerier(t)
	router := api.NewRouter(packMediatorMock, orderMediatorMock, repositoryMock)

	t.Run("Wrong JSON body", func(t *testing.T) {
		// Arrange
		httpRecorder := httptest.NewRecorder()
		reqBody := struct{ Order int }{
			Order: 2,
		}
		requestBytes, _ := json.Marshal(reqBody)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/order", bytes.NewBuffer(requestBytes))

		// Act
		router.ServeHTTP(httpRecorder, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, httpRecorder.Code)
		orderMediatorMock.AssertExpectations(t)

		// Clean up
		orderMediatorMock.ExpectedCalls = make([]*mock.Call, 0)
	})

	t.Run("Unparsable JSON", func(t *testing.T) {
		// Arrange
		httpRecorder := httptest.NewRecorder()
		requestBytes := []byte("{quantity 2")
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/order", bytes.NewBuffer(requestBytes))

		// Act
		router.ServeHTTP(httpRecorder, req)

		// Assert
		require.Equal(t, http.StatusUnprocessableEntity, httpRecorder.Code)
		orderMediatorMock.AssertExpectations(t)

		// Clean up
		orderMediatorMock.ExpectedCalls = make([]*mock.Call, 0)
	})

	t.Run("Unknown error creating order", func(t *testing.T) {
		// Arrange
		httpRecorder := httptest.NewRecorder()
		reqBody := viewmodel.OrderRequest{
			OrderQuantity: 2,
		}
		requestBytes, _ := json.Marshal(reqBody)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/order", bytes.NewBuffer(requestBytes))
		orderMediatorMock.On("CreateOrder", mock.Anything, mock.Anything).Return(errors.New("unknown error creating order"))

		// Act
		router.ServeHTTP(httpRecorder, req)

		// Assert
		require.Equal(t, http.StatusInternalServerError, httpRecorder.Code)
		orderMediatorMock.AssertExpectations(t)

		// Clean up
		orderMediatorMock.ExpectedCalls = make([]*mock.Call, 0)
	})

	t.Run("Unknown error calculating packs for order", func(t *testing.T) {
		// Arrange
		httpRecorder := httptest.NewRecorder()
		reqBody := viewmodel.OrderRequest{
			OrderQuantity: 2,
		}
		requestBytes, _ := json.Marshal(reqBody)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/order", bytes.NewBuffer(requestBytes))
		orderMediatorMock.On("CreateOrder", mock.Anything, mock.Anything).Return(nil)
		orderMediatorMock.On("CalculateOrderPacks", mock.Anything, mock.Anything).Return(domain_model.OrderPacks{}, errors.New("unknown error calculating packs"))

		// Act
		router.ServeHTTP(httpRecorder, req)

		// Assert
		require.Equal(t, http.StatusInternalServerError, httpRecorder.Code)
		orderMediatorMock.AssertExpectations(t)

		// Clean up
		orderMediatorMock.ExpectedCalls = make([]*mock.Call, 0)
	})
}

func Test_AddOrder_MethodsNotAllowed(t *testing.T) {
	// Arrange
	orderMediatorMock := mediator_mocks.NewOrderMediator(t)
	packMediatorMock := mediator_mocks.NewPackMediator(t)
	repositoryMock := repository_mocks.NewQuerier(t)
	router := api.NewRouter(packMediatorMock, orderMediatorMock, repositoryMock)
	httpRecorder := httptest.NewRecorder()
	reqBody := viewmodel.OrderRequest{
		OrderQuantity: 500,
	}
	requestBytes, _ := json.Marshal(reqBody)
	t.Run("Methods not implemented", func(t *testing.T) {
		for _, httpVerb := range []string{http.MethodGet, http.MethodPut, http.MethodDelete, http.MethodHead, http.MethodPatch, http.MethodOptions} {
			// Arrange
			req, _ := http.NewRequestWithContext(context.Background(), httpVerb, "/api/v1/order", bytes.NewBuffer(requestBytes))

			// Act
			router.ServeHTTP(httpRecorder, req)

			// Assert
			require.Equal(t, http.StatusMethodNotAllowed, httpRecorder.Code)
			orderMediatorMock.AssertExpectations(t)

			// Clean up
			orderMediatorMock.ExpectedCalls = make([]*mock.Call, 0)
		}
	})
}
