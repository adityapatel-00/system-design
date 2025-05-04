package consumer

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/adityapatel-00/system-design/distributed-transactions/saga/choreography/payment-service/internal/adapter/pubsub/consumer/seatsconsumer"
	"github.com/adityapatel-00/system-design/distributed-transactions/saga/choreography/payment-service/internal/app"
	"github.com/adityapatel-00/system-design/distributed-transactions/saga/utility"
)

func InitPubSubConsumers(ctx context.Context, app *app.Application) {
	consumers := make([]*utility.Consumer, 0)
	// Init seats consumer
	seatsConsumerHandler := seatsconsumer.NewSeatsConsumerHandler(app)
	seatsConsumer, err := utility.NewConsumer(fmt.Sprintf("%s-%s", app.AppName, app.ConsumerTopics["seat-reservation-status"]), app.ConsumerTopics["seat-reservation-status"], seatsConsumerHandler)
	if err != nil {
		log.Fatalf("error initializing kafka consumer")
	}

	consumers = append(consumers, seatsConsumer)

	// Start consumers
	for _, consumer := range consumers {
		slog.Info(consumer.Topic)
		consumer.StartConsumer(ctx, consumer.Topic, consumer.ConsumerHandler)
	}
}
