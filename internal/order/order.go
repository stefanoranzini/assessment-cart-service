package order

import (
	"errors"
	"fmt"

	"github.com/stefanoranzini/assessment/cart-service/internal/model"
)

var ErrInvalidOrderRequest = errors.New("invalid order request")

func Insert(orderRequest *model.OrderRequest) (*model.Order, error) {
	err := validateOrder(orderRequest)
	if err != nil {
		return nil, err
	}
	return nil, errors.New("not implemented")
}

func validateOrder(orderRequest *model.OrderRequest) error {
	if len(orderRequest.Items) == 0 {
		return fmt.Errorf("no items in request: %w", ErrInvalidOrderRequest)
	}

	return nil
}
