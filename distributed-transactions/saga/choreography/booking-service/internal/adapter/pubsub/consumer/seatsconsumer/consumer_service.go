package seatsconsumer

import (
	"encoding/json"
	"log/slog"

	"github.com/IBM/sarama"
)

func (*SeatsConsumerHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (*SeatsConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (consumer *SeatsConsumerHandler) ConsumeClaim(
	sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	slog.Info("listening to seat reservation service events")
	for event := range claim.Messages() {
		var seatReservationEvent map[string]any
		if err := json.Unmarshal(event.Value, &seatReservationEvent); err != nil {
			slog.Error("error binding seat reservation req")
			continue
		}
		slog.Info("seat reservation event", slog.Any("body", seatReservationEvent))
		//TODO
		/*
		* Consume the message in "seat-reservation-status" topic
		* Iff the status is false which indicates that the seat reservation failed for some reason, start the compensating transaction and put the booking status to "failed"
		 */

		// Temporary in-mem logic
		bookingId, ok := seatReservationEvent["booking_id"].(float64)
		if !ok {
			slog.Info("failed to get booking id")
			continue
		}
		bookingIdInt := int32(bookingId)
		if seatReservationEvent["event_type"] == "seat-reservation-failed" {
			consumer.app.Bookings[bookingIdInt].SeatReservationStatus = "failed"
			consumer.app.Bookings[bookingIdInt].BookingStatus = "failed"
		} else {
			consumer.app.Bookings[bookingIdInt].SeatReservationStatus = "pending"
		}
	}
	return nil
}
