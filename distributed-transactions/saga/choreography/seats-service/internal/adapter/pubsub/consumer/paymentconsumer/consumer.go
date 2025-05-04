package paymentconsumer

import "github.com/adityapatel-00/system-design/distributed-transactions/saga/choreography/seats-service/internal/app"

type PaymentConsumerHandler struct {
	app *app.Application
}

func NewPaymentConsumerHandler(app *app.Application) *PaymentConsumerHandler {
	return &PaymentConsumerHandler{
		app,
	}
}
