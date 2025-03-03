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

	for _, item := range orderRequest.Items {
		if item.Quantity < 1 {
			return fmt.Errorf("invalid quantity for product %d: %w", item.ProductID, ErrInvalidOrderRequest)
		}
	}

	return nil
}
