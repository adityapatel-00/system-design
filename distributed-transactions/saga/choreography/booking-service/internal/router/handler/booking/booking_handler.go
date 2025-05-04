package booking

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adityapatel-00/system-design/distributed-transactions/saga/choreography/booking-service/internal/app"
	"github.com/adityapatel-00/system-design/distributed-transactions/saga/choreography/booking-service/internal/domain"
	"github.com/adityapatel-00/system-design/distributed-transactions/saga/utility"
)

var bookingIdCounter int32

func CreateBooking(app *app.Application) http.HandlerFunc {
	type BookingRequest struct {
		UserId int32    `json:"user_id"`
		ShowId int32    `json:"show_id"`
		Seats  []string `json:"seats"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &BookingRequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// TODO
		/*
		* Start a local transaction in repo layer for booking creation, commit it and get the generated booking id
		 */

		bookingId := bookingIdCounter + 1
		// Temporary logic to hold the booking details
		app.Bookings[bookingId] = &domain.BookingDetails{
			UserId:                req.UserId,
			BookingId:             bookingId,
			ShowId:                req.ShowId,
			Seats:                 req.Seats,
			BookingStatus:         "pending",
			SeatReservationStatus: "pending",
			PaymentStatus:         "pending",
		}

		// Create and event to push to booking-status topic
		event := map[string]interface{}{
			"event_type": "booking-initiated",
			"booking_id": bookingId,
			"show_id":    req.ShowId,
			"seats":      req.Seats,
		}

		err := app.KafkaProducer.ProduceNewEvent(&utility.ProduceEventRequest{
			Topic: app.ProducerTopics["booking-status"],
			Key:   fmt.Sprintf("%d", req.ShowId),
			Body:  event,
		})
		if err != nil {
			http.Error(w, "booking failed", http.StatusInternalServerError)
		}

		w.Write([]byte("booking pending"))
	}
}

func GetBookingDetails(app *app.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := make(map[string]any)
		resp["bookings"] = app.Bookings

		respBytes, err := json.Marshal(resp)
		if err != nil {
			http.Error(w, "failed to get bookings", http.StatusInternalServerError)
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(respBytes)
	}
}
