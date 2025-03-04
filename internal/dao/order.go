package dao

import (
	"context"
	"database/sql"

	"github.com/shopspring/decimal"
)

type OrderDao struct {
	db *sql.DB
}

func NewOrderDao(db *sql.DB) *OrderDao {
	return &OrderDao{db: db}
}

func (o *OrderDao) Insert(ctx context.Context, price decimal.Decimal, vat decimal.Decimal) (int, error) {
	result, err := o.db.ExecContext(ctx, "INSERT INTO `order` (price, vat) VALUES ($1, $2)", price, vat)
	if err != nil {
		return 0, err
	}
	insertedId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(insertedId), nil
}
