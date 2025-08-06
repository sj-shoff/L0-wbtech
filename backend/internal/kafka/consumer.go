package kafka

import (
	"L0-wbtech/internal/model"
	"L0-wbtech/internal/service"
	"L0-wbtech/pkg/logger/sl"
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
	service service.Service,
	log *slog.Logger,
) *Consumer {
	log.Info("Creating Kafka consumer",
		"brokers", brokers,
		"topic", topic,
		"groupID", groupID)

	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        brokers,
			Topic:          topic,
			GroupID:        groupID,
			MinBytes:       10e3,
			MaxBytes:       10e6,
			MaxWait:        30 * time.Second,
			StartOffset:    kafka.LastOffset,
			CommitInterval: 0,
			Dialer: &kafka.Dialer{
				Timeout:   60 * time.Second,
				DualStack: true,
			},
		}),
		orderService: service,
		log:          log,
	}
}

func (c *Consumer) Start(ctx context.Context) {

	const op = "kafka.Consumer.Start"
	log := c.log.With(slog.String("op", op))

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
				log.Error("Fetch error", sl.Err(err))
				continue
			}

			c.processMessage(ctx, msg)
		}
	}
}

func (c *Consumer) processMessage(ctx context.Context, msg kafka.Message) {

	const op = "kafka.Consumer.processMessage"
	log := c.log.With(slog.String("op", op))

	var order model.Order
	if err := json.Unmarshal(msg.Value, &order); err != nil {
		log.Error("Unmarshal error", sl.Err(err), "message", string(msg.Value))
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Error("Commit error", sl.Err(err))
		}
		return
	}

	if err := order.Validate(); err != nil {
		log.Error("Invalid order data", sl.Err(err), "order_uid", order.OrderUID)
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Error("Commit error", sl.Err(err))
		}
		return
	}

	if order.OrderUID == "" {
		log.Error("Received order with empty UID")
		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			log.Error("Commit error", sl.Err(err))
		}
		return
	}

	log = log.With("order_uid", order.OrderUID)
	log.Info("Processing order")

	if err := c.orderService.CreateOrder(ctx, &order); err != nil {
		log.Error("Failed to create order", sl.Err(err))
		return
	}

	if err := c.reader.CommitMessages(ctx, msg); err != nil {
		log.Error("Commit error", sl.Err(err))
	} else {
		log.Info("Message committed")
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
