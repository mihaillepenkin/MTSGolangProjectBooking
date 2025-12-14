package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/message"
	kafka2 "github.com/segmentio/kafka-go"
	"github.com/testcontainers/testcontainers-go/modules/kafka"
	"gotest.tools/v3/assert"
)

var (
	testKafka    *kafka.KafkaContainer
	testProducer *KafkaProducer
	testConsumer *kafka2.Reader
	testTopic    = "testTopic"
	testMessage  = message.Message{
		Name:          "1",
		Email:         "1",
		OperationType: message.BookingOperationType,
		OperationInfo: "1",
		Status:        message.StatusOK,
		Msg:           "",
	}
)

func setup() error {
	start := time.Now()
	ctx := context.Background()

	kafkaContainer, err := kafka.Run(ctx,
		"confluentinc/confluent-local:7.5.0",
		kafka.WithClusterID("test-cluster"))

	if err != nil {
		return fmt.Errorf("setup kafka: %w", err)
	}

	testKafka = kafkaContainer

	brokers, err := kafkaContainer.Brokers(ctx)

	if err != nil {
		return fmt.Errorf("setup kafka: %w", err)
	}

	if err = createTopic(ctx, brokers[0]); err != nil {
		return fmt.Errorf("setup kafka: %w", err)
	}

	writer := &kafka2.Writer{
		Addr:     kafka2.TCP(brokers[0]),
		Topic:    testTopic,
		Balancer: &kafka2.RoundRobin{},
	}

	testProducer = NewKafkaProducer(writer)

	testConsumer = kafka2.NewReader(kafka2.ReaderConfig{
		Brokers:  brokers,
		Topic:    testTopic,
		GroupID:  "test-group",
		MinBytes: 10e3,
		MaxBytes: 10e6,
		MaxWait:  time.Second,
	})

	time.Sleep(2 * time.Second)

	duration := time.Now().Sub(start)
	log.Printf("setup completed in %v", duration)

	return nil
}

func createTopic(ctx context.Context, broker string) error {
	admin := kafka2.Client{Addr: kafka2.TCP(broker)}

	_, err := admin.CreateTopics(ctx, &kafka2.CreateTopicsRequest{
		Topics: []kafka2.TopicConfig{
			{
				Topic:             testTopic,
				NumPartitions:     1,
				ReplicationFactor: 1,
			},
		},
	})

	if err != nil {
		if err.Error() != "topic already exists" {
			return fmt.Errorf("create topic: %w", err)
		}
	}

	return nil
}

func terminate() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := testKafka.Terminate(ctx); err != nil {
		log.Printf("close kafka consumer: %v", err)
	}
}

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	terminate()

	os.Exit(code)
}

func TestKafkaProducer_SendMessage(t *testing.T) {
	ctx := context.Background()

	msg := testMessage

	err := testProducer.SendMessage(ctx, &msg)

	assert.Assert(t, err == nil)

	kafkaMessage, err := testConsumer.ReadMessage(ctx)
	assert.Assert(t, err == nil)

	var newMessage message.Message
	if err = json.Unmarshal(kafkaMessage.Value, &newMessage); err != nil {
		t.Fatalf("unmarshal json: %v", err)
	}

	assert.Assert(t, newMessage.Name == msg.Name)
}
