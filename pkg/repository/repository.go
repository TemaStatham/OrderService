package repository

import (
	"github.com/TemaStatham/OrderService/pkg/cache"
	"github.com/jmoiron/sqlx"
)

type Orders interface {
}

type Repository struct {
	Orders
}

func NewRepository(db *sqlx.DB, c *cache.Cache) *Repository {
	return &Repository{
		Orders: NewOrdersPostgres(db, c),
	}
}
