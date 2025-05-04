package utility

import (
	"encoding/json"
	"log/slog"

	"github.com/IBM/sarama"
)

const (
	KafkaServerAddress = "localhost:9092"
)

type Producer struct {
	Producer sarama.SyncProducer
}

type ProduceEventRequest struct {
	Topic string
	Key   string
	Body  map[string]interface{}
}

func NewProducer() (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{KafkaServerAddress}, config)
	if err != nil {
		return nil, err
	}

	return &Producer{
		Producer: producer,
	}, nil
}

func (p *Producer) ProduceNewEvent(req *ProduceEventRequest) error {

	bodyBytes, err := json.Marshal(req.Body)
	if err != nil {
		slog.Error("Error marshalling body")
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: req.Topic,
		Key:   sarama.StringEncoder(req.Key),
		Value: sarama.StringEncoder(bodyBytes),
	}

	_, _, err = p.Producer.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}
