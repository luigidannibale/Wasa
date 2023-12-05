package database

import (
	"errors"
)

func (db *appdbimpl) DeleteComment(commentID int) (string, error) {

	res, err := db.c.Exec(`DELETE OR IGNORE 
							FROM Comments
							WHERE CommentId = ?`, commentID)

	if x, y := res.RowsAffected(); x == 0 && y == nil {
		return "No commment with such id", errors.New("NotFound")
	}
	if err != nil {
		return "An error occcurred on the server" + err.Error(), errors.New("InternalServerError")
	}
	return "Comment deleted successfully", nil
}
