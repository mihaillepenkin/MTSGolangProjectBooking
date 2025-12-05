package request

type RoomUpd struct {
	Id     int64 `json:"id"`
	Number int   `json:"number"`
	Price  int   `json:"price"`
}

type HotelInfoUpdateRequestDto struct {
	Id             int64     `json:"id"`
	NewName        string    `json:"newName"`
	NewDescription string    `json:"newDescription"`
	NewLocation    string    `json:"newLocation"`
	NewOwnerId     int64     `json:"newOwnerId"`
	NewRooms       []RoomUpd `json:"newRooms"`
}
