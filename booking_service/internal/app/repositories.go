package app

import (
	"database/sql"
	"net/http"
	"time"

	bookingdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/hotel"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/payment"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/infrastructure/booking"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/infrastructure/config/paymentconfig"
	hotel2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/infrastructure/hotel"
	payment2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/infrastructure/payment"
	"google.golang.org/grpc"
)

type Repositories struct {
	BookingRepository bookingdomain.Repository
	HotelRepo         hotel.Repository
	PaymentSender     payment.PaymentSender
}

func NewRepositories(db *sql.DB, cfg paymentconfig.PaymentConfig, conn *grpc.ClientConn) *Repositories {
	bookingRepo := booking.NewBookingRepository(db)
	hotelRepo := hotel2.NewHotelClient(conn)
	paymentSender := payment2.NewPaymentSender(&http.Client{Timeout: 30 * time.Second}, cfg.CreateEndpoint)
	return &Repositories{bookingRepo, hotelRepo, paymentSender}
}
