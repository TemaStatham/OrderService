package repository

import (
	"fmt"

	"github.com/TemaStatham/OrderService/pkg/cache"
	"github.com/TemaStatham/OrderService/pkg/model"
	"github.com/jmoiron/sqlx"
)

type OrdersPostgres struct {
	db *sqlx.DB
	c  *cache.Cache
}

func NewOrdersPostgres(db *sqlx.DB, c *cache.Cache) *OrdersPostgres {
	return &OrdersPostgres{db: db, c: c}
}

func (o *OrdersPostgres) AddOrder(order *model.OrderClient) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id")

	row := o.db.QueryRow(query)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
