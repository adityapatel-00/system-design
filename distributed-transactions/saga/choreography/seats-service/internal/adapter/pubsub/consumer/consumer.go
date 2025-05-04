package consumer

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/adityapatel-00/system-design/distributed-transactions/saga/choreography/seats-service/internal/adapter/pubsub/consumer/bookingconsumer"
	"github.com/adityapatel-00/system-design/distributed-transactions/saga/choreography/seats-service/internal/adapter/pubsub/consumer/paymentconsumer"
	"github.com/adityapatel-00/system-design/distributed-transactions/saga/choreography/seats-service/internal/app"
	"github.com/adityapatel-00/system-design/distributed-transactions/saga/utility"
)

func InitPubSubConsumers(ctx context.Context, app *app.Application) {
	consumers := make([]*utility.Consumer, 0)

	// Init seats consumer
	bookingConsumerHandler := bookingconsumer.NewBookingConsumerHandler(app)
	bookingConsumer, err := utility.NewConsumer(fmt.Sprintf("%s-%s", app.AppName, app.ConsumerTopics["booking-status"]), app.ConsumerTopics["booking-status"], bookingConsumerHandler)
	if err != nil {
		log.Fatalf("error initializing kafka consumer")
	}

	// Init Payment consumer
	paymentConsumerHandler := paymentconsumer.NewPaymentConsumerHandler(app)
	paymentConsumer, err := utility.NewConsumer(fmt.Sprintf("%s-%s", app.AppName, app.ConsumerTopics["payment-status"]), app.ConsumerTopics["payment-status"], paymentConsumerHandler)
	if err != nil {
		log.Fatalf("error initializing kafka consumer")
	}

	consumers = append(consumers, bookingConsumer, paymentConsumer)

	// Start consumers
	for _, consumer := range consumers {
		slog.Info(consumer.Topic)
		consumer.StartConsumer(ctx, consumer.Topic, consumer.ConsumerHandler)
	}
}
