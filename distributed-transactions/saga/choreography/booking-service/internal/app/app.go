package app

import (
	"log"

	"github.com/adityapatel-00/system-design/distributed-transactions/saga/choreography/booking-service/internal/domain"
	"github.com/adityapatel-00/system-design/distributed-transactions/saga/utility"
)

type Application struct {
	AppName        string
	KafkaProducer  *utility.Producer
	ProducerTopics map[string]string
	ConsumerTopics map[string]string

	Bookings map[int32]*domain.BookingDetails
}

func NewApp() *Application {
	return initApp()
}

func initApp() *Application {

	producerTopics := map[string]string{"booking-status": "booking-status"}
	consumerTopics := map[string]string{"seat-reservation-status": "seat-reservation-status", "payment-status": "payment-status"}

	// Init Booking Producer
	producer, err := utility.NewProducer()
	if err != nil {
		log.Fatalf("error initializing kafka producer")
	}

	// Init App
	app := &Application{
		AppName:        "booking-service",
		KafkaProducer:  producer,
		ProducerTopics: producerTopics,
		ConsumerTopics: consumerTopics,
		Bookings:       make(map[int32]*domain.BookingDetails),
	}

	return app
}
