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
	var uid, pid int
	e := db.c.QueryRow(`SELECT UserID, PhotoID
						FROM Likes
						WHERE UserID = ? AND PhotoID = ?`, like.UserID, like.PhotoID).Scan(&uid, &pid)

	if e == nil {
		return "Like found successfully", nil
	}
	if errors.Is(e, sql.ErrNoRows) {
		return "Couldn't find the like", ErrNotFound
	}
	return e.Error(), ErrInternalServerError
}

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
