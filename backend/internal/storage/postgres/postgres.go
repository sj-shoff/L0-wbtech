package postgres

import (
	"L0-wbtech/internal/config"
	"L0-wbtech/internal/model"
	"L0-wbtech/pkg/errors"
	"context"
	"database/sql"
	stdErrors "errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func NewPostgresDB(cfg config.Postgres) (*PostgresStorage, error) {
	const op = "storage.postgres.NewPostgresDB"

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: db.Ping error: %w", op, err)
	}

	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) CreateOrder(ctx context.Context, order *model.Order) error {
	const op = "storage.postgres.CreateOrder"

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	orderQuery := `
		INSERT INTO orders (
			order_uid, track_number, entry, locale, internal_signature, 
			customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (order_uid) DO NOTHING
	`
	_, err = tx.ExecContext(ctx, orderQuery,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated,
		order.OofShard)
	if err != nil {
		return fmt.Errorf("%s: insert order failed: %w", op, err)
	}

	deliveryQuery := `
		INSERT INTO delivery (
			order_uid, name, phone, zip, city, address, region, email
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = tx.ExecContext(ctx, deliveryQuery,
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email)
	if err != nil {
		return fmt.Errorf("%s: insert delivery failed: %w", op, err)
	}

	paymentQuery := `
		INSERT INTO payment (
			order_uid, transaction, request_id, currency, provider, amount, 
			payment_dt, bank, delivery_cost, goods_total, custom_fee
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err = tx.ExecContext(ctx, paymentQuery,
		order.OrderUID,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee)
	if err != nil {
		return fmt.Errorf("%s: insert payment failed: %w", op, err)
	}

	itemQuery := `
		INSERT INTO items (
			order_uid, chrt_id, track_number, price, rid, name, 
			sale, size, total_price, nm_id, brand, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, itemQuery,
			order.OrderUID,
			item.ChrtID,
			item.TrackNumber,
			item.Price,
			item.Rid,
			item.Name,
			item.Sale,
			item.Size,
			item.TotalPrice,
			item.NmID,
			item.Brand,
			item.Status)
		if err != nil {
			return fmt.Errorf("%s: insert item failed: %w", op, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: transaction commit failed: %w", op, err)
	}

	return nil
}

func (s *PostgresStorage) GetOrder(ctx context.Context, orderUID string) (*model.Order, error) {
	const op = "storage.postgres.GetOrder"

	orderQuery := `
		SELECT
			order_uid, track_number, entry, locale,
			internal_signature, customer_id, delivery_service,
			shardkey, sm_id, date_created, oof_shard
		FROM orders
		WHERE order_uid = $1
	`
	var order model.Order
	if err := s.db.GetContext(ctx, &order, orderQuery, orderUID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	deliveryQuery := `
		SELECT name, phone, zip, city, address, region, email
		FROM delivery
		WHERE order_uid = $1
	`
	if err := s.db.GetContext(ctx, &order.Delivery, deliveryQuery, orderUID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("%s: get delivery failed: %w", op, err)
	}

	paymentQuery := `
		SELECT
			id, transaction, request_id, currency, provider,
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payment
		WHERE order_uid = $1
	`
	if err := s.db.GetContext(ctx, &order.Payment, paymentQuery, orderUID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("%s: get payment failed: %w", op, err)
	}

	itemsQuery := `
		SELECT
			chrt_id, track_number, price, rid,
			name, sale, size, total_price, nm_id, brand, status
		FROM items
		WHERE order_uid = $1
	`
	if err := s.db.SelectContext(ctx, &order.Items, itemsQuery, orderUID); err != nil {
		if err == sql.ErrNoRows {
			order.Items = []model.Item{}
		} else {
			return nil, fmt.Errorf("%s: get items failed: %w", op, err)
		}
	}

	return &order, nil
}

func (s *PostgresStorage) GetAllOrders(ctx context.Context) (map[string]*model.Order, error) {
	const op = "storage.postgres.GetAllOrders"

	uidsQuery := `SELECT order_uid FROM orders`
	var uids []string
	err := s.db.SelectContext(ctx, &uids, uidsQuery)
	if err != nil {
		return nil, fmt.Errorf("%s: get order uids failed: %w", op, err)
	}

	orders := make(map[string]*model.Order)
	for _, uid := range uids {
		order, err := s.GetOrder(ctx, uid)
		if err != nil {
			if stdErrors.Is(err, errors.ErrNotFound) {
				continue
			}
			return nil, fmt.Errorf("%s: get order %s failed: %w", op, uid, err)
		}
		orders[uid] = order
	}

	return orders, nil
}

func (s *PostgresStorage) Close() error {
	return s.db.Close()
}
