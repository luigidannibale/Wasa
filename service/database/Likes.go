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
	var exists bool
	query := `SELECT EXISTS (
					SELECT 1
					FROM Likes
					WHERE UserID = ? AND PhotoID = ?
				)`

	err := db.c.QueryRow(query, like.UserID, like.PhotoID).Scan(&exists)
	if err != nil {
		return "Error checking like existence", err
	}

	if exists {
		return "Like found successfully", nil
	} else {
		return "Couldn't find the like", ErrNotFound
	}
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
	if err != nil {
		return err.Error(), ErrInternalServerError
	}
	if x, y := res.RowsAffected(); x == 0 && y == nil {
		return "Couldn't find the like", ErrNotFound
	}

	return "Like deleted successfully", nil
}

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) GetLikersList(photoID int) ([]string, string, error) {
	var likers []string
	rows, e := db.c.Query(`SELECT Username
						FROM Likes
						JOIN Users ON UserId == Id
						WHERE PhotoID = ?`, photoID)
	if rows.Err() != nil {
		return likers, rows.Err().Error(), ErrInternalServerError
	}
	if e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return likers, "Couldn't find any like", ErrNotFound
		}
		return likers, e.Error(), ErrInternalServerError
	}
	defer rows.Close()
	for rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			return likers, err.Error(), ErrInternalServerError
		}
		likers = append(likers, username)
	}

	return likers, "List of likers found successfully", nil
}
