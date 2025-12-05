package entity

type Hotel struct {
	Id          int64
	Name        string
	Description string
	Location    string
	OwnerId     int64
	Rooms       []Room
}
