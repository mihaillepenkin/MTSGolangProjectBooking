package response

type HotelInfoResponseDto struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	OwnerId  int64  `json:"ownerId"`
	Rooms    []Room `json:"rooms"`
	Message  string `json:"message"`
	Error    string `json:"error"`
}
