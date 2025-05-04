package utility

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/IBM/sarama"
)

type Consumer struct {
	Topic           string
	ConsumerHandler sarama.ConsumerGroupHandler
	Consumer        sarama.ConsumerGroup
}

func NewConsumer(groupId string, topic string, consumerHandler sarama.ConsumerGroupHandler) (*Consumer, error) {
	config := sarama.NewConfig()
	consumerGroup, err := sarama.NewConsumerGroup(
		[]string{KafkaServerAddress}, groupId, config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize consumer group: %w", err)
	}

	return &Consumer{
		Topic:           topic,
		ConsumerHandler: consumerHandler,
		Consumer:        consumerGroup,
	}, nil
}

func (c *Consumer) StartConsumer(ctx context.Context, topic string, consumerHandler sarama.ConsumerGroupHandler) {
	go func() {
		for {
			slog.Info("successfully cosuming", slog.String("topic", topic))
			err := c.Consumer.Consume(ctx, []string{topic}, consumerHandler)
			if err != nil {
				log.Printf("error from consumer: %v", err)
				return
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()
}
