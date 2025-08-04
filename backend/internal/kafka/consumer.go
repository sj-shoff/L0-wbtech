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
	orderService service.OrderService
	log          *slog.Logger
}

func NewConsumer(
	brokers []string,
	topic string,
	groupID string,
	orderService service.OrderService,
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
	const op = "kafka.Consumer.Start"
	log := c.log.With(slog.String("op", op))
	log.Info("Starting Kafka consumer")

	go func() {
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
	}()
}

func (c *Consumer) processMessage(ctx context.Context, msg kafka.Message) {
	const op = "kafka.Consumer.processMessage"
	log := c.log.With(slog.String("op", op))

	var order model.Order
	if err := json.Unmarshal(msg.Value, &order); err != nil {
		log.Error("Unmarshal error",
			sl.Err(err),
			"message", string(msg.Value))
		return
	}

	log = log.With(slog.String("order_uid", order.OrderUID))
	log.Info("Processing order")

	if err := c.orderService.CreateOrder(ctx, &order); err != nil {
		log.Error("Failed to create order", sl.Err(err))
		return
	}

	if err := c.reader.CommitMessages(ctx, msg); err != nil {
		log.Error("Commit error", sl.Err(err))
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
