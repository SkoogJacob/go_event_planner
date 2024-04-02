package models

import (
	"database/sql"
	"errors"
	"event_planner_api/authentication"
	"event_planner_api/event_db"
	"fmt"
	"log"
)

type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
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
	password, err := authentication.HashPassword(u.Password)
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
	query := fmt.Sprintf("SELECT id, password FROM %s WHERE email=?", event_db.USERS_TABLE_NAME)
	stmt, err := event_db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Printf("Failed to close select user password statement: %s\n", err)
		}
	}(stmt)

	row := stmt.QueryRow(u.Email)
	var toCompare string
	err = row.Scan(&u.ID, &toCompare)
	if err != nil {
		return err
	}
	matches := authentication.ComparePlainToHashed(u.Password, toCompare)
	if matches {
		return nil
	}
	return errors.New("the password did not match")
}
