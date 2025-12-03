package app

import (
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/application/usecase/bookingprovider"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/application/usecase/bookingsaver"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/application/usecase/jwtservice"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/application/usecase/transactionmanager"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user"
)

type Services struct {
	BookingSaver    booking.Saver
	BookingProvider booking.Provider
	TokenService    user.TokenService
}

func NewServices(cfg *Config, repos *Repositories, txManager transactionmanager.TransactionManager[string]) *Services {
	bookingSaver := bookingsaver.NewBookingSaver(repos.BookingRepository, txManager, repos.HotelRepo, repos.PaymentSender, cfg.HTTPConfig.Host+":"+cfg.HTTPConfig.Port+cfg.HTTPConfig.WebhookHandlerEndpoint)
	bookingProvider := bookingprovider.NewBookingProvider(repos.BookingRepository, repos.HotelRepo)
	tokenService := jwtservice.NewJwtService(cfg.JWTConfig.SecretKey)
	return &Services{bookingSaver, bookingProvider, tokenService}
}
