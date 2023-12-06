package database

import (
	"database/sql"
	"errors"

	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) GetFollow(follow utils.Follow) (string, error) {

	e := db.c.QueryRow(`SELECT (FollowerID, FollowedID)
						FROM Follows
						WHERE FollowerID = ? AND FollowedID = ?`, follow.FollowerID, follow.FollowedID)
	if e.Err() == nil {
		return "Follow found successfully", nil
	}
	if errors.Is(e.Err(), sql.ErrNoRows) {
		return "Couldn't find the follow", ErrNotFound
	}
	return e.Err().Error(), ErrInternalServerError
}
