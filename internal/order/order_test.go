package order

import (
	"testing"

	"github.com/stefanoranzini/assessment/cart-service/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {

	t.Run("given an order request with no items, it should return an error", func(t *testing.T) {
		orderRequest := &model.OrderRequest{}

		err := validateOrder(orderRequest)

		assert.ErrorIs(t, err, ErrInvalidOrderRequest)
	})

	t.Run("given an order request having items with quantity less than 1, it should return an error", func(t *testing.T) {
		orderRequest := &model.OrderRequest{
			Items: []model.ItemRequest{
				{
					ProductId: 1,
					Quantity:  0,
				},
				{
					ProductId: 2,
					Quantity:  1,
				},
			},
		}

		err := validateOrder(orderRequest)

		assert.ErrorIs(t, err, ErrInvalidOrderRequest)
	})
}
