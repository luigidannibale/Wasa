package database

import (
	"database/sql"
	"errors"

	"github.com/luigidannibale/Wasa/service/utils"
)

/*
Errors that can be returned: (NotFound, InternalServerError)
*/
func (db *appdbimpl) GetComment(commentID int) (utils.Comment, string, error) {
	var comment utils.Comment
	e := db.c.QueryRow(`SELECT Id,UserID,PhotoID,Content
						FROM Comments
						WHERE Id = ?`, commentID).Scan(&comment.Id, &comment.UserID, &comment.PhotoID, &comment.Content)
	if e == nil {
		return comment, "Comment found successfully", nil
	}
	if errors.Is(e, sql.ErrNoRows) {
		return comment, "Couldn't find the user", NotFound
	}
	return comment, e.Error(), InternalServerError
}
