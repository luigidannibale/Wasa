package database

import (
	"errors"
)

func (db *appdbimpl) CreateUser(username string) (int, string, error) {
	var userID int
	u, s, e := db.GetUserByUsername(username)
	if e != nil {
		if e.Error() == "NotFound" {
			err := db.c.QueryRow(`	INSERT INTO Users(Username)
									VALUES (?)
									RETURNING Id`, username).Scan(&userID)
			if err == nil {
				return userID, "Created", nil
			}
			return userID, "An error occured on the server " + err.Error(), errors.New("InternalServerError")
		}
		if e.Error() == "InternalServerError" {
			return userID, s, errors.New("InternalServerError")
		}
	}
	return u.Id, "Logged", nil
}
