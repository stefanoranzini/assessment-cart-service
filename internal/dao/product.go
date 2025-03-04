package dao

import (
	"context"
	"database/sql"

	"github.com/shopspring/decimal"
)

type ProductDao struct {
	db *sql.DB
}

const fetchProductPriceQuery = `SELECT price FROM product WHERE id = $1`

func NewProductDao(db *sql.DB) *ProductDao {
	return &ProductDao{db: db}
}

func (p *ProductDao) FetchProductPrice(ctx context.Context, productID int) (decimal.Decimal, error) {
	var productPrice decimal.Decimal
	err := p.db.QueryRowContext(ctx, fetchProductPriceQuery, productID).Scan(&productPrice)
	if err != nil {
		if err == sql.ErrNoRows {
			return decimal.Zero, NewErrNoProductFound(productID)
		}
		return decimal.Zero, err
	}

	return productPrice, nil
}
