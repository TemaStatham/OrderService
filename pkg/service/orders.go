package service

import "github.com/TemaStatham/OrderService/pkg/repository"

type OrderSerice struct {
	r *repository.Orders
}

func NewOrderService(r *repository.Orders) *OrderSerice {
	return &OrderSerice{
		r: r,
	}
}

func (o *OrderSerice) AddNewItemInCache() {

}
