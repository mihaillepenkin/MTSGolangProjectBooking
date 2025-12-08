package messagebroker

import (
	"context"
	"encoding/json"
	"fmt"
	"notification_service/internal/application/dto"
	"notification_service/internal/application/interface/service"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func NewBroker(brokers []string, topic string, group string, logger *zap.Logger, notService service.NotificationService) *Broker {
	return &Broker{
		brokers:    brokers,
		topic:      topic,
		group:      group,
		logger:     logger,
		notService: notService,
	}
}

type Broker struct {
	brokers    []string
	topic      string
	group      string
	Reader     *kafka.Reader
	logger     *zap.Logger
	notService service.NotificationService
}

func (b *Broker) StartConsuming(parentCtx context.Context) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        b.brokers,
		Topic:          b.topic,
		GroupID:        b.group,
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: 0,
	})
	b.Reader = reader
	defer reader.Close()

	b.logger.Info("Started consuming commands", zap.String("topic", b.topic))

	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		for {
			select {
			case <-parentCtx.Done():
				b.logger.Info("Shutdown signal received, stopping consumer")
				return
			default:
				msg, err := b.Reader.FetchMessage(parentCtx)
				if err != nil {
					if err == context.Canceled || err == context.DeadlineExceeded {
						return
					}
					errCh <- fmt.Errorf("fetch message: %w", err)
					return
				}

				if err := b.handleMessage(parentCtx, msg); err != nil {
					b.logger.Error("Handle message failed, skipping commit", zap.Error(err))
					continue
				}

				if err := b.Reader.CommitMessages(parentCtx, msg); err != nil {
					b.logger.Warn("Commit failed", zap.Error(err))
				}
			}
		}
	}()

	select {
	case <-parentCtx.Done():
		b.logger.Info("Graceful shutdown initiated")
		return parentCtx.Err()
	case err := <-errCh:
		b.logger.Error("Consumer failed", zap.Error(err))
		return err
	}
}

func (b *Broker) handleMessage(parentCtx context.Context, msg kafka.Message) error {
	var notification *dto.SendNotificationInputDto
	if err := json.Unmarshal(msg.Value, &notification); err != nil {
		b.logger.Error("Invalid command format",
			zap.Error(err),
			zap.ByteString("value", msg.Value),
			zap.String("topic", msg.Topic),
			zap.Int64("offset", msg.Offset),
		)
		return err
	}

	_, cancel := context.WithTimeout(parentCtx, 60*time.Second)
	defer cancel()

	//пока без вывода ошибки
	b.notService.SendNotification(notification)
	return nil
}