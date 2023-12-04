package database

import (
	"errors"
	"fmt"
)

func (db *appdbimpl) DeleteBan(userID int, userToUnbanID int) (string, error) {
	user, s, e := db.GetUser(userID)
	if e != nil {
		if e.Error() == "NotFound" {
			return "Couldn't find the user", errors.New("UserNotFound")
		}
		if e.Error() == "InternalServerError" {
			return s, errors.New("InternalServerError")
		}
	}
	userToUnban, s, e := db.GetUser(userToUnbanID)
	if e != nil {
		if e.Error() == "NotFound" {
			return "Couldn't find the user to unban", errors.New("BannedNotFound")
		}
		if e.Error() == "InternalServerError" {
			return s, errors.New("InternalServerError")
		}
	}

	res, err := db.c.Exec(`DELETE FROM Bans WHERE BannerID = ? AND BannedID = ?`, user.Id, userToUnban.Id)
	if x, y := res.RowsAffected(); x == 0 && y == nil {
		return fmt.Sprintf("User %s was not banned by %s", userToUnban.Username, user.Username), errors.New("BannedNotFound")
	}
	if err != nil {
		return "An error occurred on the server" + err.Error(), errors.New("InternalServerError")
	}

	return fmt.Sprintf("User %s has unbanned %s", user.Username, userToUnban.Username), nil

}
