package main

import (
	"context"
	"fmt"
	"notification_service/internal/application/usecase"
	"notification_service/internal/config"
	messagebroker "notification_service/internal/infrastructure/message_broker"
	"notification_service/internal/infrastructure/repository"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	conf, err := config.LoadConfig()
	if (err != nil) {
		fmt.Println(err)
		os.Exit(1)
	}
	logger, _ := zap.NewProduction()
	rep := repository.NewRep(conf.SMTP.FromEmail, conf.SMTP.Password, conf.SMTP.SMTPHost, conf.SMTP.SMTPPort)
	serv := usecase.NewService(rep)
	kafka := messagebroker.NewBroker(conf.Kafka.Brokers, conf.Kafka.Topic, conf.Kafka.GroupID, logger, serv)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		logger.Info("Shutdown signal received", zap.String("signal", sig.String()))
		cancel()
	}()
	if err := kafka.StartConsuming(ctx); err != nil && err != context.Canceled {
		logger.Fatal("Consumer error", zap.Error(err))
	}

	logger.Info("Service stopped gracefully.")
}