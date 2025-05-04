package paymentconsumer

import (
	"encoding/json"
	"log/slog"

	"github.com/IBM/sarama"
)

func (*PaymentConsumerHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (*PaymentConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (consumer *PaymentConsumerHandler) ConsumeClaim(
	sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	slog.Info("listening to payment events")
	for event := range claim.Messages() {
		slog.Info("payment event", slog.Any("event", event))
		//TODO
		/*
		* Consume the message in "payment-status" topic
		* If the status is false which indicates that the payent processing failed for some reason, start the compensating transaction and put the booking status to "failed"
		* Else, update the booking status to "success"
		 */
		var paymentEvent map[string]any
		if err := json.Unmarshal(event.Value, &paymentEvent); err != nil {
			slog.Error("error binding seat reservation req")
			continue
		}

		bookingId, ok := paymentEvent["booking_id"].(float64)
		if !ok {
			slog.Info("failed to get booking id")
			continue
		}
		bookingIdInt := int32(bookingId)
		if paymentEvent["event_type"] == "payment-failed" {
			consumer.app.Bookings[bookingIdInt].SeatReservationStatus = "failed"
			consumer.app.Bookings[bookingIdInt].BookingStatus = "failed"
			consumer.app.Bookings[bookingIdInt].PaymentStatus = "failed"
		} else {
			consumer.app.Bookings[bookingIdInt].BookingStatus = "success"
			consumer.app.Bookings[bookingIdInt].SeatReservationStatus = "success"
			consumer.app.Bookings[bookingIdInt].PaymentStatus = "success"
		}

	}
	return nil
}
