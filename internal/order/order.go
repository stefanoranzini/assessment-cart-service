package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/stefanoranzini/assessment/cart-service/internal/dao"
	"github.com/stefanoranzini/assessment/cart-service/internal/model"
)

var ErrInvalidOrderRequest = errors.New("invalid order request")

// VAT is a country based concept, country handling is not in the domain of this assessment. Using a default value
var defaultVatPercentage = decimal.NewFromFloat(0.2)

type OrderService struct {
	productDao *dao.ProductDao
	orderDao   *dao.OrderDao
}

func NewOrderService(productDao *dao.ProductDao, orderDao *dao.OrderDao) *OrderService {
	return &OrderService{
		productDao: productDao,
		orderDao:   orderDao,
	}
}

func (s *OrderService) Insert(ctx context.Context, orderRequest *model.OrderRequest) (*model.Order, error) {
	err := validateOrder(orderRequest)
	if err != nil {
		return nil, err
	}

	var items []*model.Item
	var orderPrice = decimal.Zero
	var orderVat = decimal.Zero
	for _, item := range orderRequest.Items {
		itemPrice, err := s.productDao.FetchProductPrice(ctx, item.ProductID)
		if err != nil {
			return nil, err
		}

		itemVat := itemPrice.Mul(defaultVatPercentage)
		items = append(items, &model.Item{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     itemPrice,
			Vat:       itemVat,
		})

		orderPrice = orderPrice.Add(itemPrice.Mul(decimal.NewFromInt(int64(item.Quantity))))
		orderVat = orderVat.Add(itemVat.Mul(decimal.NewFromInt(int64(item.Quantity))))
	}

	orderID, err := s.orderDao.Insert(ctx, orderPrice, orderVat)
	if err != nil {
		return nil, err
	}

	return &model.Order{
		OrderId:    orderID,
		OrderPrice: orderPrice,
		OrderVat:   orderVat,
		Items:      items,
	}, nil
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
