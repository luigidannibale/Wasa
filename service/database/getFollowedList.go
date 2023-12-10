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
	rows, e := db.c.Query(`SELECT Id, Username, Name, Surname, DateOfBirth
						FROM Follows
						JOIN Users ON FollowedID = Id
						WHERE FollowerID = ?`, userID)
	if e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return followed, "Couldn't find any followed", ErrNotFound
		}
		return followed, e.Error(), ErrInternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		var u utils.User
		err := rows.Scan(&u.Id, &u.Username, &u.Name, &u.Surname, &u.DateOfBirth)
		if err != nil {
			return followed, e.Error(), ErrInternalServerError
		}
		followed = append(followed, u)
	}

	return followed, "List of followed found successfully", nil

}
