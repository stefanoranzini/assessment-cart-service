package model

type OrderRequest struct {
	Items []ItemRequest `json:"items"`
}

type ItemRequest struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
