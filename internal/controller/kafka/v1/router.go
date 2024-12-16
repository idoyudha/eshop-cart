package v1

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/idoyudha/eshop-cart/internal/usecase"
)

type KafkaRouter struct {
	uc usecase.Cart
}

func NewKafkaRouter(uc usecase.Cart) *KafkaRouter {
	return &KafkaRouter{
		uc: uc,
	}
}

type KafkaProductUpdatedMessage struct {
	ProductID    string  `json:"product_id"`
	ProductName  string  `json:"product_name"`
	ProductPrice float64 `json:"product_price"`
}

func (r *KafkaRouter) HandleProductUpdated(msg *kafka.Message) error {
	var message KafkaProductUpdatedMessage

	if err := json.Unmarshal(msg.Value, &message); err != nil {
		return err
	}

	log.Println("product updated", message)
	// TODO: update product name and price in redis then mysql
	return nil
}
