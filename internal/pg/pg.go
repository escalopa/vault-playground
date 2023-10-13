package pg

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderItem struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Desc     string  `json:"desc"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type Facade interface {
	GetOrderItems(ctx context.Context, orderId int) ([]OrderItem, error)
}

type facade struct {
	pg *pgxpool.Pool
}

func New(pg *pgxpool.Pool) Facade {
	return &facade{
		pg: pg,
	}
}

func (f *facade) GetOrderItems(ctx context.Context, orderId int) ([]OrderItem, error) {
	if orderId < 1 {
		return nil, errors.New("invalid order ID")
	}

	conn, err := f.pg.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	rows, err := conn.Query(ctx, `
			SELECT oi.orderItemID, p.name, p.description, p.price, oi.quantity
			FROM OrderItems oi
			JOIN Products p ON oi.productID = p.productID
			WHERE oi.orderID = $1
			ORDER BY oi.orderItemID;
	`, orderId)
	if err != nil {
		return nil, err
	}

	var items []OrderItem
	for rows.Next() {
		var item OrderItem
		err := rows.Scan(&item.ID, &item.Name, &item.Desc, &item.Price, &item.Quantity)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
