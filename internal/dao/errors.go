package dao

import (
	"fmt"
)

type ErrNoProductFound struct {
	productId int
}

func NewErrNoProductFound(productId int) *ErrNoProductFound {
	return &ErrNoProductFound{productId: productId}
}

func (e *ErrNoProductFound) Error() string {
	return fmt.Sprintf("product with id %d not found", e.productId)
}
