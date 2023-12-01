package database

import (
	"database/sql"
	"errors"

	"github.com/luigidannibale/Wasa/service/utils"
)

func (db *appdbimpl) UpdateUser(user utils.User) (utils.User, string, error) {
	oldUser, s, e := db.GetUser(user.Id)

	if e != nil {
		if e.Error() == "NotFound" {
			return user, "Could not find the user with such id", errors.New("NotFound")
		}
		if e.Error() == "InternalServerError" {
			return user, s, errors.New("InternalServerError")
		}
	}

	userNull, _, e1 := db.GetUserByUsername(user.Username)
	if e1 == nil && userNull.Id != user.Id {
		return user, "This username is already used by someone else", errors.New("UsernameTaken")
	}

	if user.Username == oldUser.Username {
		return user, "The username was already set so", errors.New("AlreadySo")
	}
	if &user.Name == nil {
		user.Name = oldUser.Name
	}
	if &user.Surname == nil {
		user.Surname = oldUser.Surname
	}
	if &user.DateOfBirth == nil {
		user.DateOfBirth = oldUser.DateOfBirth
	}
	_, err := db.c.Exec(`	UPDATE Users
							SET Username = ? ,Name = ? ,Surname = ?,DateOfBirth = ?
							Where Id = ?`, user.Username, user.Name, user.Surname, utils.DateToString(user.DateOfBirth), user.Id)

	if err == nil {
		return user, "Successful operation", nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return user, "An error occured on the server " + err.Error(), errors.New("InternalServerError")
	}
	return user, "An error occured on the server " + err.Error(), errors.New("InternalServerError")

}
