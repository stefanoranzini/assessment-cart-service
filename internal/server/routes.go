package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/stefanoranzini/assessment/cart-service/internal/model"
	"github.com/stefanoranzini/assessment/cart-service/internal/order"
)

const applicationJsonContentType = "application/json"

type orderService interface {
	Insert(ctx context.Context, orderRequest *model.OrderRequest) (*model.Order, error)
}

type inserOrderHandler struct {
	orderService orderService
}

func (h *inserOrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var orderRequest model.OrderRequest
	err := json.NewDecoder(r.Body).Decode(&orderRequest)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode request")
		writeErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	insertedOrder, err := h.orderService.Insert(r.Context(), &orderRequest)
	if err != nil {
		if errors.Is(err, order.ErrInvalidOrderRequest) {
			log.Error().Err(err).Msg("Invalid order request")
			writeErrorResponse(w, http.StatusBadRequest, err)
			return
		}
		log.Error().Err(err).Msg("Failed to insert order")
		writeErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	jsonResponse, err := json.Marshal(insertedOrder)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal response")
		writeErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	log.Debug().Msg("Order created successfully")

	w.Header().Set("Content-Type", applicationJsonContentType)
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}

func writeErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	response := map[string]string{"error": err.Error()}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", applicationJsonContentType)
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}
