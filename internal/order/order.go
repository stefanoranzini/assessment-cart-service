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

	var order = model.Order{
		OrderPrice: decimal.Zero,
		OrderVat:   decimal.Zero,
		Items:      make([]*model.Item, len(orderRequest.Items)),
	}

	for i, item := range orderRequest.Items {
		s.processItem(ctx, &order, &item, i)
	}

	orderID, err := s.orderDao.Insert(ctx, order.OrderPrice, order.OrderVat)
	if err != nil {
		return nil, err
	}
	order.OrderId = orderID

	return &order, nil
}

func (s *OrderService) processItem(ctx context.Context, order *model.Order, item *model.ItemRequest, itemPosition int) error {
	itemPrice, err := s.productDao.FetchProductPrice(ctx, item.ProductID)
	if err != nil {
		return err
	}

	itemVat := itemPrice.Mul(defaultVatPercentage)
	order.Items[itemPosition] = &model.Item{
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
		Price:     itemPrice,
		Vat:       itemVat,
	}

	quantity := decimal.NewFromInt(int64(item.Quantity))
	order.OrderPrice = order.OrderPrice.Add(itemPrice.Mul(quantity))
	order.OrderVat = order.OrderVat.Add(itemVat.Mul(quantity))
	return nil
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
