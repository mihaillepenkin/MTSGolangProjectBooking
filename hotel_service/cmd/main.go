package main

import (
	"context"
	"hotel_service/internal/config"
	"hotel_service/internal/infrastructure/grpc"
	"hotel_service/internal/infrastructure/http"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Vlad-Ali/MTSGolangProjectBooking-protos/gen/proto/hotel"
	"google.golang.org/grpc"
)

func main() {
	textHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelError,
	})
	logger := slog.New(textHandler)
	slog.SetDefault(logger)

	db := config.ConfigureDb()
	if db == nil {
		slog.Error("Error connecting to database")
		return
	}
	err := config.RunMigrations(db)
	if err != nil {
		slog.Error("Error running migrations")
		return
	}

	hotelGrpcService := grpc2.HotelGrpcService{}
	hotelGrpcService.Initialize(db)
	grpcServer := grpc.NewServer()

	hotel.RegisterHotelServer(grpcServer, &hotelGrpcService)

	listener, err := net.Listen("tcp", ":9090")
	if err != nil {
		slog.Error("Error listening on port 9090")
		return
	}
	if err = grpcServer.Serve(listener); err != nil {
		slog.Error("Error serving gRPC server")
		return
	}

	hotelHandler := api.HotelHandler{}
	hotelHandler.Initialize(db)

	mux := api.CreateRouting(&hotelHandler)
	handler := api.AuthMiddleware(api.CORSMiddleware(mux))

	server := &http.Server{
		Addr:    ":8082",
		Handler: handler,
	}
	go func() {
		err := server.ListenAndServe()
		slog.Info("Server is running on port 8082...")
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	<-osSignals
	slog.Info("Graceful shutdown...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error(err.Error())
	}
}
