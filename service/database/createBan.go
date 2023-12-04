package database

import (
	"errors"
	"fmt"
)

func (db *appdbimpl) CreateBan(userID int, userToBanID int) (string, error) {

	user, s, e := db.GetUser(userID)
	if e != nil {
		if e.Error() == "NotFound" {
			return "Couldn't find the user", errors.New("UserNotFound")
		}
		if e.Error() == "InternalServerError" {
			return s, errors.New("InternalServerError")
		}
	}
	userToBan, s, e := db.GetUser(userToBanID)
	if e != nil {
		if e.Error() == "NotFound" {
			return "Couldn't find the user to ban", errors.New("BannedNotFound")
		}
		if e.Error() == "InternalServerError" {
			return s, errors.New("InternalServerError")
		}
	}

	res, err := db.c.Exec(`INSERT OR IGNORE INTO Bans(BannerID,BannedID)
									VALUES (?,?)`, user.Id, userToBan.Id)
	if x, y := res.RowsAffected(); x == 0 && y == nil {
		return fmt.Sprintf("User %s already banned %s", user.Username, userToBan.Username), errors.New("AlreadyFollowed")
	}
	if err != nil {
		return err.Error(), errors.New("InternalServerError")
	}
	return fmt.Sprintf("User %s banned %s", user.Username, userToBan.Username), nil

}
