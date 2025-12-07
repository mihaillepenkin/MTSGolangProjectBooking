package api

import (
	"database/sql"
	"encoding/json"
	"hotel_service/internal/application/dto/request"
	"hotel_service/internal/application/usecase"
	"log/slog"
	"net/http"
)

type HotelHandler struct {
	hotelService usecase.HotelService
}

func (hh *HotelHandler) Initialize(db *sql.DB) {
	hh.hotelService.Initialize(db)
}

func (hh *HotelHandler) GetAllHotels(w http.ResponseWriter, r *http.Request) {
	res, err := hh.hotelService.GetAllHotels()

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		slog.Error("Ошибка в HotelService, метод GetAllHotels: " + err.Error())
		switch res.Error {
		case "500":
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		slog.Error("Ошибка в HotelHandler, метод GetAllHotels: " + err.Error())
		return
	}
}

func (hh *HotelHandler) AddHotelInfo(w http.ResponseWriter, r *http.Request) {
	var body request.HotelInfoAdditionRequestDto
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		slog.Error("Ошибка в HotelHandler, метод AddHotelInfo: " + err.Error())
		http.Error(w, "Неверный формат входных данных", http.StatusBadRequest)
		return
	}

	value := r.Context().Value("claims")
	claims := value.(*JWTClaims)
	res, err := hh.hotelService.AddHotelInfo(&body, claims.UserId)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		slog.Error("Ошибка в HotelService, метод AddHotelInfo: " + err.Error())
		switch res.Error {
		case "500":
			w.WriteHeader(http.StatusInternalServerError)
		case "409":
			w.WriteHeader(http.StatusConflict)
		}
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		slog.Error("Ошибка в HotelHandler, метод AddHotelInfo: " + err.Error())
		return
	}
}

func (hh *HotelHandler) UpdateHotelInfo(w http.ResponseWriter, r *http.Request) {
	var body request.HotelInfoUpdateRequestDto
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		slog.Error("Ошибка в HotelHandler, метод UpdateHotelInfo: " + err.Error())
		http.Error(w, "Неверный формат входных данных", http.StatusBadRequest)
		return
	}

	value := r.Context().Value("claims")
	claims := value.(*JWTClaims)
	res, err := hh.hotelService.UpdateHotelInfo(&body, claims.UserId)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		slog.Error("Ошибка в HotelService, метод UpdateHotelInfo: " + err.Error())
		switch res.Error {
		case "500":
			w.WriteHeader(http.StatusInternalServerError)
		case "403":
			w.WriteHeader(http.StatusForbidden)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		slog.Error("Ошибка в HotelHandler, метод UpdateHotelInfo: " + err.Error())
		return
	}
}
