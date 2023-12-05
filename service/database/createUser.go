package database

import (
	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (InternalServerError)
*/
func (db *appdbimpl) CreateUserByUsername(username string) (int, string, error) {
	var userID int
	u, s, e := db.GetUserByUsername(username)
	if e != nil {
		switch e {
		case NotFound:
			err := db.c.QueryRow(`INSERT INTO Users(Username) VALUES (?) RETURNING Id`, username).Scan(&userID)
			if err == nil {
				return userID, "Created", nil
			}
			return userID, err.Error(), InternalServerError
		case InternalServerError:
			return userID, s, InternalServerError
		}
	}
	return u.Id, "Logged", nil
}

func (db *appdbimpl) CreateUser(user utils.User) (int, string, error) {
	var userID int
	userID, s, e := db.CreateUserByUsername(user.Username)
	if e != nil {
		return userID, s, e
	}
	user.Id = userID
	_, s, e = db.UpdateUser(user)
	return userID, s, e
}
