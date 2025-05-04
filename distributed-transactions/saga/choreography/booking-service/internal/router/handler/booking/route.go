package booking

import (
	"net/http"

	"github.com/adityapatel-00/system-design/distributed-transactions/saga/choreography/booking-service/internal/app"
)

func RegisterRoutes(router *http.ServeMux, app *app.Application) {
	router.Handle("/booking", CreateBooking(app))
	router.Handle("/list-bookings", GetBookingDetails(app))
}
