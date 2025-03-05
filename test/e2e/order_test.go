package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stefanoranzini/assessment/cart-service/internal/model"
	"github.com/stefanoranzini/assessment/cart-service/internal/server"
	"github.com/stefanoranzini/assessment/cart-service/test/helper"
	"github.com/stefanoranzini/assessment/cart-service/tools/migration/migrate"
	"github.com/stretchr/testify/assert"
)

func TestOrder(t *testing.T) {
	testPort, err := randomPort()
	if err != nil {
		t.Fatalf("failed to get random free port: %v", err)
	}

	tempFile, err := os.CreateTemp("", "testdb-*.db")
	if err != nil {
		t.Fatal(err)
	}
	db := helper.PrepareTestDb(t, tempFile.Name())
	migrate.MigrateDB(db)
	helper.InsertProduct(t, db, 1, decimal.NewFromInt(1))

	go func() {
		server.New(testPort, tempFile.Name()).Start()
	}()

	t.Run("given a valid request then should return a status created", func(t *testing.T) {
		actualResponse := callApi(t, testPort, &model.ItemRequest{
			ProductId: 1,
			Quantity:  2,
		})

		assert.Equal(t, http.StatusCreated, actualResponse.StatusCode)
		assert.Equal(t, "application/json", actualResponse.Header.Get("Content-Type"))
		defer actualResponse.Body.Close()

		var response model.Order
		err := json.NewDecoder(actualResponse.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		assert.Equal(t, 1, response.OrderId)
		assert.Equal(t, decimal.NewFromInt(2), response.OrderPrice)
		assert.Equal(t, decimal.NewFromFloat(0.4), response.OrderVat)
		assert.Len(t, response.Items, 1)
		assert.Equal(t, 1, response.Items[0].ProductId)
		assert.Equal(t, 2, response.Items[0].Quantity)
		assert.Equal(t, decimal.NewFromInt(1), response.Items[0].Price)
		assert.Equal(t, decimal.NewFromFloat(0.2), response.Items[0].Vat)
	})

	t.Run("given an invalid request then should return a bad request", func(t *testing.T) {
		actualResponse := callApi(t, testPort)

		assert.Equal(t, http.StatusBadRequest, actualResponse.StatusCode)
		assert.Equal(t, "application/json", actualResponse.Header.Get("Content-Type"))
		defer actualResponse.Body.Close()

		var response map[string]string
		err := json.NewDecoder(actualResponse.Body).Decode(&response)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		assert.Equal(t, "no items in request: invalid order request", response["error"])
	})

}

func callApi(t *testing.T, testPort int, items ...*model.ItemRequest) *http.Response {
	orderRequest := &model.OrderRequest{Items: items}
	jsonBody, err := json.Marshal(orderRequest)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}

	respose, err := http.Post(fmt.Sprintf("http://localhost:%d/order", testPort), "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	return respose
}

func randomPort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}
