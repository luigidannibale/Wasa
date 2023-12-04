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
		user.DateOfBirth = utils.StringToDate(dateOfBirth.String)
	}
	if e == nil {
		return user, "User found successfully", nil
	}
	if errors.Is(e, sql.ErrNoRows) {
		return user, "Couldn't find the user", errors.New("NotFound")
	}
	return user, "An error occured on the server : " + e.Error(), errors.New("InternalServerError")
}