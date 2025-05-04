package bookingconsumer

import "github.com/adityapatel-00/system-design/distributed-transactions/saga/choreography/seats-service/internal/app"

type BookingConsumerHandler struct {
	app *app.Application
}

func NewBookingConsumerHandler(app *app.Application) *BookingConsumerHandler {
	return &BookingConsumerHandler{
		app,
	}
}
