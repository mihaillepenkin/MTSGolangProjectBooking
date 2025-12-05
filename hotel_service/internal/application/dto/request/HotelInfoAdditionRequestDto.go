package request

type Room struct {
	Number int `json:"number"`
	Price  int `json:"price"`
}

type HotelInfoAdditionRequestDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Location    string `json:"location"`
	OwnerId     int64  `json:"ownerId"`
	Rooms       []Room `json:"rooms"`
}
