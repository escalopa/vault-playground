package service

import (
	"context"
	"errors"

	"github.com/escalopa/vault-playground/internal/pg"
)

type OrderItem struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Desc     string  `json:"desc"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type OrderService interface {
	GetOrderItems(ctx context.Context, orderId int) ([]OrderItem, error)
}

type ordersService struct {
	facade pg.Facade
}

func NewOrdersService(facade pg.Facade) OrderService {
	return &ordersService{
		facade: facade,
	}
}

func (s *ordersService) GetOrderItems(ctx context.Context, orderId int) ([]OrderItem, error) {
	if orderId < 1 {
		return nil, errors.New("invalid order ID")
	}

	pgOrderItems, err := s.facade.GetOrderItems(ctx, orderId)
	if err != nil {
		return nil, err
	}

	var orderItems []OrderItem
	for _, orderItem := range pgOrderItems {
		orderItems = append(orderItems, OrderItem{
			ID:       orderItem.ID,
			Name:     orderItem.Name,
			Desc:     orderItem.Desc,
			Price:    orderItem.Price,
			Quantity: orderItem.Quantity,
		})
	}

	return orderItems, nil
}
