package response

import boookingdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking"

type BookingsResponse struct {
	Bookings []*BookingResponse `json:"bookings"`
}

func NewBookingsResponse(bookings []*boookingdomain.Booking) *BookingsResponse {
	bookingsResponse := &BookingsResponse{Bookings: make([]*BookingResponse, 0)}
	for _, booking := range bookings {
		bookingsResponse.Bookings = append(bookingsResponse.Bookings, &BookingResponse{UserID: booking.UserID, HotelID: booking.HotelID,
			RoomNumber: booking.RoomNumber, TotalPrice: booking.TotalPrice, Currency: booking.Currency, CheckIn: booking.CheckIn, CheckOut: booking.CheckOut})
	}

	return bookingsResponse
}
