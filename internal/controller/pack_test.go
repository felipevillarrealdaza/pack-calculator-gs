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
	mediator_mocks "github.com/felipevillarrealdaza/go-service-template/internal/mediator/mocks"
	repository_mocks "github.com/felipevillarrealdaza/go-service-template/internal/repository/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_AddPack_OK(t *testing.T) {
	// Arrange
	orderMediatorMock := mediator_mocks.NewOrderMediator(t)
	packMediatorMock := mediator_mocks.NewPackMediator(t)
	repositoryMock := repository_mocks.NewQuerier(t)
	router := api.NewRouter(packMediatorMock, orderMediatorMock, repositoryMock)
	httpRecorder := httptest.NewRecorder()
	reqBody := viewmodel.PackRequest{
		Size: 2,
	}
	requestBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/pack", bytes.NewBuffer(requestBytes))
	packMediatorMock.On("AddPack", mock.Anything, reqBody.Size).Return(nil)

	// Act
	router.ServeHTTP(httpRecorder, req)

	// Assert
	packMediatorMock.AssertExpectations(t)
	require.Equal(t, http.StatusCreated, httpRecorder.Code)

	// Clean up
	packMediatorMock.ExpectedCalls = make([]*mock.Call, 0)
}

func Test_AddPack_Errors(t *testing.T) {
	// Set up
	orderMediatorMock := mediator_mocks.NewOrderMediator(t)
	packMediatorMock := mediator_mocks.NewPackMediator(t)
	repositoryMock := repository_mocks.NewQuerier(t)
	router := api.NewRouter(packMediatorMock, orderMediatorMock, repositoryMock)

	t.Run("Wrong JSON body", func(t *testing.T) {
		// Arrange
		httpRecorder := httptest.NewRecorder()
		reqBody := struct{ PackSize int }{
			PackSize: 2,
		}
		requestBytes, _ := json.Marshal(reqBody)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/pack", bytes.NewBuffer(requestBytes))

		// Act
		router.ServeHTTP(httpRecorder, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, httpRecorder.Code)
		packMediatorMock.AssertExpectations(t)

		// Clean up
		packMediatorMock.ExpectedCalls = make([]*mock.Call, 0)
	})

	t.Run("Unparsable JSON", func(t *testing.T) {
		// Arrange
		httpRecorder := httptest.NewRecorder()
		requestBytes := []byte("{size: 2")
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/pack", bytes.NewBuffer(requestBytes))

		// Act
		router.ServeHTTP(httpRecorder, req)

		// Assert
		require.Equal(t, http.StatusUnprocessableEntity, httpRecorder.Code)
		packMediatorMock.AssertExpectations(t)

		// Clean up
		packMediatorMock.ExpectedCalls = make([]*mock.Call, 0)
	})

	t.Run("Pack already exists", func(t *testing.T) {
		// Arrange
		httpRecorder := httptest.NewRecorder()
		reqBody := viewmodel.PackRequest{
			Size: 2,
		}
		requestBytes, _ := json.Marshal(reqBody)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/api/v1/pack", bytes.NewBuffer(requestBytes))
		packMediatorMock.On("AddPack", mock.Anything, reqBody.Size).Return(errors.New("Pack already exists!"))

		// Act
		router.ServeHTTP(httpRecorder, req)

		// Assert
		require.Equal(t, http.StatusConflict, httpRecorder.Code)
		packMediatorMock.AssertExpectations(t)

		// Clean up
		packMediatorMock.ExpectedCalls = make([]*mock.Call, 0)
	})
}

