package database

import (
	"errors"
	"fmt"
)

func (db *appdbimpl) CreateUser(username string) (int, string, error) {

	var userID int

	userID, s, e := db.GetUserByUsername(username)

	if e != nil {
		if e.Error() == "NotFound" {
			err := db.c.QueryRow(`	INSERT INTO Users(username)
									VALUES (?)
									RETURNING userID`, username).Scan(&userID)
			if err == nil {
				fmt.Println("Lo creo")
				return userID, "Created", nil
			}
			fmt.Println("Non lo creo")
			return userID, "An error occured on the server", errors.New("InternalServerError")
		}
		if e.Error() == "InternalServerError" {
			return userID, s, errors.New("InternalServerError")
		}
	}
	return userID, "Logged", nil
}
