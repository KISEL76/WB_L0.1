package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"wb_test/internal/model"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	errUnknownType = errors.New("unknown event type")
	flushTimeout   = 5000
)

type Producer struct {
	producer *ckafka.Producer
}

func NewProducer(address []string) (*Producer, error) {
	cfg := &ckafka.ConfigMap{
		"bootstrap.servers": strings.Join(address, ","),
	}

	p, err := ckafka.NewProducer(cfg)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Something wrong with new producer: %v", err)
	}
	return &Producer{producer: p}, nil
}

func (p *Producer) Produce(ctx context.Context, order *model.Order, topic string) error {
	payload, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("[ERROR] Something wrong with data marshalling: %v", err)
	}

	kafkaMsg := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{
			Topic:     &topic,
			Partition: ckafka.PartitionAny,
		},
		Value: payload,
		Key:   nil,
	}

	kafkaChan := make(chan ckafka.Event)
	if err := p.producer.Produce(kafkaMsg, kafkaChan); err != nil {
		return fmt.Errorf("[ERROR] Something wrong with producing of message: %v", err)
	}

	e := <-kafkaChan
	switch event := e.(type) {
	case *ckafka.Message:
		return nil
	case ckafka.Error:
		return event
	default:
		return errUnknownType
	}
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}
