package database

import (
	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (NotFound, InternalServerError, UsernameTaken)
*/
func (db *appdbimpl) UpdateUser(user utils.User) (utils.User, string, error) {

	userNull, _, e1 := db.GetUserByUsername(user.Username)

	if e1 == nil && userNull.Id != user.Id {
		return user, "This username is already used by someone else", ErrUsernameTaken
	}

	res, err := db.c.Exec(`	UPDATE Users
							SET Username = ? ,Name = ? ,Surname = ?,DateOfBirth = ?
							Where Id = ?`, user.Username, user.Name, user.Surname, user.DateOfBirth.String(), user.Id)

	if err != nil {
		return user, err.Error(), ErrInternalServerError
	}
	x, e := res.RowsAffected()
	if x == 0 {
		return user, e.Error(), ErrNotFound
	}
	if e != nil {
		return user, e.Error(), ErrInternalServerError
	}
	return user, "User updated successfully", nil

}
