package service

import (
	"github.com/TemaStatham/OrderService/pkg/model"
	"github.com/TemaStatham/OrderService/pkg/repository"
)

type Order interface {
	GetOrder(orderID string) (*model.OrderClient, error)
	AddOrder(order *model.OrderClient) (string, error)
	GetRecentOrders(countOrders int) ([]*model.OrderClient, error)
}

type Service struct {
	Order
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Order: NewOrderService(repos.Orders),
	}
}
