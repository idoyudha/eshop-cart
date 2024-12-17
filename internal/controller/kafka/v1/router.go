package v1

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/idoyudha/eshop-cart/internal/usecase"
	"github.com/idoyudha/eshop-cart/pkg/kafka"
	"github.com/idoyudha/eshop-cart/pkg/logger"
)

func KafkaNewRouter(
	ucp usecase.Cart,
	l logger.Interface,
	c *kafka.ConsumerServer,
) {
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
		default:
			ev, err := c.Consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// log.Println("CONSUME CART SERVICE!!")
				// Errors are informational and automatically handled by the consumer
				continue
			}
			log.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))
		}
	}
}

type KafkaProductUpdatedMessage struct {
	ProductID    string  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductPrice float64 `json:"product_price"`
}

// func (r *KafkaRouter) HandleProductUpdated(msg *kafka.Message) error {
// 	var message KafkaProductUpdatedMessage

// 	if err := json.Unmarshal(msg.Value, &message); err != nil {
// 		return err
// 	}

// 	log.Println("product updated", message)
// 	// TODO: update product name and price in redis then mysql
// 	return nil
// }
