package dao

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stefanoranzini/assessment/cart-service/test/helper"
	"github.com/stretchr/testify/assert"
)

func TestOrderDao(t *testing.T) {

	t.Run("given a valid and filled database", func(t *testing.T) {
		testDb := helper.PerpareTemporaryTestDb(t)
		sut := NewOrderDao(testDb)

		t.Run("given a valid price and vat, it should insert the order", func(t *testing.T) {
			insertedId, err := sut.Insert(t.Context(), decimal.NewFromFloat(1.7), decimal.NewFromFloat(0.22))

			assert.NoError(t, err)
			assert.NotEqual(t, 0, insertedId)

			var acutalPrice, actualVat decimal.Decimal
			testDb.QueryRow("SELECT price, vat FROM `order` WHERE id = $1", insertedId).Scan(&acutalPrice, &actualVat)
			assert.Equal(t, decimal.NewFromFloat(1.7), acutalPrice)
			assert.Equal(t, decimal.NewFromFloat(0.22), actualVat)
		})
	})

	// This test simulate one of the possible scenario for generic error performing the query
	t.Run("given a not reachable/closed database should return an error", func(t *testing.T) {
		testDb := helper.PerpareTemporaryTestDb(t)
		testDb.Close()
		sut := NewOrderDao(testDb)

		_, err := sut.Insert(t.Context(), decimal.NewFromInt(1), decimal.NewFromInt(1))

		assert.Error(t, err)
	})

}
