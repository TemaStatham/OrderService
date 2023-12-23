package service

import (
	"fmt"

	"github.com/TemaStatham/OrderService/pkg/model"
	"github.com/TemaStatham/OrderService/pkg/repository"
)

type OrderSerice struct {
	r repository.Orders
}

func NewOrderService(r repository.Orders) *OrderSerice {
	return &OrderSerice{
		r: r,
	}
}

func (o *OrderSerice) GetOrder(orderID string) (*model.OrderClient, error) {
	if orderID == "" {
		return nil, fmt.Errorf("orderID is nil")
	}
	return o.r.GetOrder(orderID)
}

func (o *OrderSerice) AddOrder(order *model.OrderClient) (string, error) {
	if order == nil {
		return "", fmt.Errorf("orderptr is nil")
	}

	return o.r.AddOrder(order)
}

func (o *OrderSerice) GetRecentOrders(countOrders int) ([]*model.OrderClient, error) {
	if countOrders <= 0 {
		return []*model.OrderClient{}, fmt.Errorf("count orders unvalue\n")
	}

	return o.r.GetRecentOrders(countOrders)
}
