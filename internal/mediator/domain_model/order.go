package domain_model

import (
	"math"

	"github.com/felipevillarrealdaza/go-service-template/internal/controller/viewmodel"
	"github.com/google/uuid"
)

type Order struct {
	OrderId  uuid.UUID
	Quantity int
}

func (o Order) ToViewModel() {

}

type AvailablePacks []int
type OrderPack map[int]int

func (op OrderPack) TotalItemsAndPackages() (int, int) {
	totalItemsPackaged, totalPackages := 0, 0
	for packKey, packQuantity := range op {
		totalItemsPackaged += packKey * packQuantity
		totalPackages += packQuantity
	}
	return totalItemsPackaged, totalPackages
}

type Pack struct {
	PackId   uuid.UUID
	PackSize int
}

type OrderPacks struct {
	OrderId          uuid.UUID
	OrderQuantity    int
	AvailablePacks   AvailablePacks
	ResultGrid       map[int]OrderPack
	BestItemQuantity int
	BestPackQuantity int
	OptimalOrderPack OrderPack
}

func (o *OrderPacks) ResetItemsAndPackageQuantities() {
	o.BestItemQuantity = math.MaxInt32
	o.BestPackQuantity = math.MaxInt32
}

func (o *OrderPacks) UseAsOptimalSolution(currentPackArrangement OrderPack) {
	items, packs := currentPackArrangement.TotalItemsAndPackages()
	o.BestItemQuantity = items
	o.BestPackQuantity = packs
	o.OptimalOrderPack = currentPackArrangement
}

func (o OrderPacks) ToViewModel() viewmodel.OrderResponse {
	var orderPacksResponse viewmodel.OrderResponse
	for packSize, packQuantity := range o.ResultGrid[o.OrderQuantity] {
		orderPacksResponse.Packs = append(orderPacksResponse.Packs, viewmodel.OrderPack{Size: packSize, Quantity: packQuantity})
	}
	return orderPacksResponse
}
