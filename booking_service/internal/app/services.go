package app

import (
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/usecase/case/bookingprovider"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/usecase/case/bookingsaver"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/usecase/case/eventsaver"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/usecase/case/jwtservice"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/usecase/case/transactionmanager"
)

type Services struct {
	BookingSaver    booking.Saver
	BookingProvider booking.Provider
	TokenService    user.TokenService
	EventSaver      booking.Saver
}

func NewServices(cfg *Config, repos *Repositories, txManager transactionmanager.TransactionManager[string]) *Services {
	bookingSaver := bookingsaver.NewBookingSaver(repos.BookingRepository, txManager, repos.HotelRepo, repos.PaymentSender, "http://booking-service:"+cfg.HTTPConfig.Port+cfg.HTTPConfig.WebhookHandlerEndpoint)
	bookingProvider := bookingprovider.NewBookingProvider(repos.BookingRepository, repos.HotelRepo)
	tokenService := jwtservice.NewJwtService(cfg.JWTConfig.SecretKey)
	eventSaver := eventsaver.NewEventSaver(repos.Producer, bookingSaver)
	return &Services{bookingSaver, bookingProvider, tokenService, eventSaver}
}
