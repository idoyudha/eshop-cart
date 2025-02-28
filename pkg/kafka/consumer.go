package kafka

import (
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/idoyudha/eshop-cart/config"
)

const (
	ProductGroup       = "product-group"
	ProductUpdateTopic = "product-updated"
	maxRetries         = 5
	retryDelay         = 2 * time.Second
)

type ConsumerServer struct {
	Consumer *kafka.Consumer
}

func NewKafkaConsumer(kafkaCfg config.Kafka) (*ConsumerServer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     kafkaCfg.Broker,
		"group.id":              ProductGroup,
		"auto.offset.reset":     "earliest",
		"session.timeout.ms":    6000,
		"heartbeat.interval.ms": 2000,
		"metadata.max.age.ms":   900000,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %v", err)
	}

	var subscribeErr error
	for i := 0; i < maxRetries; i++ {
		subscribeErr = c.SubscribeTopics([]string{ProductUpdateTopic}, nil)
		if subscribeErr == nil {
			break
		}
		log.Printf("attempt %d: failed to subscribe to topics: %v. retrying in %v...", i+1, subscribeErr, retryDelay)
		time.Sleep(retryDelay)
	}

	if subscribeErr != nil {
		c.Close()
		return nil, fmt.Errorf("failed to subscribe to topics after %d attempts: %v", maxRetries, subscribeErr)
	}

	return &ConsumerServer{
		Consumer: c,
	}, nil
}

func (c *ConsumerServer) Close() error {
	return c.Consumer.Close()
}
