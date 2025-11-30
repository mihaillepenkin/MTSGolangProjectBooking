package entity

type Hotel struct {
	Id       int64
	Name     string
	Location string
	OwnerId  int64
	Rooms    []Room
}
