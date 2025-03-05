package dao

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/shopspring/decimal"
)

type ProductDao struct {
	db *sql.DB
}

const fetchProductPriceQuery = `SELECT price FROM product WHERE id = $1`

func NewProductDao(db *sql.DB) *ProductDao {
	return &ProductDao{db: db}
}

func (p *ProductDao) FetchProductPrice(ctx context.Context, productId int) (decimal.Decimal, error) {
	var productPrice decimal.Decimal
	err := p.db.QueryRowContext(ctx, fetchProductPriceQuery, productId).Scan(&productPrice)
	if err != nil {
		if err == sql.ErrNoRows {
			return decimal.Zero, fmt.Errorf("id %d: %w", productId, ErrNoProductFound)
		}
		return decimal.Zero, err
	}

	return productPrice, nil
}
