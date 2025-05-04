package seatsconsumer

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/IBM/sarama"
	"github.com/adityapatel-00/system-design/distributed-transactions/saga/utility"
)

func (*SeatsConsumerHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (*SeatsConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (consumer *SeatsConsumerHandler) ConsumeClaim(
	sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	slog.Info("listening to seat reservation service")
	for event := range claim.Messages() {
		//TODO
		/*
		* Consume the message in "seat-reservation-status" topic
		* Iff the status is true which indicates that the seat reservation success, start the payment processing.

		* If the payment process fails for some reason, push an event to "payment-status" with status false
		 */

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
		if seatReservationEvent["event_type"] == "seat-reservation-success" {
			bookingId, ok := seatReservationEvent["booking_id"].(float64)
			if !ok {
				slog.Info("failed to get booking id")
				continue
			}
			bookingIdInt := int32(bookingId)

			showId, ok := seatReservationEvent["show_id"].(float64)
			if !ok {
				slog.Info("failed to get show id")
				continue
			}
			showIdInt := int32(showId)

			// Produce a payment-success event
			event := map[string]any{
				"event_type": "payment-success",
				"booking_id": bookingIdInt,
				"show_id":    showIdInt,
			}

			err := consumer.app.KafkaProducer.ProduceNewEvent(&utility.ProduceEventRequest{
				Topic: consumer.app.ProducerTopics["payment-status"],
				Key:   fmt.Sprintf("%d", showIdInt),
				Body:  event,
			})
			if err != nil {
				slog.Info("error producing payment-status event")
				continue
			}

			slog.Info("successfully produced payment-status event")
		}
	}
	return nil
}
