package database

import (
	"database/sql"
	"errors"

	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) GetLike(like utils.Like) (string, error) {

	e := db.c.QueryRow(`SELECT (UserID, PhotoID)
						FROM Likes
						WHERE UserID = ? AND PhotoID = ?`, like.UserID, like.PhotoID)
	if e.Err() == nil {
		return "Like found successfully", nil
	}
	if errors.Is(e.Err(), sql.ErrNoRows) {
		return "Couldn't find the like", ErrNotFound
	}
	return e.Err().Error(), ErrInternalServerError
}
