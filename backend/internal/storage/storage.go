package storage

import (
	"L0-wbtech/internal/model"
	"context"
)

type Storage interface {
	CreateOrder(ctx context.Context, order *model.Order) error
	GetOrder(ctx context.Context, orderUID string) (*model.Order, error)
	GetAllOrders(ctx context.Context) (map[string]*model.Order, error)
	Close() error
}
