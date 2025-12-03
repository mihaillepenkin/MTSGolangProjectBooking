package booking

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/adapter/booking/request"
	response2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/adapter/booking/response"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/adapter/userkey"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking"
	error2 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/error"
	"github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/booking/object"
	error3 "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/hotel/error"
	userdomain "github.com/mihaillepenkin/MTSGolangProjectBooking/booking_service/internal/domain/user"
)

type BookingHandler struct {
	bookingSaver    booking.Saver
	bookingProvider booking.Provider
	logger          *slog.Logger
}

func NewBookingHandler(saver booking.Saver, provider booking.Provider) *BookingHandler {
	return &BookingHandler{bookingSaver: saver, bookingProvider: provider, logger: slog.Default().With("component", "booking_handler")}
}

func (b *BookingHandler) BookRoom(w http.ResponseWriter, r *http.Request) {
	user, err := userkey.ExtractUserFromReq(r)
	if err != nil {
		b.logger.Error("Error extracting user from request ", "Error", err)
		http.Error(w, "Failed to book room", http.StatusUnauthorized)
		return
	}

	if !userdomain.IsClient(user) {
		b.logger.Error("User role is not allowed", "User", user)
		http.Error(w, "User role is not allowed", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		b.logger.Error("Error reading body ", "Error", err)
		http.Error(w, "Failed to book room", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()
	var bookRoomRequest request.BookRoomRequest
	err = json.Unmarshal(body, &bookRoomRequest)
	if err != nil {
		b.logger.Error("Error unmarshalling body ", "Error", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	bookingInfo := &object.BookingInfo{User: *user, HotelName: bookRoomRequest.HotelName, RoomNumber: bookRoomRequest.RoomNumber, CheckIn: bookRoomRequest.CheckIn, CheckOut: bookRoomRequest.CheckOut}
	url, err := b.bookingSaver.BookRoom(r.Context(), bookingInfo)
	if err != nil {
		b.logger.Error("Error booking room ", "Error", err)
		if errors.Is(err, error2.ErrBookingIsIntersected) {
			http.Error(w, "Booking is intersected", http.StatusConflict)
			return
		} else if errors.Is(err, error3.ErrHotelRoomIsNotFound) {
			http.Error(w, "Hotel room is not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to book room", http.StatusInternalServerError)
		}
		return
	}

	response := response2.BookRoomResponse{URL: url}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		b.logger.Error("Error encoding response ", "Error", err)
		return
	}
}

func (b *BookingHandler) GetHotelBookings(w http.ResponseWriter, r *http.Request) {
	user, err := userkey.ExtractUserFromReq(r)
	if err != nil {
		b.logger.Error("Error extracting user from request ", "Error", err)
		http.Error(w, "Failed to get hotel bookings", http.StatusUnauthorized)
		return
	}

	if !userdomain.IsHotelier(user) {
		b.logger.Error("User role is not allowed", "User", user)
		http.Error(w, "Failed to get hotel bookings", http.StatusUnauthorized)
		return
	}

	hotelName := GetHotelNameFromRequest(r)
	if hotelName == "" {
		b.logger.Error("Invalid hotel booking name")
		http.Error(w, "Invalid hotel parameter", http.StatusBadRequest)
		return
	}

	bookings, err := b.bookingProvider.GetBookingsByHotelier(r.Context(), user, hotelName)
	if err != nil {
		b.logger.Error("Error getting bookings by hotelier ", "Error", err)
		if errors.Is(err, error2.ErrHotelierIsNotValid) {
			http.Error(w, "Hotelier is not valid", http.StatusUnauthorized)
		} else {
			http.Error(w, "Failed to get bookings", http.StatusInternalServerError)
		}
		return
	}

	response := response2.NewBookingsResponse(bookings)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		b.logger.Error("Error encoding response ", "Error", err)
		return
	}
}

func (b *BookingHandler) GetUserBookings(w http.ResponseWriter, r *http.Request) {
	user, err := userkey.ExtractUserFromReq(r)
	if err != nil {
		b.logger.Error("Error extracting user from request ", "Error", err)
		http.Error(w, "Failed to get bookings", http.StatusUnauthorized)
		return
	}

	if !userdomain.IsClient(user) {
		b.logger.Error("User role is not allowed", "User", user)
		http.Error(w, "User role is not allowed", http.StatusUnauthorized)
		return
	}

	bookings, err := b.bookingProvider.GetBookingsByUser(r.Context(), user)
	if err != nil {
		b.logger.Error("Error getting bookings by user ", "Error", err)
		http.Error(w, "Failed to get bookings", http.StatusInternalServerError)
		return
	}

	response := response2.NewBookingsResponse(bookings)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		b.logger.Error("Error encoding response ", "Error", err)
		return
	}
}

func (b *BookingHandler) GetOccupiedRoomDurations(w http.ResponseWriter, r *http.Request) {
	hotelName := GetHotelNameFromRequest(r)
	if hotelName == "" {
		b.logger.Error("Invalid hotel booking name")
		http.Error(w, "Invalid hotel parameter", http.StatusBadRequest)
		return
	}

	roomNumber := GetRoomNumberFromRequest(r)
	if roomNumber == "" {
		b.logger.Error("Invalid room number ")
		http.Error(w, "Invalid room number", http.StatusBadRequest)
		return
	}

	durations, err := b.bookingProvider.GetOccupiedRoomDurations(r.Context(), hotelName, roomNumber)
	if err != nil {
		b.logger.Error("Error getting occupied room ", "Error", err)
		http.Error(w, "Failed to get occupied room durations", http.StatusInternalServerError)
		return
	}

	response := response2.RoomDurationsResponse{Durations: durations}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		b.logger.Error("Error encoding response ", "Error", err)
		return
	}
}

func GetHotelNameFromRequest(r *http.Request) string {
	return r.URL.Query().Get("hotelName")
}

func GetRoomNumberFromRequest(r *http.Request) string {
	return r.URL.Query().Get("roomNumber")
}
