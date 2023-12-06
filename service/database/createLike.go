package database

import (
	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (InternalServerError, AlreadyDone)
*/
func (db *appdbimpl) CreateLike(like utils.Like) (string, error) {

	res, err := db.c.Exec(`INSERT OR IGNORE INTO Likes(UserID,PhotoID) VALUES (?,?)`, like.UserID, like.PhotoID)

	if err != nil {
		return err.Error(), ErrInternalServerError
	}
	if x, y := res.RowsAffected(); x == 0 && y == nil {
		return "This like already exists", ErrAlreadyDone
	}
	return "Like created successfully", nil
}
