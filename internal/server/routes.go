package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/stefanoranzini/assessment/cart-service/internal/model"
	"github.com/stefanoranzini/assessment/cart-service/internal/order"
)

const applicationJsonContentType = "application/json"

func insertOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var orderRequest model.OrderRequest
	err := json.NewDecoder(r.Body).Decode(&orderRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	insertedOrder, err := order.Insert(&orderRequest)
	if err != nil {
		if errors.Is(err, order.ErrInvalidOrderRequest) {
			w.WriteHeader(http.StatusBadRequest)
			writeErrorResponse(w, err)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		writeErrorResponse(w, err)
		return
	}

	jsonResponse, err := json.Marshal(insertedOrder)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", applicationJsonContentType)
	w.Write(jsonResponse)
}

func writeErrorResponse(w http.ResponseWriter, err error) {
	response := map[string]string{"error": err.Error()}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", applicationJsonContentType)
	w.Write(jsonResponse)
}
