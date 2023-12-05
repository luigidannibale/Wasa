package database

import (
	"errors"

	"github.com/luigidannibale/Wasa/service/utils"
)

func (db *appdbimpl) CreateUserByUsername(username string) (int, string, error) {
	var userID int
	u, s, e := db.GetUserByUsername(username)
	if e != nil {
		if e.Error() == "NotFound" {
			err := db.c.QueryRow(`	INSERT INTO Users(Username)
									VALUES (?)
									RETURNING Id`, username).Scan(&userID)
			if err == nil {
				return userID, "Created", nil
			}
			return userID, "An error occured on the server " + err.Error(), errors.New("InternalServerError")
		}
		if e.Error() == "InternalServerError" {
			return userID, s, errors.New("InternalServerError")
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