func Test_MethodsNotAllowed(t *testing.T) {
	// Set up
	orderMediatorMock := mediator_mocks.NewOrderMediator(t)
	packMediatorMock := mediator_mocks.NewPackMediator(t)
	repositoryMock := repository_mocks.NewQuerier(t)
	router := api.NewRouter(packMediatorMock, orderMediatorMock, repositoryMock)
	httpRecorder := httptest.NewRecorder()

	t.Run("Methods not implemented", func(t *testing.T) {
		for _, httpVerb := range []string{http.MethodGet, http.MethodPut, http.MethodHead, http.MethodPatch, http.MethodOptions} {
			// Arrange
			reqBody := struct{ PackSize int }{
				PackSize: 2,
			}
			requestBytes, _ := json.Marshal(reqBody)
			req, _ := http.NewRequestWithContext(context.Background(), httpVerb, "/api/v1/pack", bytes.NewBuffer(requestBytes))

			// Act
			router.ServeHTTP(httpRecorder, req)

			// Assert
			require.Equal(t, http.StatusNotFound, httpRecorder.Code)
			packMediatorMock.AssertExpectations(t)

			// Clean up
			packMediatorMock.ExpectedCalls = make([]*mock.Call, 0)
		}
	})
}

func Test_RemovePack_OK(t *testing.T) {
	// Arrange
	orderMediatorMock := mediator_mocks.NewOrderMediator(t)
	packMediatorMock := mediator_mocks.NewPackMediator(t)
	repositoryMock := repository_mocks.NewQuerier(t)
	router := api.NewRouter(packMediatorMock, orderMediatorMock, repositoryMock)
	httpRecorder := httptest.NewRecorder()
	reqBody := viewmodel.PackRequest{
		Size: 2,
	}
	requestBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/api/v1/pack", bytes.NewBuffer(requestBytes))
	packMediatorMock.On("RemovePack", mock.Anything, reqBody.Size).Return(nil)
	// Act
	router.ServeHTTP(httpRecorder, req)

	// Assert
	orderMediatorMock.AssertExpectations(t)
	require.Equal(t, http.StatusOK, httpRecorder.Code)

	// Clean up
	packMediatorMock.ExpectedCalls = make([]*mock.Call, 0)
}

func Test_RemovePack_Errors(t *testing.T) {
	// Set up
	orderMediatorMock := mediator_mocks.NewOrderMediator(t)
	packMediatorMock := mediator_mocks.NewPackMediator(t)
	repositoryMock := repository_mocks.NewQuerier(t)
	router := api.NewRouter(packMediatorMock, orderMediatorMock, repositoryMock)

	t.Run("Wrong JSON body", func(t *testing.T) {
		// Arrange
		httpRecorder := httptest.NewRecorder()
		reqBody := struct{ PackSize int }{
			PackSize: 2,
		}
		requestBytes, _ := json.Marshal(reqBody)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/api/v1/pack", bytes.NewBuffer(requestBytes))

		// Act
		router.ServeHTTP(httpRecorder, req)

		// Assert
		require.Equal(t, http.StatusBadRequest, httpRecorder.Code)
		packMediatorMock.AssertExpectations(t)

		// Clean up
		packMediatorMock.ExpectedCalls = make([]*mock.Call, 0)
	})

	t.Run("Unparsable JSON", func(t *testing.T) {
		// Arrange
		httpRecorder := httptest.NewRecorder()
		requestBytes := []byte("{size: 2")
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/api/v1/pack", bytes.NewBuffer(requestBytes))

		// Act
		router.ServeHTTP(httpRecorder, req)

		// Assert
		require.Equal(t, http.StatusUnprocessableEntity, httpRecorder.Code)
		packMediatorMock.AssertExpectations(t)

		// Clean up
		packMediatorMock.ExpectedCalls = make([]*mock.Call, 0)
	})

	t.Run("Unknown error", func(t *testing.T) {
		// Arrange
		httpRecorder := httptest.NewRecorder()
		reqBody := viewmodel.PackRequest{
			Size: 2,
		}
		requestBytes, _ := json.Marshal(reqBody)
		req, _ := http.NewRequestWithContext(context.Background(), http.MethodDelete, "/api/v1/pack", bytes.NewBuffer(requestBytes))
		packMediatorMock.On("RemovePack", mock.Anything, reqBody.Size).Return(errors.New("unexpected error happened"))

		// Act
		router.ServeHTTP(httpRecorder, req)

		// Assert
		require.Equal(t, http.StatusInternalServerError, httpRecorder.Code)
		packMediatorMock.AssertExpectations(t)

		// Clean up
		packMediatorMock.ExpectedCalls = make([]*mock.Call, 0)
	})
}
