package models

import "time"

type Event struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	DateTime    time.Time `json:"date_time" binding:"required"`
	UserID      uint      `json:"user_id"`
}

var events = make([]*Event, 0, 128)

func (e *Event) Save() {
	events = append(events, e)
}

func GetEvents() []*Event {
	return events
}
