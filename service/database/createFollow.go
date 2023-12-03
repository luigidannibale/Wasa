package database

import (
	"errors"
	"fmt"
)

func (db *appdbimpl) CreateFollow(userID int, userToFollowID int) (string, error) {
	user, s, e := db.GetUser(userID)
	if e != nil {
		if e.Error() == "NotFound" {
			return "Couldn't find the user", errors.New("UserNotFound")
		}
		if e.Error() == "InternalServerError" {
			return s, errors.New("InternalServerError")
		}
	}

	userToFollow, s, e := db.GetUser(userToFollowID)
	if e != nil {
		if e.Error() == "NotFound" {
			return "Couldn't find the user to follow", errors.New("FollowedNotFound")
		}
		if e.Error() == "InternalServerError" {
			return s, errors.New("InternalServerError")
		}
	}

	res, err := db.c.Exec(`INSERT OR IGNORE INTO Follows(FollowerID,FollowedID)
									VALUES (?,?)`, user.Id, userToFollow.Id)
	if x, y := res.RowsAffected(); x == 0 && y == nil {
		return fmt.Sprintf("User %s was already following %s", user.Username, userToFollow.Username), errors.New("AlreadyFollowed")
	}
	if err != nil {
		return err.Error(), errors.New("InternalServerError")
	}
	return fmt.Sprintf("User %s started following %s", user.Username, userToFollow.Username), nil

}
