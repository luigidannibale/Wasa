package database

import (
	"database/sql"
	"errors"
	"fmt"
)

func (db *appdbimpl) GetUserByUsername(username string) (int, string, error) {

	var userID int

	e := db.c.QueryRow(`SELECT userID 
						FROM Users
						WHERE username = ?`, username).Scan(&userID)
	if e == nil {
		fmt.Println("Lo trovo")
		return userID, "UserID found successfully", nil
	}
	if errors.Is(e, sql.ErrNoRows) {
		fmt.Println("Non lo trovo")
		return userID, "Couldn't find the user", errors.New("NotFound")
	}
	return userID, "An error occured on the server : " + e.Error(), errors.New("InternalServerError")

}
