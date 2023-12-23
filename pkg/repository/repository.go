package repository

import (
	"github.com/TemaStatham/OrderService/pkg/model"
	"github.com/jmoiron/sqlx"
)

type Orders interface {
	GetOrder(orderID string) (*model.OrderClient, error)
	GetPayment(orderID string) (*model.Payment, error)
	GetItems(orderID string) ([]model.Item, error)
	GetDelivery(orderID string) (*model.Delivery, error)
	AddOrder(order *model.OrderClient) (string, error)
}

type Repository struct {
	Orders
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Orders: NewOrdersPostgres(db),
	}
}
