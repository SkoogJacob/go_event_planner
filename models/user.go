package models

import (
	"database/sql"
	"event_planner_api/event_db"
	"event_planner_api/hashing"
	"fmt"
	"log"
)

type User struct {
	ID       int64
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) Save() error {
	query := fmt.Sprintf("INSERT INTO %s(email, password) VALUES (?, ?)", event_db.USERS_TABLE_NAME)
	stmt, err := event_db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Printf("Failed to close insert user statement: %s\n", err)
		}
	}(stmt)
	password, err := hashing.HashPassword(u.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Email, password)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	u.ID = id
	return err
}

func (u *User) Login() error {
	query := fmt.Sprintf("SELECT password FROM %s WHERE email=?", event_db.USERS_TABLE_NAME)
	stmt, err := event_db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Printf("Failed to close insert user statement: %s\n", err)
		}
	}(stmt)
	password, err := hashing.HashPassword(u.Password)
	if err != nil {
		return err
	}

	result, err := stmt.Exec(u.Email, password)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	u.ID = id
	return err
}
