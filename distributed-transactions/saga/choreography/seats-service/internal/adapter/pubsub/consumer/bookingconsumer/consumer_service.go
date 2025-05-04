package bookingconsumer

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/IBM/sarama"
	"github.com/adityapatel-00/system-design/distributed-transactions/saga/utility"
)

func (*BookingConsumerHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (*BookingConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (consumer *BookingConsumerHandler) ConsumeClaim(
	sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	slog.Info("consuming messages from booking service", slog.String("topic", consumer.app.ConsumerTopics["booking-status"]))
	for msg := range claim.Messages() {
		var bookingReq map[string]any
		if err := json.Unmarshal(msg.Value, &bookingReq); err != nil {
			slog.Error("error binding booking req")
			continue
		}
		slog.Info("booking req", slog.Any("body", bookingReq))

		if bookingReq["event_type"] == "booking-initiated" {
			slog.Info("Booking initiated")
			slog.Info("Starting seat reservation process")

			showId, ok := bookingReq["show_id"].(float64)
			if !ok {
				slog.Info("failed to get show id")
			}

			bookingId, ok := bookingReq["booking_id"].(float64)
			if !ok {
				slog.Info("failed to get show id")
			}

			/*
			* Validate user req on the selected seats, make sure they are available.
			* If available, reserve them and push an even to "seat-reservation-status" with status true
			* Else, push an event to "seat-reservation-status" with status false so that the booking will be set to "failed"
			 */

			event := map[string]any{
				"event_type": "seat-reservation-success",
				"show_id":    showId,
				"booking_id": bookingId,
			}

			// For the time being, lets push an event that the seat reservation is failed
			err := consumer.app.KafkaProducer.ProduceNewEvent(&utility.ProduceEventRequest{
				Topic: consumer.app.ProducerTopics["seat-reservation-status"],
				Key:   fmt.Sprintf("%f", showId),
				Body:  event,
			})
			if err != nil {
				slog.Info("failed to push seat reservation event", slog.Any("err", err.Error()))
				continue
			}

			slog.Info("Sent seat reservation event", slog.Any("event_type", event["event_type"]))
		}
	}
	return nil
}
