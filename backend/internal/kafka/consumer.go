package kafka

import (
	"L0-wbtech/internal/model"
	"L0-wbtech/internal/service"
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader       *kafka.Reader
	orderService service.Service
	log          *slog.Logger
}

func NewConsumer(
	brokers []string,
	topic string,
	groupID string,
	orderService service.Service,
	log *slog.Logger,
) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        brokers,
			Topic:          topic,
			GroupID:        groupID,
			MinBytes:       10e3,
			MaxBytes:       10e6,
			MaxWait:        1 * time.Second,
			StartOffset:    kafka.LastOffset,
			CommitInterval: 0,
		}),
		orderService: orderService,
		log:          log,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	log := c.log.With("operation", "kafka.Consumer.Start")
	log.Info("Starting Kafka consumer")

	for {
		select {
		case <-ctx.Done():
			log.Info("Stopping Kafka consumer")
			return
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Error("Fetch error", "error", err)
				continue
			}

			c.processMessage(ctx, msg)
		}
	}
}

func (c *Consumer) processMessage(ctx context.Context, msg kafka.Message) {
	log := c.log.With("operation", "kafka.Consumer.processMessage")

	var order model.Order
	if err := json.Unmarshal(msg.Value, &order); err != nil {
		log.Error("Unmarshal error", "error", err, "message", string(msg.Value))
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Error("Commit error", "error", err)
		}
		return
	}

	if err := order.Validate(); err != nil {
		log.Error("Invalid order data", "error", err, "order_uid", order.OrderUID)
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Error("Commit error", "error", err)
		}
		return
	}

	if order.OrderUID == "" {
		log.Error("Received order with empty UID")
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Error("Commit error", "error", err)
		}
		return
	}

	log = log.With("order_uid", order.OrderUID)
	log.Info("Processing order")

	if err := c.orderService.CreateOrder(ctx, &order); err != nil {
		log.Error("Failed to create order", "error", err)
		return
	}

	if err := c.reader.CommitMessages(ctx, msg); err != nil {
		log.Error("Commit error", "error", err)
	} else {
		log.Info("Message committed")
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
