package mediator_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/felipevillarrealdaza/go-service-template/internal/mediator"
	"github.com/felipevillarrealdaza/go-service-template/internal/mediator/domain_model"
	"github.com/felipevillarrealdaza/go-service-template/internal/repository"
	repository_mocks "github.com/felipevillarrealdaza/go-service-template/internal/repository/mocks"
	"github.com/google/uuid"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_CreateOrder_OK(t *testing.T) {
	// Set Up
	repositoryMock := repository_mocks.NewQuerier(t)
	orderMediator := mediator.NewOrderMediator(mediator.WithOrderRepository(repositoryMock))

	// Arrange
	order := domain_model.Order{
		OrderId:  uuid.New(),
		Quantity: 500,
	}
	orderRepositoryParams := repository.AddOrderParams{
		OrderID:       order.OrderId,
		OrderQuantity: int32(order.Quantity),
	}
	repositoryMock.On("AddOrder", mock.Anything, orderRepositoryParams).Return(nil)

	// Act
	creationErr := orderMediator.CreateOrder(context.Background(), order)

	// Assert
	repositoryMock.AssertExpectations(t)
	require.NoError(t, creationErr)

	// Clean up
	repositoryMock.ExpectedCalls = make([]*mock.Call, 0)
}

func Test_CreateOrder_InvalidOrderQuantity(t *testing.T) {
	// Set Up
	repositoryMock := repository_mocks.NewQuerier(t)
	orderMediator := mediator.NewOrderMediator(mediator.WithOrderRepository(repositoryMock))

	t.Run("Zero order quantity", func(t *testing.T) {
		// Arrange
		order := domain_model.Order{
			OrderId:  uuid.New(),
			Quantity: 0,
		}

		// Act
		creationErr := orderMediator.CreateOrder(context.Background(), order)

		// Assert
		repositoryMock.AssertExpectations(t)
		require.Error(t, creationErr)

		// Clean up
		repositoryMock.ExpectedCalls = make([]*mock.Call, 0)
	})

	t.Run("Negative order quantity", func(t *testing.T) {
		// Arrange
		order := domain_model.Order{
			OrderId:  uuid.New(),
			Quantity: -50,
		}

		// Act
		creationErr := orderMediator.CreateOrder(context.Background(), order)

		// Assert
		repositoryMock.AssertExpectations(t)
		require.Error(t, creationErr)

		// Clean up
		repositoryMock.ExpectedCalls = make([]*mock.Call, 0)
	})
}

func Test_CalculateOrderPacks_OK(t *testing.T) {
	// Set Up
	repositoryMock := repository_mocks.NewQuerier(t)
	orderMediator := mediator.NewOrderMediator(mediator.WithOrderRepository(repositoryMock))
	useCases := []struct {
		order          repository.Order
		availablePacks []int32
		optimalResult  domain_model.OrderPack
	}{
		{
			order:          repository.Order{OrderID: uuid.New(), OrderQuantity: 8},
			availablePacks: []int32{2, 5},
			optimalResult:  domain_model.OrderPack{2: 4, 5: 0},
		},
		{
			order:          repository.Order{OrderID: uuid.New(), OrderQuantity: 5},
			availablePacks: []int32{2, 5},
			optimalResult:  domain_model.OrderPack{2: 0, 5: 1},
		},
		{
			order:          repository.Order{OrderID: uuid.New(), OrderQuantity: 12},
			availablePacks: []int32{2, 5},
			optimalResult:  domain_model.OrderPack{2: 1, 5: 2},
		},
		{
			order:          repository.Order{OrderID: uuid.New(), OrderQuantity: 19},
			availablePacks: []int32{2, 5},
			optimalResult:  domain_model.OrderPack{2: 2, 5: 3},
		},
		{
			order:          repository.Order{OrderID: uuid.New(), OrderQuantity: 21},
			availablePacks: []int32{2, 5},
			optimalResult:  domain_model.OrderPack{2: 3, 5: 3},
		},
		{
			order:          repository.Order{OrderID: uuid.New(), OrderQuantity: 8},
			availablePacks: []int32{15, 33, 50},
			optimalResult:  domain_model.OrderPack{15: 1, 33: 0, 50: 0},
		},
		{
			order:          repository.Order{OrderID: uuid.New(), OrderQuantity: 100},
			availablePacks: []int32{15, 33, 50},
			optimalResult:  domain_model.OrderPack{15: 0, 33: 0, 50: 2},
		},
		{
			order:          repository.Order{OrderID: uuid.New(), OrderQuantity: 111},
			availablePacks: []int32{15, 33, 50},
			optimalResult:  domain_model.OrderPack{15: 3, 33: 2, 50: 0},
		},
		{
			order:          repository.Order{OrderID: uuid.New(), OrderQuantity: 82},
			availablePacks: []int32{15, 33, 50},
			optimalResult:  domain_model.OrderPack{15: 0, 33: 1, 50: 1},
		},
		{
			order:          repository.Order{OrderID: uuid.New(), OrderQuantity: 333},
			availablePacks: []int32{15, 33, 50},
			optimalResult:  domain_model.OrderPack{15: 0, 33: 1, 50: 6},
		},
	}

	for _, useCase := range useCases {
		t.Run(fmt.Sprintf("Packages: [%+v], Quantity: [%+v]", useCase.availablePacks, useCase.order.OrderQuantity), func(t *testing.T) {
			// Arrange
			repositoryMock.
				On("RetrieveOrderById", mock.Anything, useCase.order.OrderID).
				Return(repository.Order{OrderID: useCase.order.OrderID, OrderQuantity: useCase.order.OrderQuantity}, nil)
			repositoryMock.On("RetrievePacks", mock.Anything).Return(useCase.availablePacks, nil)
			repositoryMock.On("AddOrderPack", mock.Anything, mock.Anything).Return(nil)

			// Act
			orderPacks, calculationErr := orderMediator.CalculateOrderPacks(context.Background(), useCase.order.OrderID)

			// Assert
			repositoryMock.AssertExpectations(t)
			require.NoError(t, calculationErr)
			require.Equal(t, useCase.optimalResult, orderPacks.OptimalOrderPack)

			// Clean up
			repositoryMock.ExpectedCalls = make([]*mock.Call, 0)
		})
	}
}

func Test_CreateOrder_Errors(t *testing.T) {
	// Set Up
	repositoryMock := repository_mocks.NewQuerier(t)
	orderMediator := mediator.NewOrderMediator(mediator.WithOrderRepository(repositoryMock))

	// Arrange
	order := domain_model.Order{
		OrderId:  uuid.New(),
		Quantity: 500,
	}
	orderRepositoryParams := repository.AddOrderParams{
		OrderID:       order.OrderId,
		OrderQuantity: int32(order.Quantity),
	}
	repositoryMock.On("AddOrder", mock.Anything, orderRepositoryParams).Return(nil)

	// Act
	creationErr := orderMediator.CreateOrder(context.Background(), order)

	// Assert
	repositoryMock.AssertExpectations(t)
	require.NoError(t, creationErr)

	// Clean up
	repositoryMock.ExpectedCalls = make([]*mock.Call, 0)
}
