package database

import (
	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) DeleteLike(like utils.Like) (string, error) {

	res, err := db.c.Exec(`DELETE FROM Likes WHERE UserID = ? AND PhotoID = ?`, like.UserID, like.PhotoID)
	if x, y := res.RowsAffected(); x == 0 && y == nil {
		return "Couldn't find the like", ErrNotFound
	}
	if err != nil {
		return err.Error(), ErrInternalServerError
	}
	return "Like deleted successfully", nil
}
