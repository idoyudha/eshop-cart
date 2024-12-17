package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	ProductGroup       = "product-group"
	ProductUpdateTopic = "product-updated"
)

type ConsumerServer struct {
	Consumer *kafka.Consumer
}

func NewKafkaConsumer(brokerURL string) (*ConsumerServer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokerURL,
		"group.id":          ProductGroup,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %v", err)
	}

	err = c.SubscribeTopics([]string{
		ProductUpdateTopic,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to topics: %v", err)
	}

	return &ConsumerServer{
		Consumer: c,
	}, nil
}

func (c *ConsumerServer) Close() error {
	return c.Consumer.Close()
}
