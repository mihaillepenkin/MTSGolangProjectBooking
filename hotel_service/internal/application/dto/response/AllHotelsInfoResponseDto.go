package response

import "hotel_service/internal/application/dto/request"

type Hotel struct {
	Id       int64          `json:"id"`
	Name     string         `json:"name"`
	Location string         `json:"location"`
	OwnerId  int64          `json:"ownerId"`
	Rooms    []request.Room `json:"rooms"`
}

type AllHotelsInfoResponseDto struct {
	Hotels  []Hotel `json:"hotels"`
	Message string  `json:"message"`
	Error   string  `json:"error"`
}
