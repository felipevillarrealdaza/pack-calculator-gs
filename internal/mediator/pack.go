package mediator

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/felipevillarrealdaza/go-service-template/internal/repository"
)

type PackMediatorDeps func(mediator *packMediator)

func WithPackRepository(repository repository.Querier) PackMediatorDeps {
	return func(mediator *packMediator) {
		mediator.packRepository = repository
	}
}

type PackMediator interface {
	AddPack(ctx context.Context, size int) error
	RemovePack(ctx context.Context, size int) error
}

type packMediator struct {
	packRepository repository.Querier
}

func NewPackMediator(deps ...PackMediatorDeps) PackMediator {
	packMediator := packMediator{}
	for _, opt := range deps {
		opt(&packMediator)
	}
	return packMediator
}

func (pm packMediator) AddPack(ctx context.Context, size int) error {
	// Add pack in db
	if addErr := pm.packRepository.AddPack(ctx, int32(size)); addErr != nil {
		return errors.Wrap(addErr, fmt.Sprintf("could not add pack of size [%v]", size))
	}
	return nil
}

func (pm packMediator) RemovePack(ctx context.Context, size int) error {
	// Remove pack from db
	if removeErr := pm.packRepository.RemovePackBySize(ctx, int32(size)); removeErr != nil {
		return errors.Wrap(removeErr, fmt.Sprintf("could not remove pack of size [%v]", size))
	}
	return nil
}
