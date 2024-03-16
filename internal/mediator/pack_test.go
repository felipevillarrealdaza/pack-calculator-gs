package mediator_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/felipevillarrealdaza/go-service-template/internal/mediator"
	repository_mocks "github.com/felipevillarrealdaza/go-service-template/internal/repository/mocks"
	"github.com/pkg/errors"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_AddPack_OK(t *testing.T) {
	// Set Up
	repositoryMock := repository_mocks.NewQuerier(t)
	packMediator := mediator.NewPackMediator(mediator.WithPackRepository(repositoryMock))

	// Arrange
	repositoryMock.On("AddPack", mock.Anything, int32(10)).Return(nil)

	// Act
	addPackErr := packMediator.AddPack(context.Background(), 10)

	// Assert
	repositoryMock.AssertExpectations(t)
	require.NoError(t, addPackErr)

	// Clean up
	repositoryMock.ExpectedCalls = make([]*mock.Call, 0)
}

func Test_AddPack_Errors(t *testing.T) {
	// Set Up
	repositoryMock := repository_mocks.NewQuerier(t)
	packMediator := mediator.NewPackMediator(mediator.WithPackRepository(repositoryMock))

	t.Run("Pack of size zero", func(t *testing.T) {
		// Act
		creationErr := packMediator.AddPack(context.Background(), 0)

		// Assert
		repositoryMock.AssertExpectations(t)
		require.Error(t, creationErr)

		// Clean up
		repositoryMock.ExpectedCalls = make([]*mock.Call, 0)
	})

	t.Run("Pack of negative size", func(t *testing.T) {
		// Act
		creationErr := packMediator.AddPack(context.Background(), -5)

		// Assert
		repositoryMock.AssertExpectations(t)
		require.Error(t, creationErr)

		// Clean up
		repositoryMock.ExpectedCalls = make([]*mock.Call, 0)
	})

	t.Run("Error saving the pack", func(t *testing.T) {
		// Arrange
		repositoryMock.On("AddPack", mock.Anything, int32(10)).Return(errors.New(fmt.Sprintf("could not add pack of size [%v]", 10)))

		// Act
		creationErr := packMediator.AddPack(context.Background(), 10)

		// Assert
		repositoryMock.AssertExpectations(t)
		require.Error(t, creationErr)

		// Clean up
		repositoryMock.ExpectedCalls = make([]*mock.Call, 0)
	})
}

func Test_RemovePack_OK(t *testing.T) {
	// Set Up
	repositoryMock := repository_mocks.NewQuerier(t)
	packMediator := mediator.NewPackMediator(mediator.WithPackRepository(repositoryMock))

	// Arrange
	repositoryMock.On("RemovePackBySize", mock.Anything, int32(10)).Return(nil)

	// Act
	removePackErr := packMediator.RemovePack(context.Background(), 10)

	// Assert
	repositoryMock.AssertExpectations(t)
	require.NoError(t, removePackErr)

	// Clean up
	repositoryMock.ExpectedCalls = make([]*mock.Call, 0)
}

func Test_RemovePack_Errors(t *testing.T) {
	// Set Up
	repositoryMock := repository_mocks.NewQuerier(t)
	packMediator := mediator.NewPackMediator(mediator.WithPackRepository(repositoryMock))

	t.Run("Error removing the pack", func(t *testing.T) {
		// Arrange
		repositoryMock.On("RemovePackBySize", mock.Anything, int32(10)).Return(errors.New(fmt.Sprintf("could not remove pack of size [%v]", 10)))

		// Act
		creationErr := packMediator.RemovePack(context.Background(), 10)

		// Assert
		repositoryMock.AssertExpectations(t)
		require.Error(t, creationErr)

		// Clean up
		repositoryMock.ExpectedCalls = make([]*mock.Call, 0)
	})
}
