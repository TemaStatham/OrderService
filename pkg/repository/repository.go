package repository

import "github.com/jmoiron/sqlx"

type Orders interface {
}

type Repository struct {
	Orders
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Orders: NewOrdersPostgres(db),
	}
}
