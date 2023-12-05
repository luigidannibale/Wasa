package database

import (
	"fmt"

	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (InternalServerError, AlreadyDone)
*/
func (db *appdbimpl) CreateFollow(follow utils.Follow) (string, error) {
	userID, userToFollowID := follow.FollowerID, follow.FollowedID
	res, err := db.c.Exec(`INSERT OR IGNORE INTO Follows(FollowerID,FollowedID) VALUES (?,?)`, userID, userToFollowID)

	if err != nil {
		return err.Error(), InternalServerError
	}
	if x, y := res.RowsAffected(); x == 0 && y == nil {
		return "Already following", AlreadyDone
	}
	return fmt.Sprintf("User %d started following %d", userID, userToFollowID), nil

}
