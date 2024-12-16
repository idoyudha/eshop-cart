package kafka

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/idoyudha/eshop-cart/pkg/logger"
)

type ConsumerServer struct {
	consumer *kafka.Consumer
	l        logger.Interface
}

func NewKafkaConsumer(brokerURL, groupID string, topics []string, l logger.Interface) (*ConsumerServer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokerURL,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %v", err)
	}

	err = c.SubscribeTopics(topics, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to topics: %v", err)
	}

	return &ConsumerServer{
		consumer: c,
		l:        l,
	}, nil
}

func (c *ConsumerServer) Close() error {
	return c.consumer.Close()
}

func (s *ConsumerServer) Consume(ctx context.Context, handler func(*kafka.Message) error) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := s.consumer.ReadMessage(-1)
			if err != nil {
				s.l.Error(fmt.Errorf("error reading message: %w", err))
				continue
			}

			if err := handler(msg); err != nil {
				s.l.Error(fmt.Errorf("error handling message: %w", err))
			}
		}
	}
}
