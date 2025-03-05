package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stefanoranzini/assessment/cart-service/internal/model"
	"github.com/stefanoranzini/assessment/cart-service/internal/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) Insert(ctx context.Context, orderRequest *model.OrderRequest) (*model.Order, error) {
	args := m.Called(ctx, orderRequest)

	order := args.Get(0)
	if order == nil {
		return nil, args.Error(1)
	}
	return order.(*model.Order), args.Error(1)
}

func TestInserOrderHandler(t *testing.T) {
	mockOrderService := new(MockOrderService)
	sut := &inserOrderHandler{
		orderService: mockOrderService,
	}

	t.Run("given a request with unsupported method then should return a method not allowed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/order", nil)
		responseRecorder := httptest.NewRecorder()

		sut.ServeHTTP(responseRecorder, req)

		assert.Equal(t, http.StatusMethodNotAllowed, responseRecorder.Code)
		mockOrderService.AssertNotCalled(t, "Insert", mock.Anything, mock.Anything)
	})

	t.Run("given an invalid json request then should return a bad request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewBufferString("invalid json"))
		responseRecorder := httptest.NewRecorder()

		sut.ServeHTTP(responseRecorder, req)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		assert.Equal(t, applicationJsonContentType, responseRecorder.Header().Get("Content-Type"))
		assert.JSONEq(t, `{"error":"invalid character 'i' looking for beginning of value"}`, responseRecorder.Body.String())
		mockOrderService.AssertNotCalled(t, "Insert", mock.Anything, mock.Anything)
	})

	t.Run("given an invalid order request inserting and order then should return a bad request", func(t *testing.T) {
		mockOrderService.On("Insert", mock.Anything, mock.Anything).Return(nil, order.ErrInvalidOrderRequest).Once()
		req := prepareHttpRequest(t)
		responseRecorder := httptest.NewRecorder()

		sut.ServeHTTP(responseRecorder, req)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		assert.Equal(t, applicationJsonContentType, responseRecorder.Header().Get("Content-Type"))
		assert.JSONEq(t, `{"error":"invalid order request"}`, responseRecorder.Body.String())
		mockOrderService.AssertCalled(t, "Insert", mock.Anything, mock.MatchedBy(func(orderRequest *model.OrderRequest) bool {
			return len(orderRequest.Items) == 1 && orderRequest.Items[0].ProductId == 1 && orderRequest.Items[0].Quantity == 1
		}))
	})

	t.Run("given an error inserting an order then should return an internal server error", func(t *testing.T) {
		mockOrderService.On("Insert", mock.Anything, mock.Anything).Return(nil, errors.New("generic error")).Once()
		req := prepareHttpRequest(t)
		responseRecorder := httptest.NewRecorder()

		sut.ServeHTTP(responseRecorder, req)

		assert.Equal(t, http.StatusInternalServerError, responseRecorder.Code)
		assert.Equal(t, applicationJsonContentType, responseRecorder.Header().Get("Content-Type"))
		assert.JSONEq(t, `{"error":"generic error"}`, responseRecorder.Body.String())
		mockOrderService.AssertCalled(t, "Insert", mock.Anything, mock.MatchedBy(func(orderRequest *model.OrderRequest) bool {
			return len(orderRequest.Items) == 1 && orderRequest.Items[0].ProductId == 1 && orderRequest.Items[0].Quantity == 1
		}))
	})

	t.Run("given a valid order request then should return a created response", func(t *testing.T) {
		mockOrderService.On("Insert", mock.Anything, mock.Anything).Return(&model.Order{
			OrderId:    1,
			OrderPrice: decimal.NewFromInt(77),
			OrderVat:   decimal.NewFromInt(7),
			Items: []*model.Item{
				{
					ProductId: 1,
					Quantity:  1,
					Price:     decimal.NewFromInt(70),
					Vat:       decimal.NewFromInt(7),
				},
			},
		}, nil).Once()
		req := prepareHttpRequest(t)
		responseRecorder := httptest.NewRecorder()

		sut.ServeHTTP(responseRecorder, req)

		assert.Equal(t, http.StatusCreated, responseRecorder.Code)
		assert.Equal(t, applicationJsonContentType, responseRecorder.Header().Get("Content-Type"))
		assert.JSONEq(t, `{"order_id":1,"order_price":"77","order_vat":"7","items":[{"product_id":1,"quantity":1,"price":"70","vat":"7"}]}`, responseRecorder.Body.String())
		mockOrderService.AssertCalled(t, "Insert", mock.Anything, mock.MatchedBy(func(orderRequest *model.OrderRequest) bool {
			return len(orderRequest.Items) == 1 && orderRequest.Items[0].ProductId == 1 && orderRequest.Items[0].Quantity == 1
		}))
	})
}

func prepareHttpRequest(t *testing.T) *http.Request {
	orderRequest := &model.OrderRequest{
		Items: []*model.ItemRequest{
			{
				ProductId: 1,
				Quantity:  1,
			},
		},
	}
	jsonBody, err := json.Marshal(orderRequest)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	return httptest.NewRequest(http.MethodPost, "/order", bytes.NewBuffer(jsonBody))
}
