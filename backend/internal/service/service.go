package service

import (
	"L0-wbtech/internal/model"
	"context"
)

type Service interface {
	CreateOrder(ctx context.Context, order *model.Order) error
	GetOrder(ctx context.Context, orderUID string) (*model.Order, error)
	RestoreCache(ctx context.Context) error
	Close() error
}
