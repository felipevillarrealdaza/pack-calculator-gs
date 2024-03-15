package mediator

import (
	"context"
	"fmt"
	"math"

	"github.com/pkg/errors"

	"github.com/felipevillarrealdaza/go-service-template/internal/mediator/domain_model"
	"github.com/felipevillarrealdaza/go-service-template/internal/repository"
	"github.com/google/uuid"
)

type OrderMediatorDeps func(mediator *orderMediator)

func WithOrderRepository(repository repository.Querier) OrderMediatorDeps {
	return func(mediator *orderMediator) {
		mediator.orderRepository = repository
	}
}

type OrderMediator interface {
	CreateOrder(ctx context.Context, order domain_model.Order) error
	CalculateOrderPacks(ctx context.Context, orderId uuid.UUID) (domain_model.OrderPacks, error)
}

type orderMediator struct {
	orderRepository repository.Querier
}

func NewOrderMediator(deps ...OrderMediatorDeps) OrderMediator {
	orderMediator := orderMediator{}
	for _, opt := range deps {
		opt(&orderMediator)
	}
	return orderMediator
}

func (om orderMediator) CreateOrder(ctx context.Context, order domain_model.Order) error {
	// Create order in db
	params := repository.AddOrderParams{OrderID: order.OrderId, OrderQuantity: int32(order.Quantity)}
	if addErr := om.orderRepository.AddOrder(ctx, params); addErr != nil {
		return errors.Wrap(addErr, fmt.Sprintf("could not add order with [%v] items", params.OrderQuantity))
	}
	return nil
}

func (om orderMediator) CalculateOrderPacks(ctx context.Context, orderId uuid.UUID) (domain_model.OrderPacks, error) {
	// Retrieve order info
	order, retrieveOrderErr := om.orderRepository.RetrieveOrderById(ctx, orderId)
	if retrieveOrderErr != nil {
		return domain_model.OrderPacks{}, errors.Wrap(retrieveOrderErr, fmt.Sprintf("could not retrieve order [%v]", orderId))
	}

	// Retrieve packs info
	packs, retrievePacksErr := om.orderRepository.RetrievePacks(ctx)
	if retrievePacksErr != nil {
		return domain_model.OrderPacks{}, errors.Wrap(retrievePacksErr, "could not retrieve packs")
	}

	// Translate to domain models
	orderPacks := translateToDomainModel(order, packs)

	// Make pack calculations
	orderPacksResult := calculateOrderPacks(orderPacks)

	// Save OrderPacks in db
	if saveOrderPackersErr := om.saveEachOrderPack(ctx, orderPacksResult); saveOrderPackersErr != nil {
		return domain_model.OrderPacks{}, errors.Wrap(retrieveOrderErr, fmt.Sprintf("could not save order packs for order [%v]", orderId))
	}

	return orderPacksResult, nil
}

// Translate from repository models to domain models
func translateToDomainModel(order repository.Order, packs []int32) domain_model.OrderPacks {
	orderPacks := domain_model.OrderPacks{
		OrderId:          order.OrderID,
		OrderQuantity:    int(order.OrderQuantity),
		ResultGrid:       make(map[int]domain_model.OrderPack),
		BestItemQuantity: math.MaxInt32,
		BestPackQuantity: math.MaxInt32,
	}

	for _, pack := range packs {
		orderPacks.AvailablePacks = append(orderPacks.AvailablePacks, int(pack))
	}

	return orderPacks
}

// calculate the optimal way to package the order quantity, based on the packs configured.
func calculateOrderPacks(orderPacks domain_model.OrderPacks) domain_model.OrderPacks {
	for gridIndex := 1; gridIndex <= orderPacks.OrderQuantity; gridIndex++ {
		orderPacks.ResetItemsAndPackageQuantities()

		for _, pack := range orderPacks.AvailablePacks {
			currentPackArrangement := make(domain_model.OrderPack)

			if gridIndex <= pack {
				for _, packSize := range orderPacks.AvailablePacks {
					currentPackArrangement[packSize] = 0
				}
			} else {
				for packSize, packValue := range orderPacks.ResultGrid[gridIndex-pack] {
					currentPackArrangement[packSize] = packValue
				}
			}

			currentPackArrangement[pack]++
			totalItemsPackaged, totalPackages := currentPackArrangement.TotalItemsAndPackages()
			if totalItemsPackaged < orderPacks.BestItemQuantity {
				orderPacks.UseAsOptimalSolution(currentPackArrangement)
			} else if totalItemsPackaged == orderPacks.BestItemQuantity && totalPackages < orderPacks.BestPackQuantity {
				orderPacks.UseAsOptimalSolution(currentPackArrangement)
			}
		}

		orderPacks.ResultGrid[gridIndex] = orderPacks.OptimalOrderPack
	}

	return orderPacks
}

// Save each of the order packs in the database
func (om orderMediator) saveEachOrderPack(ctx context.Context, orderPacks domain_model.OrderPacks) error {
	for orderPackSize, orderPackQuantity := range orderPacks.OptimalOrderPack {
		addOrderPackParams := repository.AddOrderPackParams{
			OrderPacksID: uuid.New(),
			OrderID:      orderPacks.OrderId,
			PackSize:     int32(orderPackSize),
			PackQuantity: int32(orderPackQuantity),
		}

		if addOrderPackErr := om.orderRepository.AddOrderPack(ctx, addOrderPackParams); addOrderPackErr != nil {
			return errors.Wrap(addOrderPackErr, fmt.Sprintf("could not save amount of packs of size [%v] for order [%v]", addOrderPackParams.PackSize, orderPacks.OrderId))
		}
	}

	return nil
}
