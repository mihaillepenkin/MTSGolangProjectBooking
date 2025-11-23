package api

import (
	"database/sql"
	"encoding/json"
	"hotel_service/internal/application/dto/request"
	"hotel_service/internal/application/usecase"
	"log/slog"
	"net/http"
	"strconv"
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
		slog.Error(err.Error())
		switch res.Error {
		case "418":
			w.WriteHeader(http.StatusTeapot)
		case "502":
			w.WriteHeader(http.StatusBadGateway)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func (hh *HotelHandler) AddHotelInfo(w http.ResponseWriter, r *http.Request) {
	var body request.HotelInfoAdditionRequestDto
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Неверный формат входных данных", http.StatusBadRequest)
		return
	}

	res, err := hh.hotelService.AddHotelInfo(&body)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		slog.Error(err.Error())
		switch res.Error {
		case "418":
			w.WriteHeader(http.StatusTeapot)
		case "502":
			w.WriteHeader(http.StatusBadGateway)
		}
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func (hh *HotelHandler) UpdateHotelInfo(w http.ResponseWriter, r *http.Request) {
	var body request.HotelInfoUpdateRequestDto
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Неверный формат входных данных", http.StatusBadRequest)
		return
	}

	res, err := hh.hotelService.UpdateHotelInfo(&body)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		slog.Error(err.Error())
		switch res.Error {
		case "418":
			w.WriteHeader(http.StatusTeapot)
		case "502":
			w.WriteHeader(http.StatusBadGateway)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func (hh *HotelHandler) GetHotelById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	hotelId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "Некорректный ID отеля", http.StatusBadRequest)
		return
	}

	res, err := hh.hotelService.GetHotelById(hotelId)
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		slog.Error(err.Error())
		switch res.Error {
		case "418":
			w.WriteHeader(http.StatusTeapot)
		case "502":
			w.WriteHeader(http.StatusBadGateway)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}
