package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"micro-warehouse/warehouse-service/repository"

	"github.com/gofiber/fiber/v2/log"
	"github.com/streadway/amqp"
)

type RabbitMQConsumer struct {
	url  string
	conn *amqp.Connection
	ch   *amqp.Channel
	repo repository.WarehouseProductRepositoryInterface
}

type StockReductionEvent struct {
	WarehouseID uint      `json:"warehouse_id"`
	ProductID   uint      `json:"product_id"`
	Stock       int       `json:"stock"`
	MerchantID  uint      `json:"merchant_id"`
	Timestamp   time.Time `json:"timestamp"`
}

const (
	ExchangeName = "warehouse_events"
	QueueName    = "stock_reduce_queue"
	RoutingKey   = "stock_reduction"
)

func NewRabbitMQConsumer(rabbitMQURL string, repo repository.WarehouseProductRepositoryInterface) *RabbitMQConsumer {
	return &RabbitMQConsumer{
		url:  rabbitMQURL,
		repo: repo,
	}
}

func (rc *RabbitMQConsumer) connect() error {
	conn, err := amqp.Dial(rc.url)
	if err != nil {
		return fmt.Errorf("failed to connect RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return fmt.Errorf("failed to open channel: %v", err)
	}

	err = ch.ExchangeDeclare(
		ExchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	q, err := ch.QueueDeclare(
		QueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	err = ch.QueueBind(
		q.Name,
		RoutingKey,
		ExchangeName,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	rc.conn = conn
	rc.ch = ch
	return nil
}

func (rc *RabbitMQConsumer) close() {
	if rc.ch != nil {
		rc.ch.Close()
		rc.ch = nil
	}
	if rc.conn != nil {
		rc.conn.Close()
		rc.conn = nil
	}
}

func (rc *RabbitMQConsumer) StartConsuming(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Infof("[RabbitMQConsumer] Stopping consumer due to context cancellation")
			rc.close()
			return
		default:
		}

		if err := rc.connect(); err != nil {
			log.Errorf("[RabbitMQConsumer] Connection failed: %v, retrying in 5s...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Infof("[RabbitMQConsumer] Connected to RabbitMQ, start consuming")

		msgs, err := rc.ch.Consume(
			QueueName,
			"",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Errorf("[RabbitMQConsumer] Failed to start consuming: %v, reconnecting in 5s...", err)
			rc.close()
			time.Sleep(5 * time.Second)
			continue
		}

		consumed := rc.consumeLoop(ctx, msgs)
		rc.close()

		if !consumed {
			log.Warnf("[RabbitMQConsumer] Channel closed, reconnecting in 5s...")
			time.Sleep(5 * time.Second)
		}
	}
}

func (rc *RabbitMQConsumer) consumeLoop(ctx context.Context, msgs <-chan amqp.Delivery) bool {
	for {
		select {
		case <-ctx.Done():
			log.Infof("[RabbitMQConsumer] Stopping consumer due to context cancellation")
			return true
		case msg, ok := <-msgs:
			if !ok {
				return false
			}
			rc.handleMessage(ctx, msg)
		}
	}
}

func (rc *RabbitMQConsumer) handleMessage(ctx context.Context, msg amqp.Delivery) {
	if len(msg.Body) == 0 {
		msg.Nack(false, false)
		return
	}

	var event StockReductionEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		log.Errorf("[RabbitMQConsumer] handleMessage - invalid JSON: %v", err)
		msg.Nack(false, false)
		return
	}

	if err := rc.processStockReduction(ctx, event); err != nil {
		log.Errorf("[RabbitMQConsumer] handleMessage - process failed: %v", err)
		msg.Nack(false, false)
		return
	}

	msg.Ack(false)
}

func (rc *RabbitMQConsumer) processStockReduction(ctx context.Context, event StockReductionEvent) error {
	warehouseProduct, err := rc.repo.GetWarehouseProductByWarehouseIDAndProductID(ctx, event.WarehouseID, event.ProductID)
	if err != nil {
		log.Errorf("[RabbitMQConsumer] processStockReduction - 1: %v", err)
		return err
	}

	newStock := warehouseProduct.Stock - event.Stock
	if newStock < 0 {
		return errors.New("stock not enough")
	}

	warehouseProduct.Stock = newStock

	if err := rc.repo.UpdateWarehouseProduct(ctx, warehouseProduct); err != nil {
		log.Errorf("[RabbitMQConsumer] processStockReduction - 2: %v", err)
		return err
	}

	return nil
}
