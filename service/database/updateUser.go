package database

import (
	"database/sql"
	"errors"

	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (NotFound, InternalServerError, UsernameTaken)
*/
func (db *appdbimpl) UpdateUser(user utils.User) (utils.User, string, error) {
	oldUser, s, e := db.GetUser(user.Id)

	if e != nil {
		switch e {
		case ErrNotFound:
			return user, "Could not find the user with such id", ErrNotFound
		case ErrInternalServerError:
			return user, s, ErrInternalServerError
		}
	}

	userNull, _, e1 := db.GetUserByUsername(user.Username)

	if e1 == nil && userNull.Id != user.Id {
		return user, "This username is already used by someone else", ErrUsernameTaken
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
		return user, "User updated successfully", nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return user, err.Error(), ErrInternalServerError
	}
	return user, err.Error(), ErrInternalServerError

}
