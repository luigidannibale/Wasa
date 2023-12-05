package database

import (
	"fmt"

	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) DeleteFollow(follow utils.Follow) (string, error) {
	userID, userToUnfollowID := follow.FollowerID, follow.FollowedID

	res, err := db.c.Exec(`DELETE FROM Follows WHERE FollowerID = ? AND FollowedID = ?`, userID, userToUnfollowID)
	if x, y := res.RowsAffected(); x == 0 && y == nil {
		return fmt.Sprintf("Couldn't find the user to unfollow %d in the following of user %d", userID, userToUnfollowID), NotFound
	}
	if err != nil {
		return err.Error(), InternalServerError
	}

	return fmt.Sprintf("User %d is no longer following %d", userID, userToUnfollowID), nil

}
