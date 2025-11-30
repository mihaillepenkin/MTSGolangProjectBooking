package response

import "time"

type RoomDurationsResponse struct {
	Durations [][]time.Time `json:"durations"`
}
