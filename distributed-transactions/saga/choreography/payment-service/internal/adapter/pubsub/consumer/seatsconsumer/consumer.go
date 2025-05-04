package seatsconsumer

import "github.com/adityapatel-00/system-design/distributed-transactions/saga/choreography/payment-service/internal/app"

type SeatsConsumerHandler struct {
	app *app.Application
}

func NewSeatsConsumerHandler(app *app.Application) *SeatsConsumerHandler {
	return &SeatsConsumerHandler{
		app,
	}
}
