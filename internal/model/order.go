package model

import "github.com/shopspring/decimal"

type Order struct {
	OrderId    int             `json:"order_id"`
	OrderPrice decimal.Decimal `json:"order_price"`
	OrderVat   decimal.Decimal `json:"order_vat"`
	Items      []*Item         `json:"items"`
}

type Item struct {
	ProductID int             `json:"product_id"`
	Quantity  int             `json:"quantity"`
	Price     decimal.Decimal `json:"price"`
	Vat       decimal.Decimal `json:"vat"`
}
