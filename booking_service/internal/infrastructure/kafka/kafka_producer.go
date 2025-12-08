package kafka

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/message"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
	logger *slog.Logger
}

func NewKafkaProducer(writer *kafka.Writer) *KafkaProducer {
	return &KafkaProducer{writer: writer, logger: slog.Default().With("component", "KafkaProducer")}
}

func (k *KafkaProducer) SendMessage(ctx context.Context, message *message.Message) error {
	jsonDta, err := json.Marshal(*message)
	if err != nil {
		k.logger.Error("Error marshalling message", "error", err)
		return err
	}

	err = k.writer.WriteMessages(ctx, kafka.Message{Value: jsonDta})
	if err != nil {
		k.logger.Error("Error writing message", "error", err)
		return err
	}

	return nil
}
