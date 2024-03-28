package event_db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strings"
)

const TABLE_NAME = "events"

var DB *sql.DB
var InsertStmt *sql.Stmt
var GetEventWithIdStmt *sql.Stmt
var UpdateEventStmt *sql.Stmt
var DeleteEventStmt *sql.Stmt

func InitDb(path string) {
	dB, dbErr := sql.Open("sqlite3", path)
	DB = dB
	if dbErr != nil {
		log.Fatalf("Unable to open DB: %s", dbErr)
	}
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	createTables()

	toPrepare := fmt.Sprintf("SELECT * FROM %s WHERE id=?", TABLE_NAME)
	stmt, err := DB.Prepare(toPrepare)
	if err != nil {
		log.Fatalf("Unable to prepare get event query: %s\n", err)
	}
	GetEventWithIdStmt = stmt

	toPrepare = fmt.Sprintf(
		"INSERT INTO %s (name, description, location, date_time, user_id) VALUES (?, ?, ?, ?, ?)",
		TABLE_NAME)
	stmt, err = DB.Prepare(toPrepare)
	if err != nil {
		log.Fatalf("Unable to prepare insert event query: %s\n", err)
	}
	InsertStmt = stmt

	toPrepare = strings.TrimSpace(fmt.Sprintf(`
		UPDATE %s
		SET name = ?, description = ?, location = ?, date_time = ?
		WHERE id = ?
		`, TABLE_NAME))
	stmt, err = DB.Prepare(toPrepare)
	if err != nil {
		log.Fatalf("Unable to prepare update event query: %s\n", err)
	}
	UpdateEventStmt = stmt

	toPrepare = strings.TrimSpace(fmt.Sprintf(`
		DELETE FROM %s WHERE id = ?
		`, TABLE_NAME))
	stmt, err = DB.Prepare(toPrepare)
	if err != nil {
		log.Fatalf("Unable to prepare delete event query: %s\n", err)
	}
	DeleteEventStmt = stmt
}

func createTables() {
	createEventsTable := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    name TEXT NOT NULL,
		    description TEXT NOT NULL,
		    location TEXT NOT NULL,
		    date_time DATETIME NOT NULL,
		    user_id INTEGER
		);`, TABLE_NAME)
	_, err := DB.Exec(createEventsTable)
	if err != nil {
		log.Fatalf("Unable to create events table in SQLite: %s\n", err)
	}
}

func CloseDb() {
	_ = InsertStmt.Close()
	_ = GetEventWithIdStmt.Close()
	_ = UpdateEventStmt.Close()
	_ = DeleteEventStmt.Close()
	err := DB.Close()
	if err != nil {
		log.Fatalf("Failed to close DB: %s", err)
	}
}
