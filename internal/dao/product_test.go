package dao

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stefanoranzini/assessment/cart-service/test/helper"
	"github.com/stretchr/testify/assert"
)

const notExistingProductId = 999

func TestProductDao(t *testing.T) {

	t.Run("given a valid and filled database", func(t *testing.T) {
		testDb := helper.PerpareTemporaryTestDb(t)
		sut := NewProductDao(testDb)

		t.Run("given an existing product id, it should return the product price", func(t *testing.T) {
			helper.InsertProduct(t, testDb, 1, decimal.NewFromInt(1))

			acutalPrice, err := sut.FetchProductPrice(t.Context(), 1)

			assert.NoError(t, err)
			assert.Equal(t, decimal.NewFromInt(1), acutalPrice)
		})

		t.Run("given a non-existing product id, it should return an error", func(t *testing.T) {
			_, err := sut.FetchProductPrice(t.Context(), notExistingProductId)

			assert.ErrorIs(t, err, ErrNoProductFound)
		})
	})

	// This test simulate one of the possible scenario for generic error performing the query
	t.Run("given a not reachable/closed database should return an error", func(t *testing.T) {
		testDb := helper.PerpareTemporaryTestDb(t)
		testDb.Close()
		sut := NewProductDao(testDb)

		_, err := sut.FetchProductPrice(t.Context(), 1)

		assert.Error(t, err)
	})

}
