package database

import (
	"database/sql"
	"errors"

	"github.com/luigidannibale/Wasa/service/utils"
)

/*Returns user searched by the id,
a string message,
and an error */
func (db *appdbimpl) GetUser(userID int) (utils.User, string, error) {
	var user utils.User
	e := db.c.QueryRow(`SELECT * 
						FROM Users
						WHERE Id = ?`, userID).Scan(&user)
	if e == nil {
		return user, "User found successfully", nil
	}
	if errors.Is(e, sql.ErrNoRows) {
		return user, "Couldn't find the user", errors.New("NotFound")
	}
	return user, "An error occured on the server : " + e.Error(), errors.New("InternalServerError")
}
