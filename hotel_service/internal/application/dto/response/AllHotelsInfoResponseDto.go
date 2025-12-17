package response

type Room struct {
	Id     int64 `json:"id"`
	Number int   `json:"number"`
	Price  int   `json:"price"`
}

type Hotel struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Location    string `json:"location"`
	OwnerId     string `json:"ownerId"`
	Rooms       []Room `json:"rooms"`
}

type AllHotelsInfoResponseDto struct {
	Hotels  []Hotel `json:"hotels"`
	Message string  `json:"message"`
	Error   string  `json:"error"`
}
