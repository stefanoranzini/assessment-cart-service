package order

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stefanoranzini/assessment/cart-service/internal/dao"
	"github.com/stefanoranzini/assessment/cart-service/internal/model"
	"github.com/stefanoranzini/assessment/cart-service/test/helper"
	"github.com/stretchr/testify/assert"
)

const notExistingProductId = 999

func TestOrderValidation(t *testing.T) {

	t.Run("given an order request with no items, it should return an error", func(t *testing.T) {
		orderRequest := &model.OrderRequest{}

		err := validateOrder(orderRequest)

		assert.ErrorIs(t, err, ErrInvalidOrderRequest)
	})

	t.Run("given an order request having items with quantity less than 1, it should return an error", func(t *testing.T) {
		orderRequest := &model.OrderRequest{
			Items: []*model.ItemRequest{
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
func TestOrderInsertion(t *testing.T) {

	testDb := helper.PerpareTemporaryTestDb(t)
	productDao := dao.NewProductDao(testDb)
	orderDao := dao.NewOrderDao(testDb)

	sut := NewOrderService(productDao, orderDao)

	t.Run("given a valid order request then should return a new order computing total price and vat", func(t *testing.T) {
		helper.TruncateProductTable(t, testDb)
		helper.InsertProduct(t, testDb, 1, decimal.NewFromInt(1))
		helper.InsertProduct(t, testDb, 2, decimal.NewFromFloat(10.2))
		orderRequest := &model.OrderRequest{
			Items: []*model.ItemRequest{
				{
					ProductId: 1,
					Quantity:  2,
				},
				{
					ProductId: 2,
					Quantity:  1,
				},
			},
		}

		actualOrder, err := sut.Insert(t.Context(), orderRequest)

		assert.NoError(t, err)

		assert.Equal(t, &model.Order{
			OrderId:    1,
			OrderPrice: decimal.NewFromFloat(12.2),
			OrderVat:   decimal.NewFromFloat(2.44),
			Items: []*model.Item{
				{
					ProductId: 1,
					Quantity:  2,
					Price:     decimal.NewFromInt(1),
					Vat:       decimal.NewFromFloat(0.2),
				},
				{
					ProductId: 2,
					Quantity:  1,
					Price:     decimal.NewFromFloat(10.2),
					Vat:       decimal.NewFromFloat(2.04),
				},
			}}, actualOrder)
	})

	t.Run("given an order with not existing product then should return an error", func(t *testing.T) {
		helper.TruncateProductTable(t, testDb)
		helper.InsertProduct(t, testDb, 1, decimal.NewFromInt(1))
		orderRequest := &model.OrderRequest{
			Items: []*model.ItemRequest{
				{
					ProductId: 1,
					Quantity:  2,
				},
				{
					ProductId: notExistingProductId,
					Quantity:  1,
				},
			},
		}

		_, err := sut.Insert(t.Context(), orderRequest)

		assert.ErrorIs(t, err, dao.ErrNoProductFound)
	})

	t.Run("given an invalid order request then should return an error", func(t *testing.T) {
		orderRequest := &model.OrderRequest{}

		_, err := sut.Insert(t.Context(), orderRequest)

		assert.ErrorIs(t, err, ErrInvalidOrderRequest)
	})

	t.Run("given an error inserting the order then should return an error", func(t *testing.T) {
		closedTestDb := helper.PerpareTemporaryTestDb(t)
		closedTestDb.Close()
		productDao := dao.NewProductDao(testDb)
		orderDao := dao.NewOrderDao(closedTestDb)

		sut := NewOrderService(productDao, orderDao)

		orderRequest := &model.OrderRequest{
			Items: []*model.ItemRequest{
				{
					ProductId: 1,
					Quantity:  2,
				},
			},
		}

		_, err := sut.Insert(t.Context(), orderRequest)

		assert.Error(t, err)
	})
}
