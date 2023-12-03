package database

import (
	"errors"
	"fmt"
)

func (db *appdbimpl) DeleteFollow(userID int, userToUnfollowID int) (string, error) {
	user, s, e := db.GetUser(userID)
	if e != nil {
		if e.Error() == "NotFound" {
			return "Couldn't find the user", errors.New("UserNotFound")
		}
		if e.Error() == "InternalServerError" {
			return s, errors.New("InternalServerError")
		}
	}
	userToUnfollow, s, e := db.GetUser(userToUnfollowID)
	if e != nil {
		if e.Error() == "NotFound" {
			return "Couldn't find the user to unfollow", errors.New("FollowedNotFound")
		}
		if e.Error() == "InternalServerError" {
			return s, errors.New("InternalServerError")
		}
	}

	res, err := db.c.Exec(`DELETE FROM Follows WHERE FollowerID = ? AND FollowedID = ?`, user.Id, userToUnfollow.Id)
	if x, y := res.RowsAffected(); x == 0 && y == nil {
		return fmt.Sprintf("User %s is not a follower of %s", user.Username, userToUnfollow.Username), errors.New("FollowedNotFound")
	}
	if err != nil {
		return err.Error(), errors.New("InternalServerError")
	}

	return fmt.Sprintf("User %s is no longer following %s", user.Username, userToUnfollow.Username), nil

}
