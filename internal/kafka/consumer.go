package kafka

import (
	"fmt"
	"log"
	"strings"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

const (
	sessionTimeout = 500
	noTimeout      = -1
)

type Consumer struct {
	consumer *ckafka.Consumer
	handler  Handler
	stop     bool
}

type Handler interface {
	HandleMessage(message *ckafka.Message, offset ckafka.Offset) error
}

func NewConsumer(handler Handler, address []string, topic string, consumerGroup string) (*Consumer, error) {
	cfg := &ckafka.ConfigMap{
		"bootstrap.servers":        strings.Join(address, ","),
		"group.id":                 consumerGroup,
		"auto.offset.reset":        "earliest",
		"enable.auto.commit":       true,
		"auto.commit.interval.ms":  5000,
		"session.timeout.ms":       10000,
		"enable.auto.offset.store": false,
	}

	c, err := ckafka.NewConsumer(cfg)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Something wrong with new consumer: %v", err)
	}

	if err := c.Subscribe(topic, nil); err != nil {
		return nil, fmt.Errorf("[ERROR] Consumer can't subscribe: %v", err)
	}

	return &Consumer{consumer: c, handler: handler}, nil
}

func (c *Consumer) Start() {
	for {
		if c.stop {
			break
		}
		kafkaMsg, err := c.consumer.ReadMessage(noTimeout)
		if err != nil {
			log.Printf("%v", err)
		}
		if kafkaMsg == nil {
			continue
		}
		if err := c.handler.HandleMessage(kafkaMsg, kafkaMsg.TopicPartition.Offset); err != nil {
			log.Printf("%v", err)
			continue
		}
		if _, err := c.consumer.StoreMessage(kafkaMsg); err != nil {
			log.Printf("%v", err)
			continue
		}
	}
}

func (c *Consumer) Stop() error {
	c.stop = true
	if _, err := c.consumer.Commit(); err != nil {
		return err
	}
	return c.consumer.Close()
}
