package service

import (
	"L0-wbtech/internal/cache"
	"L0-wbtech/internal/model"
	"L0-wbtech/internal/storage"
	"L0-wbtech/pkg/errors"
	"L0-wbtech/pkg/logger/sl"
	"context"
	stdErrors "errors"
	"fmt"
	"log/slog"
)

type orderService struct {
	storage storage.Storage
	cache   cache.Cache
	log     *slog.Logger
}

func NewOrderService(
	storage storage.Storage,
	cache cache.Cache,
	log *slog.Logger,
) Service {
	return &orderService{
		storage: storage,
		cache:   cache,
		log:     log,
	}
}

func (s *orderService) CreateOrder(ctx context.Context, order *model.Order) error {
	const op = "service.orderService.CreateOrder"
	log := s.log.With(
		slog.String("op", op),
		slog.String("order_uid", order.OrderUID),
	)

	if order.OrderUID == "" {
		log.Error("Order UID is empty")
		return fmt.Errorf("%s: %w", op, errors.ErrInvalidInput)
	}

	if err := s.storage.CreateOrder(ctx, order); err != nil {
		log.Error("Failed to create order",
			sl.Err(err),
			"order_uid", order.OrderUID)
		return fmt.Errorf("%s: %w", op, err)
	}

	s.cache.Set(order)
	log.Info("Order created and cached")
	return nil
}

func (s *orderService) GetOrder(ctx context.Context, orderUID string) (*model.Order, error) {
	const op = "service.orderService.GetOrder"
	log := s.log.With(
		slog.String("op", op),
		slog.String("order_uid", orderUID),
	)

	if order, ok := s.cache.Get(orderUID); ok {
		log.Info("Order retrieved from cache")
		return order, nil
	}

	order, err := s.storage.GetOrder(ctx, orderUID)
	if err != nil {
		if stdErrors.Is(err, errors.ErrNotFound) {
			log.Warn("Order not found in storage")
			return nil, fmt.Errorf("%s: %w", op, errors.ErrNotFound)
		}

		log.Error("Failed to get order from storage", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s.cache.Set(order)
	log.Info("Order retrieved from storage and cached")
	return order, nil
}

func (s *orderService) RestoreCache(ctx context.Context) error {
	const op = "service.orderService.RestoreCache"
	log := s.log.With(slog.String("op", op))

	orders, err := s.storage.GetAllOrders(ctx)
	if err != nil {
		log.Error("Failed to get all orders", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, order := range orders {
		s.cache.Set(order)
	}

	log.Info("Cache restored", "orders_count", len(orders))
	return nil
}

func (s *orderService) Close() error {
	const op = "service.orderService.Close"
	log := s.log.With(slog.String("op", op))

	if err := s.storage.Close(); err != nil {
		log.Error("Failed to close storage", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Service closed")
	return nil
}
