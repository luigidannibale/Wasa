package database

import (
	"errors"
	"fmt"

	"github.com/luigidannibale/Wasa/service/utils"
)

func (db *appdbimpl) DeleteLike(like utils.Like) (string, error) {
	user, s, e := db.GetUser(like.UserID)
	if e != nil {
		if e.Error() == "NotFound" {
			return "Couldn't find the user", errors.New("UserNotFound")
		}
		if e.Error() == "InternalServerError" {
			return s, errors.New("InternalServerError")
		}
	}
	_, s, e = db.GetPhoto(like.PhotoID)
	if e != nil {
		if e.Error() == "NotFound" {
			return "Couldn't find the photo", errors.New("PhotoNotFound")
		}
		if e.Error() == "InternalServerError" {
			return s, errors.New("InternalServerError")
		}
	}
	res, err := db.c.Exec(`DELETE FROM Likes 
				WHERE UserID = ? AND PhotoID = ?`, like.UserID, like.PhotoID)
	if x, y := res.RowsAffected(); x == 0 && y == nil {
		return fmt.Sprintf("There is no like of %s to the photo %d", user.Username, like.PhotoID), errors.New("BannedNotFound")
	}
	if err != nil {
		return "An error occurred on the server" + err.Error(), errors.New("InternalServerError")
	}
	return fmt.Sprintf("User %s has unliked %d", user.Username, like.PhotoID), nil
}
