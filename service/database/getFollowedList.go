package database

import (
	"database/sql"
	"errors"

	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) GetFollowedList(userID int) ([]utils.User, string, error) {
	var followed []utils.User
	rows, e := db.c.Query(`SELECT FollowedID
						FROM Follows
						WHERE FollowerID = ?`, userID)
	if e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return followed, "Couldn't find any followed", ErrNotFound
		}
		return followed, e.Error(), ErrInternalServerError
	}
	e = rows.Scan(&followed)
	if e != nil {
		return followed, e.Error(), ErrInternalServerError
	}
	return followed, "List of followed found successfully", nil
}
