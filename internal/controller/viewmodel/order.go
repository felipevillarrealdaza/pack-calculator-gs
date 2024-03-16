package viewmodel

type OrderRequest struct {
	OrderQuantity int `json:"quantity" validate:"required"`
}

type OrderPack struct {
	Size     int `json:"size"`
	Quantity int `json:"quantity"`
}

type OrderResponse struct {
	Packs []OrderPack `json:"packs"`
}
