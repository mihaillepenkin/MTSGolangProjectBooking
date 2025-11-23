package request

type RoomUpd struct {
	Id     int64
	Number int `json:"number"`
	Price  int `json:"price"`
}

type HotelInfoUpdateRequestDto struct {
	Id          int64     `json:"id"`
	NewName     string    `json:"newName"`
	NewLocation string    `json:"newLocation"`
	NewOwnerId  int64     `json:"newOwnerId"`
	NewRooms    []RoomUpd `json:"newRooms"`
}
