package service

import "github.com/TemaStatham/OrderService/pkg/repository"

type Order interface {
}

type Service struct {
	Order
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Order: NewOrderService(&repos.Orders),
	}
}
