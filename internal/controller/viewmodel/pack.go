package viewmodel

type PackRequest struct {
	Size int `json:"size" validate:"required"`
}
