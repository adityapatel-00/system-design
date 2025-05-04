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
	slog.Info("listening to payment-status")
	for event := range claim.Messages() {
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
			// TODO
			slog.Info("payment failed for", slog.Any("bookingId", bookingIdInt))
			slog.Info("rolling back seat reservation")
			// Compute the compensating transaction for seat reservation
		}
	}
	return nil
}
