package app

import (
	"log"

	"github.com/adityapatel-00/system-design/distributed-transactions/saga/utility"
)

type Application struct {
	AppName        string
	KafkaProducer  *utility.Producer
	ProducerTopics map[string]string
	ConsumerTopics map[string]string
}

func NewApp() *Application {
	return initApp()
}

func initApp() *Application {

	producerTopics := map[string]string{"payment-status": "payment-status"}
	consumerTopics := map[string]string{"seat-reservation-status": "seat-reservation-status"}

	// Init Payment Producer
	producer, err := utility.NewProducer()
	if err != nil {
		log.Fatalf("error initializing kafka producer")
	}

	// Init App
	app := &Application{
		AppName:        "payment-service",
		KafkaProducer:  producer,
		ProducerTopics: producerTopics,
		ConsumerTopics: consumerTopics,
	}

	return app
}
