// Package models The models package contains various business data types
package models

import (
	"database/sql"
	"event_planner_api/event_db"
	"fmt"
	"log"
	"time"
)

// The Event struct contains data describing an event.
//
// An event has an ID that uniquely identifies it,
//
// # A Name for the event
//
// # A Description describing the purpose of the event
//
// # A Location for where the event will take place
//
// # A DateTime for the date and time of the event
//
// And a UserID identifying the user that created the event
type Event struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	DateTime    time.Time `json:"date_time" binding:"required"`
	UserID      int64     `json:"user_id"`
}

// Save Saves the Event to the database
func (e *Event) Save() error {
	result, err := event_db.InsertStmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	e.ID = id
	return err
}

func GetEvent(id int64) (Event, error) {
	row := event_db.GetEventWithIdStmt.QueryRow(id)
	var event Event
	err := scanEvent(&event, row)
	return event, err
}

// GetEvents gets all the events saved in the database
func GetEvents() ([]Event, error) {
	query := fmt.Sprintf("SELECT * FROM %s", event_db.EVENTS_TABLE_NAME)
	rows, err := event_db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Error closing rows: %s\n", err)
		}
	}(rows)

	var events = make([]Event, 0, 32)
	for rows.Next() {
		var event Event
		err = rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
		if err != nil {
			log.Println("Error scanning row into even struct")
			continue
		}
		events = append(events, event)
	}
	return events, nil
}

func (e *Event) UpdateEvent() error {
	_, err := event_db.UpdateEventStmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.ID)
	return err
}

func (e *Event) Delete() error {
	_, err := event_db.DeleteEventStmt.Exec(e.ID)
	return err
}

func (e *Event) Register(userId int64) error {
	query := fmt.Sprintf("INSERT INTO %s(event_id, user_id) VALUES (?, ?)",
		event_db.REGISTRATIONS_TABLE_NAME)
	stmt, err := event_db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Printf("error closing insert registration statement: %s\n", err)
		}
	}(stmt)
	_, err = stmt.Exec(e.ID, userId)
	return err
}

func (e *Event) DeRegister(userId int64) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE event_id = ? AND user_id = ?", event_db.REGISTRATIONS_TABLE_NAME)
	stmt, err := event_db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Printf("error closing delete registration statement: %s\n", err)
		}
	}(stmt)
	_, err = stmt.Exec(e.ID, userId)
	return err
}

func scanEvent(event *Event, row *sql.Row) error {
	return row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
}
