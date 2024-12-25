package v1

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/idoyudha/eshop-cart/internal/entity"
	"github.com/idoyudha/eshop-cart/internal/usecase"
	kafkaConSrv "github.com/idoyudha/eshop-cart/pkg/kafka"
	"github.com/idoyudha/eshop-cart/pkg/logger"
)

type kafkaConsumerRoutes struct {
	ucp usecase.Cart
	l   logger.Interface
}

func KafkaNewRouter(
	ucp usecase.Cart,
	l logger.Interface,
	c *kafkaConSrv.ConsumerServer,
) error {
	routes := &kafkaConsumerRoutes{
		ucp: ucp,
		l:   l,
	}

	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	run := true
	for run {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: terminating\n", sig)
			run = false
			return nil
		default:
			// l.Debug("Attempting to read message...")
			ev, err := c.Consumer.ReadMessage(3 * time.Second)
			if err != nil {
				// log.Println("CONSUME CART SERVICE!!")
				// Errors are informational and automatically handled by the consumer
				if kerr, ok := err.(kafka.Error); ok && kerr.Code() == kafka.ErrTimedOut {
					// l.Debug("Timeout waiting for message, continuing...")
					continue
				}
				l.Error("Error reading message: ", err)
				continue
			}

			switch *ev.TopicPartition.Topic {
			case kafkaConSrv.ProductUpdateTopic:
				if err := routes.handleProductUpdated(ev); err != nil {
					l.Error("Failed to handle product update: %w", err)
				}
			default:
				l.Info("Unknown topic: %s", *ev.TopicPartition.Topic)
			}

			log.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
		}
	}

	return nil
}

type KafkaProductUpdatedMessage struct {
	ProductID    uuid.UUID `json:"product_id"`
	ProductName  string    `json:"product_name"`
	ProductPrice float64   `json:"product_price"`
}

func (r *kafkaConsumerRoutes) handleProductUpdated(msg *kafka.Message) error {
	var message KafkaProductUpdatedMessage

	if err := json.Unmarshal(msg.Value, &message); err != nil {
		r.l.Error(err, "http - v1 - kafkaConsumerRoutes - handleProductUpdated")
		return err
	}

	cart := &entity.Cart{
		ProductID:    message.ProductID,
		ProductName:  message.ProductName,
		ProductPrice: message.ProductPrice,
	}

	if err := r.ucp.UpdateProductNameAndPriceCart(context.Background(), cart); err != nil {
		r.l.Error(err, "http - v1 - kafkaConsumerRoutes - handleProductUpdated")
		return err
	}

	r.l.Info("Product updated", "http - v1 - kafkaConsumerRoutes - handleProductUpdated")

	return nil
}
