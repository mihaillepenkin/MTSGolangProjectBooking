package api

import "net/http"

func CreateRouting(hotelHandler *HotelHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/hotels", hotelHandler.GetAllHotels)

	mux.HandleFunc("POST /api/v1/hotels", hotelHandler.AddHotelInfo)

	mux.HandleFunc("PUT /api/v1/hotels", hotelHandler.UpdateHotelInfo)

	return mux
}
