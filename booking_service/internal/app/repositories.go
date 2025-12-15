package app

import (
	"database/sql"
	"net/http"
	"time"

	hotel3 "github.com/Vlad-Ali/MTSGolangProjectBooking-protos/gen/proto/hotel"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/config/payment"
	bookingdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/hotel"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/message"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/payment"
	hotel2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/infrastructure/grpc"
	payment2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/infrastructure/http"
	kafka2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/infrastructure/kafka"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/infrastructure/postgres"
	"github.com/segmentio/kafka-go"
	"google.golang.org/grpc"
)

type Repositories struct {
	BookingRepository bookingdomain.Repository
	HotelRepo         hotel.Repository
	PaymentSender     payment.PaymentSender
	Producer          message.Producer
}

func NewRepositories(db *sql.DB, cfg paymentconfig.PaymentConfig, conn *grpc.ClientConn, writer *kafka.Writer) *Repositories {
	bookingRepo := postgres.NewBookingRepository(db)
	hotelClient := hotel3.NewHotelClient(conn)
	hotelRepo := hotel2.NewHotelClient(hotelClient)
	paymentSender := payment2.NewPaymentSender(&http.Client{Timeout: 30 * time.Second}, cfg.CreateEndpoint)
	producer := kafka2.NewKafkaProducer(writer)
	return &Repositories{bookingRepo, hotelRepo, paymentSender, producer}
}
