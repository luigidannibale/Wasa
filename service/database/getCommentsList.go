package database

import (
	"database/sql"
	"errors"

	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) GetCommentsList(photoID int) ([]utils.Comment, string, error) {
	var comments []utils.Comment
	rows, e := db.c.Query(`SELECT CommentID
						FROM Comments
						WHERE PhotoID = ?`, photoID)
	if e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return comments, "Couldn't find any followed", ErrNotFound
		}
		return comments, e.Error(), ErrInternalServerError
	}
	e = rows.Scan(&comments)
	if e != nil {
		return comments, e.Error(), ErrInternalServerError
	}
	return comments, "List of comments found successfully", nil
}
