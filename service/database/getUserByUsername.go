package database

import (
	"database/sql"
	"errors"
)

/*Returns userID of the user*/
func (db *appdbimpl) GetUserByUsername(username string) (int, string, error) {
	var userID int
	e := db.c.QueryRow(`SELECT Id 
						FROM Users
						WHERE Username = ?`, username).Scan(&userID)
	if e == nil {
		return userID, "UserID found successfully", nil
	}
	if errors.Is(e, sql.ErrNoRows) {
		return userID, "Couldn't find the user", errors.New("NotFound")
	}
	return userID, "An error occured on the server : " + e.Error(), errors.New("InternalServerError")
}
