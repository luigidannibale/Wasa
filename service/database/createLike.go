package database

import (
	"errors"
	"fmt"

	"github.com/luigidannibale/Wasa/service/utils"
)

func (db *appdbimpl) CreateLike(like utils.Like) (string, error) {
	user, s, e := db.GetUser(like.UserID)
	if e != nil {
		if e.Error() == "NotFound" {
			return "Couldn't find the user", errors.New("UserNotFound")
		}
		if e.Error() == "InternalServerError" {
			return s, errors.New("InternalServerError")
		}
	}
	_, s, e = db.GetPhoto(like.UserID)
	if e != nil {
		if e.Error() == "NotFound" {
			return "Couldn't find the photo", errors.New("PhotoNotFound")
		}
		if e.Error() == "InternalServerError" {
			return s, errors.New("InternalServerError")
		}
	}

	res, err := db.c.Exec(`INSERT OR IGNORE INTO Likes(UserID,PhotoID)
									VALUES (?,?)`, like.UserID, like.PhotoID)
	if x, y := res.RowsAffected(); x == 0 && y == nil {
		return fmt.Sprintf("User %s had already liked photo %d", user.Username, like.PhotoID), errors.New("AlreadyLiked")
	}
	if err != nil {
		return err.Error(), errors.New("InternalServerError")
	}
	return "Like created successfully", nil
}
