package database

import (
	"database/sql"
	"errors"

	"github.com/luigidannibale/Wasa/service/utils"
	"github.com/rickb777/date"
)

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) GetUser(userID int) (utils.User, string, error) {
	var user utils.User
	var name, surname, dateOfBirth sql.NullString
	e := db.c.QueryRow(`SELECT Id,Username,Name,Surname,DateOfBirth
						FROM Users
						WHERE Id = ?`, userID).Scan(&user.Id, &user.Username, &name, &surname, &dateOfBirth)

	if name.Valid {
		user.Name = name.String
	}
	if surname.Valid {
		user.Surname = surname.String
	}
	if dateOfBirth.Valid {
		d, er := date.AutoParse(dateOfBirth.String)
		if er != nil {
			return user, er.Error(), ErrInternalServerError
		}
		user.DateOfBirth = d
	}
	if e == nil {
		return user, "User found successfully", nil
	}
	if errors.Is(e, sql.ErrNoRows) {
		return user, "Couldn't find the user", ErrNotFound
	}
	return user, e.Error(), ErrInternalServerError
}
