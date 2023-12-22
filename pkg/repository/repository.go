package repository

import (
	"github.com/TemaStatham/OrderService/pkg/cache"
	"github.com/TemaStatham/OrderService/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Orders interface {
	AddOrder(order *model.OrderClient) (int, error)
}

type Repository struct {
	Orders
}

func NewRepository(db *sqlx.DB, c *cache.Cache) *Repository {
	return &Repository{
		Orders: NewOrdersPostgres(db, c),
	}
}
