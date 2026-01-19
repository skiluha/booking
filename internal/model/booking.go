package model

import "time"

type Booking struct {
	ID        int       `json:"id"`
	RoomID    int       `json:"room_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	CreatedAt time.Time `json:"created_at"`
}
